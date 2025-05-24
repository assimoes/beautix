package graph

import (
	"html/template"
	"net/http"
)

// SandboxHTML contains the Apollo Sandbox HTML
const SandboxHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GraphQL Sandbox - Beautix</title>
    <style>
        body {
            margin: 0;
            overflow: hidden;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        .loading {
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background: #f5f5f5;
        }
        .loading h1 {
            color: #333;
            margin: 0;
        }
        #embeddableExplorer {
            height: 100vh;
            width: 100vw;
        }
    </style>
</head>
<body>
    <div class="loading">
        <h1>Loading GraphQL Sandbox...</h1>
    </div>
    <div id="embeddableExplorer"></div>
    
    <script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
    <script>
        new window.EmbeddedSandbox({
            target: "#embeddableExplorer",
            initialEndpoint: "{{.Endpoint}}",
            includeCookies: true,
            initialState: {
                document: ` + "`" + `
                    query GetUsers {
                        users {
                            id
                            email
                            firstName
                            lastName
                            fullName
                            role
                            isActive
                            createdAt
                            updatedAt
                        }
                    }
                    
                    query GetCurrentUser {
                        currentUser {
                            id
                            email
                            firstName
                            lastName
                            fullName
                            role
                            isActive
                        }
                    }
                    
                    mutation CreateUser($input: CreateUserInput!) {
                        createUser(input: $input) {
                            id
                            email
                            firstName
                            lastName
                            fullName
                            role
                            isActive
                            createdAt
                        }
                    }
                ` + "`" + `,
                variables: {
                    "input": {
                        "email": "user@example.com",
                        "firstName": "John",
                        "lastName": "Doe",
                        "role": "USER"
                    }
                },
                headers: {
                    "Content-Type": "application/json"
                }
            }
        });
        
        // Hide loading message once sandbox loads
        document.querySelector('.loading').style.display = 'none';
    </script>
</body>
</html>
`

// SandboxHandler creates an HTTP handler for the GraphQL Sandbox
func SandboxHandler(graphqlEndpoint string) http.HandlerFunc {
	tmpl := template.Must(template.New("sandbox").Parse(SandboxHTML))
	
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		
		data := struct {
			Endpoint string
		}{
			Endpoint: graphqlEndpoint,
		}
		
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render sandbox", http.StatusInternalServerError)
			return
		}
	}
}