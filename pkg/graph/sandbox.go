package graph

import (
	"net/http"
	"strings"
)

// SandboxHandler serves the Apollo Sandbox for GraphQL API exploration
func SandboxHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sandboxHTML))
	}
}

// sandboxHTML is the HTML content for Apollo Sandbox
const sandboxHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>BeautyBiz API Sandbox</title>
  <style>
    body { margin: 0; padding: 0; }
    #sandbox { position: absolute; top: 0; right: 0; bottom: 0; left: 0; }
  </style>
</head>
<body>
  <div id="sandbox"></div>
  <script src="https://embeddable-sandbox.cdn.apollographql.com/v2/embeddable-sandbox.umd.production.min.js"></script>
  <script>
    // Get the current origin
    const origin = window.location.origin;
    // Construct the GraphQL endpoint URL
    const endpoint = origin + "/graphql";

    // Initialize the sandbox
    new window.EmbeddedSandbox({
      target: "#sandbox",
      initialEndpoint: endpoint,
      initialState: {
        document: "# Welcome to BeautyBiz API Sandbox\n# Try a query like:\nquery GetProviders {\n  providers(limit: 10, offset: 0) {\n    id\n    businessName\n    description\n    city\n    user {\n      firstName\n      lastName\n    }\n  }\n}"
      },
      includeCookies: true,
    });
  </script>
</body>
</html>
`

// CorsMiddleware adds CORS headers to all responses
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware adds authentication middleware to protect routes
func AuthMiddleware(userService interface{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public paths that don't require authentication
		publicPaths := []string{
			"/graphql",
			"/sandbox",
			"/health",
		}

		// Check if the path is public
		for _, path := range publicPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Check for authentication token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Bearer token format
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Token validation would be done here
		// For now, we'll just pass through since this is handled at the GraphQL layer

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}