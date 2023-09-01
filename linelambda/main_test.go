package linelambda

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"testing"
	// "fmt"
)

var (
	body = `
  { "key": "value" }
`
	input = `
  {
    "path": "/some/path",
    "httpMethod":"POST",
    "headers": {
      "Host": "www.example.com",
      "Content-Type": "application/json",
      "X-Line-Signature": "SIG_VALUE",
      "User-Agent": "Test/0.0.0",
      "Should-Not-Propagate": "This-Header"
    },
    "IsBase64Encoded": true
  }
`
)

func TestHandler(t *testing.T) {
	var inputRequest events.APIGatewayProxyRequest
	err := json.Unmarshal([]byte(input), &inputRequest)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	inputRequest.Body = body
	inputRequest.IsBase64Encoded = false
	result, err := APIEventToHTTPRequest(inputRequest)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(result.Header["X-Line-Signature"]) != 1 || result.Header["X-Line-Signature"][0] != "SIG_VALUE" {
		t.Errorf("Incorrect X-Line-Signature: %v", result.Header["X-Line-Signature"])
	}
	if _, exists := result.Header["Should-Not-Propagate"]; exists {
		t.Errorf("Unexpected Header found in result: %v", result.Header)
	}
	// fmt.Printf("%+v\n", result)
}
