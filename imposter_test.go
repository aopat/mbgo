package mbgo_test

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/senseyeio/mbgo"
	"github.com/senseyeio/mbgo/internal/testutil"
)

func TestImposter_MarshalJSON(t *testing.T) {
	cases := []struct {
		Description string
		Imposter    mbgo.Imposter
		Expected    map[string]interface{}
		Err         error
	}{
		{
			Description: "should marshal the tcp Imposter into the expected JSON",
			Imposter: mbgo.Imposter{
				Port:           8080,
				Proto:          "tcp",
				Name:           "tcp_test_imposter",
				RecordRequests: true,
				AllowCORS:      true,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: &mbgo.TCPRequest{
									RequestFrom: net.IPv4(172, 17, 0, 1),
									Data:        "SGVsbG8sIHdvcmxkIQ==",
								},
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.TCPResponse{
									Data: "Z2l0aHViLmNvbS9zZW5zZXllaW8vbWJnbw==",
								},
							},
						},
					},
				},
			},
			Expected: map[string]interface{}{
				"port":           8080,
				"protocol":       "tcp",
				"name":           "tcp_test_imposter",
				"recordRequests": true,
				"allowCORS":      true,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"data":        "SGVsbG8sIHdvcmxkIQ==",
									"requestFrom": "172.17.0.1",
								},
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"data": "Z2l0aHViLmNvbS9zZW5zZXllaW8vbWJnbw==",
								},
							},
						},
					},
				},
			},
		},
		{
			Description: "should marshal the http Imposter into the expected JSON",
			Imposter: mbgo.Imposter{
				Port:           8080,
				Proto:          "http",
				Name:           "http_test_imposter",
				RecordRequests: true,
				AllowCORS:      true,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: mbgo.HTTPRequest{
									RequestFrom: net.IPv4(172, 17, 0, 1),
									Method:      http.MethodGet,
									Path:        "/foo",
									Query: map[string]string{
										"page": "3",
									},
									Headers: map[string]string{
										"Accept": "application/json",
									},
								},
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.HTTPResponse{
									StatusCode: http.StatusOK,
									Headers: map[string]string{
										"Content-Type": "application/json",
									},
									Body: `{"test":true}`,
								},
							},
						},
					},
				},
			},
			Expected: map[string]interface{}{
				"port":           8080,
				"protocol":       "http",
				"name":           "http_test_imposter",
				"recordRequests": true,
				"allowCORS":      true,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"requestFrom": "172.17.0.1",
									"method":      http.MethodGet,
									"path":        "/foo",
									"query": map[string]string{
										"page": "3",
									},
									"headers": map[string]string{
										"Accept": "application/json",
									},
								},
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"statusCode": 200,
									"headers": map[string]string{
										"Content-Type": "application/json",
									},
									"body": `{"test":true}`,
								},
							},
						},
					},
				},
			},
		},
		{
			Description: "should marshal non-string bodies to JSON",
			Imposter: mbgo.Imposter{
				Port:           8080,
				Proto:          "http",
				Name:           "http_test_imposter",
				RecordRequests: true,
				AllowCORS:      true,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: mbgo.HTTPRequest{
									RequestFrom: net.IPv4(172, 17, 0, 1),
									Method:      http.MethodGet,
									Path:        "/foo",
									Headers: map[string]string{
										"Accept": "application/json",
									},
								},
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.HTTPResponse{
									StatusCode: http.StatusOK,
									Headers: map[string]string{
										"Content-Type": "application/json",
									},
									Body: struct{
										Test bool `json:"test"`
									}{
										Test: true,
									},
								},
							},
						},
					},
				},
			},
			Expected: map[string]interface{}{
				"port":           8080,
				"protocol":       "http",
				"name":           "http_test_imposter",
				"recordRequests": true,
				"allowCORS":      true,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"requestFrom": "172.17.0.1",
									"method":      http.MethodGet,
									"path":        "/foo",
									"headers": map[string]string{
										"Accept": "application/json",
									},
								},
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"statusCode": 200,
									"headers": map[string]string{
										"Content-Type": "application/json",
									},
									"body": map[string]interface{}{"test": true},
								},
							},
						},
					},
				},
			},
		},
		{
			Description: "should include parameters on the predicate if specified",
			Imposter: mbgo.Imposter{
				Port:           8080,
				Proto:          "http",
				Name:           "http_test_imposter",
				RecordRequests: true,
				AllowCORS:      true,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: mbgo.HTTPRequest{
									Method: http.MethodGet,
									Path:   "/foo",
								},
								// include JSONPath parameter
								JSONPath: &mbgo.JSONPath{
									Selector: "$..test",
								},
								// include case sensitive parameter
								CaseSensitive: true,
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.HTTPResponse{
									StatusCode: http.StatusOK,
									Headers: map[string]string{
										"Content-Type": "application/json",
									},
									Body: `{"test":true}`,
								},
							},
						},
					},
				},
			},
			Expected: map[string]interface{}{
				"port":           8080,
				"protocol":       "http",
				"name":           "http_test_imposter",
				"recordRequests": true,
				"allowCORS":      true,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"method": http.MethodGet,
									"path":   "/foo",
								},
								"jsonpath": map[string]interface{}{
									"selector": "$..test",
								},
								"caseSensitive": true,
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"statusCode": 200,
									"headers": map[string]string{
										"Content-Type": "application/json",
									},
									"body": `{"test":true}`,
								},
							},
						},
					},
				},
			},
		},
		{
			Description: "should marshal the expected default response on an http imposter, if provided",
			Imposter: mbgo.Imposter{
				Proto: "http",
				Port:  8080,
				DefaultResponse: mbgo.HTTPResponse{
					StatusCode: http.StatusNotImplemented,
					Mode:       "text",
					Body:       "not implemented",
				},
			},
			Expected: map[string]interface{}{
				"protocol": "http",
				"port":     8080,
				"defaultResponse": map[string]interface{}{
					"statusCode": 501,
					"_mode":      "text",
					"body":       "not implemented",
				},
			},
		},
		{
			Description: "should marshal the expected default response on a tcp imposter, if provided",
			Imposter: mbgo.Imposter{
				Proto: "tcp",
				Port:  8080,
				DefaultResponse: mbgo.TCPResponse{
					Data: "not implemented",
				},
			},
			Expected: map[string]interface{}{
				"protocol": "tcp",
				"port":     8080,
				"defaultResponse": map[string]interface{}{
					"data": "not implemented",
				},
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			ab, err := json.Marshal(c.Imposter)
			testutil.ExpectEqual(t, err, c.Err)
			eb, err := json.Marshal(c.Expected)
			testutil.ExpectEqual(t, err, nil)

			var actual, expected map[string]interface{}
			err = json.Unmarshal(ab, &actual)
			testutil.ExpectEqual(t, err, nil)
			err = json.Unmarshal(eb, &expected)
			testutil.ExpectEqual(t, err, nil)

			testutil.ExpectEqual(t, actual, expected)
		})
	}
}

func TestImposter_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		Description string
		JSON        map[string]interface{}
		Expected    mbgo.Imposter
		Err         error
	}{
		{
			Description: "should unmarshal the JSON into the expected http Imposter",
			JSON: map[string]interface{}{
				"port":             8080,
				"protocol":         "http",
				"name":             "http_imposter",
				"numberOfRequests": 42,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"requestFrom": "172.17.0.1:58112",
									"method":      "POST",
									"path":        "/foo",
									"query": map[string]string{
										"bar": "baz",
									},
									"headers": map[string]string{
										"Content-Type": "application/json",
									},
									"body": `{"predicate":true}`,
								},
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"statusCode": 200,
									"_mode":      "text",
									"headers": map[string]string{
										"Accept":       "application/json",
										"Content-Type": "application/json",
									},
									"body": `{"response":true}`,
								},
							},
						},
					},
				},
			},
			Expected: mbgo.Imposter{
				Port:         8080,
				Proto:        "http",
				Name:         "http_imposter",
				RequestCount: 42,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: mbgo.HTTPRequest{
									RequestFrom: net.IPv4(172, 17, 0, 1),
									Method:      "POST",
									Path:        "/foo",
									Query: map[string]string{
										"bar": "baz",
									},
									Headers: map[string]string{
										"Content-Type": "application/json",
									},
									Body: `{"predicate":true}`,
								},
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.HTTPResponse{
									StatusCode: http.StatusOK,
									Mode:       "text",
									Headers: map[string]string{
										"Accept":       "application/json",
										"Content-Type": "application/json",
									},
									Body: `{"response":true}`,
								},
							},
						},
					},
				},
			},
		},
		{
			Description: "should unmarshal the JSON into the expected tcp Imposter",
			JSON: map[string]interface{}{
				"port":             8080,
				"protocol":         "tcp",
				"name":             "tcp_imposter",
				"numberOfRequests": 4,
				"stubs": []interface{}{
					map[string]interface{}{
						"predicates": []interface{}{
							map[string]interface{}{
								"equals": map[string]interface{}{
									"requestFrom": "172.17.0.1:58112",
									"data":        "SGVsbG8sIHdvcmxkIQ==",
								},
							},
						},
						"responses": []interface{}{
							map[string]interface{}{
								"is": map[string]interface{}{
									"data": "Z2l0aHViLmNvbS9zZW5zZXllaW8vbWJnbw==",
								},
							},
						},
					},
				},
			},
			Expected: mbgo.Imposter{
				Port:         8080,
				Proto:        "tcp",
				Name:         "tcp_imposter",
				RequestCount: 4,
				Stubs: []mbgo.Stub{
					{
						Predicates: []mbgo.Predicate{
							{
								Operator: "equals",
								Request: mbgo.TCPRequest{
									RequestFrom: net.IPv4(172, 17, 0, 1),
									Data:        "SGVsbG8sIHdvcmxkIQ==",
								},
							},
						},
						Responses: []mbgo.Response{
							{
								Type: "is",
								Value: mbgo.TCPResponse{
									Data: "Z2l0aHViLmNvbS9zZW5zZXllaW8vbWJnbw==",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(c.JSON)
			if err != nil {
				t.Fatal(err)
			}

			actual := mbgo.Imposter{}
			err = json.Unmarshal(b, &actual)
			testutil.ExpectEqual(t, err, c.Err)
			testutil.ExpectEqual(t, actual, c.Expected)
		})
	}
}
