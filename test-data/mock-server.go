package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/api/rulesets/evaluate", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check for Authorization header (accept any token for testing)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Status":  "Error",
				"Message": "Missing or invalid Authorization header",
				"Meta":    nil,
			})
			return
		}

		// Mock response based on the example in requirements
		response := []map[string]interface{}{
			{
				"code":     "owasp-define-error-responses-401",
				"path":     []string{"paths", "/", "get", "responses"},
				"message":  "missing response code `401` for `GET`",
				"severity": 1,
				"range": map[string]interface{}{
					"start": map[string]interface{}{
						"line":      1,
						"character": 194,
					},
					"end": map[string]interface{}{
						"line":      1,
						"character": 205,
					},
				},
				"source": "684acc5b0e08080001e72b3a",
				"api": map[string]interface{}{
					"id":   "684acc5b0e08080001e72b3a",
					"name": "testing-rest-api-2025-05",
				},
				"rule": map[string]interface{}{
					"name": "owasp-define-error-responses-401",
				},
			},
			{
				"code":     "owasp-rate-limit",
				"path":     []string{"paths", "/", "get", "responses", "200"},
				"message":  "response with code `200`, must contain one of the defined headers: `{X-RateLimit-Limit} {X-Rate-Limit-Limit} {RateLimit-Limit, RateLimit-Reset} {RateLimit} `",
				"severity": 0,
				"range": map[string]interface{}{
					"start": map[string]interface{}{
						"line":      1,
						"character": 207,
					},
					"end": map[string]interface{}{
						"line":      1,
						"character": 212,
					},
				},
				"source": "684acc5b0e08080001e72b3a",
				"api": map[string]interface{}{
					"id":   "684acc5b0e08080001e72b3a",
					"name": "testing-rest-api-2025-05",
				},
				"rule": map[string]interface{}{
					"name": "owasp-rate-limit",
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock governance service starting on :8989")
	log.Fatal(http.ListenAndServe(":8989", nil))
}
