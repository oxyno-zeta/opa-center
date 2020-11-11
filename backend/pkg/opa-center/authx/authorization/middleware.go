package authorization

import (
	"github.com/gin-gonic/gin"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/utils"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

func (s *service) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logger
		logger := log.GetLoggerFromGin(c)
		// Get user from request
		ouser := authentication.GetAuthenticatedUserFromGin(c)

		authorized, err := s.isRequestAuthorized(c.Request, ouser)
		// Check error
		if err != nil {
			logger.Error(err)
			utils.AnswerWithError(c, err)

			return
		}

		// Check if user is authorized
		if !authorized {
			err = errors.NewForbiddenError("forbidden")

			logger.Error(err)
			utils.AnswerWithError(c, err)

			return
		}

		// User is authorized

		logger.Infof("OIDC user %s authorized", ouser.GetIdentifier())
		c.Next()
	}
}
