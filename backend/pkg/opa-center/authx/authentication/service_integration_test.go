//+build integration

package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	cmocks "github.com/oxyno-zeta/opa-center/pkg/opa-center/config/mocks"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/stretchr/testify/assert"
)

func Test_authentication(t *testing.T) {
	fakeMatchingReg := regexp.MustCompile(".*fake")
	validAuthCfg := &config.OIDCAuthConfig{
		ClientID:      "client-with-secret",
		ClientSecret:  &config.CredentialConfig{Value: "565f78f2-a706-41cd-a1a0-431d7df29443"},
		CookieName:    "oidc",
		EmailVerified: false,
		RedirectURL:   "http://localhost:8080/",
		IssuerURL:     "http://localhost:8088/auth/realms/integration",
	}
	type jwtToken struct {
		IDToken string `json:"id_token"`
	}
	tests := []struct {
		name                       string
		inputUnauthorizedRegexList []*regexp.Regexp
		inputCfg                   *config.Config
		inputRequestURL            string
		inputForgeOIDCHeader       bool
		inputForgeOIDCUsername     string
		inputForgeOIDCPassword     string
		inputAuthorizationHeader   string
		wantErr                    bool
		expectedStatusCode         int
		checkBody                  bool
		expectedBody               string
		expectedHeaders            map[string]string
		expectedUser               *models.OIDCUser
	}{
		{
			name: "wrong issuer url",
			inputCfg: &config.Config{
				OIDCAuthentication: &config.OIDCAuthConfig{
					ClientID:  "fake",
					IssuerURL: "http://fake.com/",
				},
			},
			wantErr: true,
		},
		{
			name: "should redirect to login path",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputRequestURL:    "/fake",
			wantErr:            false,
			expectedStatusCode: 307,
			expectedHeaders: map[string]string{
				"Location": "/auth/oidc?rd=http%3A%2F%2Fexample.com%2Ffake",
			},
		},
		{
			name: "should redirect to oidc login page",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputRequestURL:    "/auth/oidc?rd=http%3A%2F%2Fexample.com%2Ffake",
			wantErr:            false,
			expectedStatusCode: 302,
			expectedHeaders: map[string]string{
				"Location": "http://localhost:8088/auth/realms/integration/protocol/openid-connect/auth?client_id=client-with-secret&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fauth%2Foidc%2Fcallback%3Frd%3Dhttp%253A%252F%252Fexample.com%252Ffake&response_type=code",
			},
		},
		{
			name: "invalid token should redirect to login path",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputAuthorizationHeader: "Bearer FAKE",
			inputRequestURL:          "/fake",
			wantErr:                  false,
			expectedStatusCode:       307,
			expectedHeaders: map[string]string{
				"Location": "/auth/oidc?rd=http%3A%2F%2Fexample.com%2Ffake",
			},
		},
		{
			name: "invalid token should respond unauthorized",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputUnauthorizedRegexList: []*regexp.Regexp{fakeMatchingReg},
			inputAuthorizationHeader:   "Bearer FAKE",
			inputRequestURL:            "/fake",
			wantErr:                    false,
			expectedStatusCode:         401,
			expectedHeaders: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			checkBody:    true,
			expectedBody: "{\"error\":\"unauthorized\",\"extensions\":{\"code\":\"UNAUTHORIZED\"}}",
		},
		{
			name: "invalid token should redirect to login path",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputAuthorizationHeader: "Bearer FAKE",
			inputRequestURL:          "/fake",
			wantErr:                  false,
			expectedStatusCode:       307,
			expectedHeaders: map[string]string{
				"Location": "/auth/oidc?rd=http%3A%2F%2Fexample.com%2Ffake",
			},
		},
		{
			name: "valid token should be ok",
			inputCfg: &config.Config{
				OIDCAuthentication: validAuthCfg,
			},
			inputForgeOIDCHeader:   true,
			inputForgeOIDCUsername: "user",
			inputForgeOIDCPassword: "password",
			inputRequestURL:        "/fake",
			wantErr:                false,
			expectedStatusCode:     200,
			expectedHeaders: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			checkBody:    true,
			expectedBody: "{\"email\":\"sample-user@example.com\", \"email_verified\":true, \"family_name\":\"User\", \"given_name\":\"Sample\", \"name\":\"Sample User\", \"preferred_username\":\"user\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock controller
			ctrl := gomock.NewController(t)
			cfgManager := cmocks.NewMockManager(ctrl)
			cfgManager.EXPECT().GetConfig().AnyTimes().Return(tt.inputCfg)
			cfgManager.EXPECT().AddOnChangeHook(gomock.Any()).Do(func(h func()) {})
			// Create service
			s := &service{
				cfgManager: cfgManager,
			}
			// Create fake engine
			router := gin.New()

			// Add logger middleware
			router.Use(log.Middleware(
				log.NewLogger(),
				func(c *gin.Context) string { return "fake" },
				func(c context.Context) string { return "fake" },
			))

			// Add oidc login endpoints
			err := s.OIDCEndpoints(router)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			// Add authentication middleware
			router.Use(s.Middleware(tt.inputUnauthorizedRegexList))
			// Add a fake endpoint
			router.GET("/fake", func(c *gin.Context) {
				us := GetAuthenticatedUserFromGin(c)
				c.JSON(http.StatusOK, us)
			})

			// Create writer
			w := httptest.NewRecorder()
			// Create request
			request := httptest.NewRequest("GET", tt.inputRequestURL, nil)

			if tt.inputForgeOIDCHeader {
				data := url.Values{}
				data.Set("username", tt.inputForgeOIDCUsername)
				data.Set("password", tt.inputForgeOIDCPassword)
				data.Set("client_id", "client-with-secret")
				data.Set("client_secret", "565f78f2-a706-41cd-a1a0-431d7df29443")
				data.Set("grant_type", "password")
				data.Set("scope", "openid profile email")

				authentUrlStr := "http://localhost:8088/auth/realms/integration/protocol/openid-connect/token"

				clientAuth := &http.Client{}
				r, err := http.NewRequest("POST", authentUrlStr, strings.NewReader(data.Encode())) // URL-encoded payload
				// Check err
				if err != nil {
					t.Error(err)
					return
				}

				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

				resp, err := clientAuth.Do(r)
				// Check err
				if err != nil {
					t.Error(err)
					return
				}

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				body := string(bodyBytes)

				// Check response
				if resp.StatusCode != 200 {
					t.Errorf("%d - %s", resp.StatusCode, body)
					return
				}

				var to jwtToken
				// Parse token
				err = json.Unmarshal(bodyBytes, &to)
				if err != nil {
					t.Error(err)
					return
				}

				// Add header to request
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", to.IDToken))
			}

			if tt.inputAuthorizationHeader != "" {
				request.Header.Add("Authorization", tt.inputAuthorizationHeader)
			}

			// Run server
			router.ServeHTTP(w, request)

			// Tests
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.checkBody {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			if tt.expectedHeaders != nil {
				for key, val := range tt.expectedHeaders {
					wheader := w.HeaderMap.Get(key)
					if val != wheader {
						assert.Equal(t, val, wheader, "header = "+key)
					}
				}
			}
		})
	}
}
