import { from, HttpLink, ApolloClient, InMemoryCache } from "@apollo/client";
import { onError } from "@apollo/client/link/error";
import { relayStylePagination } from "@apollo/client/utilities";

// Create error link for unauthorized handler
const errorLink = onError((error) => {
  // Open error input
  const { networkError, forward, operation } = error;
  // Check if error is 401 error
  if (networkError && networkError.statusCode === 401) {
    // Redirect to login page
    window.location.href = "/auth/oidc";
  }
  // Forward
  return forward(operation);
});

// Consolidate to 1 link
const link = from([errorLink, new HttpLink({ uri: "/api/graphql" })]);

// Create client
const client = new ApolloClient({
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          partition: {
            decisionLogs: relayStylePagination(),
            statuses: relayStylePagination(),
          },
        },
      },
    },
  }),
  link,
});

export default client;
