/*
Package linelambda implements helper functions that make it easier to use
the Line Bot Go SDK in AWS Lambda, as a serverless API function.

Line Bot Go SDK: https://github.com/line/line-bot-sdk-go/

Example code and an AWS SAM deployment template for a simple
"echo" Line Bot are here: https://github.com/JamesJJ/line-bot-demo
*/
package linelambda

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

func main() {}

// Shortcut that provides HTTP Header Keys in canonical format
func canonical(s string) string {
	return http.CanonicalHeaderKey(s)
}

// Takes the `APIGatewayProxyRequest` event provided by AWS Lambda
// and returns a pointer to a http.Request that is compatible with
// the Line Bot SDK's `ParseRequest` function.
// https://pkg.go.dev/github.com/line/line-bot-sdk-go@v1.3.0/linebot#Client.ParseRequest
func APIEventToHTTPRequest(event events.APIGatewayProxyRequest) (*http.Request, error) {
	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://%s%s",
			event.Headers[canonical("Host")],
			event.Path,
		),
		bytes.NewBufferString(event.Body),
	)
	if err != nil {
		return nil, err
	}
	for k, v := range event.Headers {
		for _, h := range []string{"Content-Type", "X-Line-Signature", "User-Agent"} {
			if canonical(k) == canonical(h) {
				httpReq.Header.Set(k, v)
			}
		}
	}
	return httpReq, nil
}

// A helper that returns a standard response for the API Gateway
// (A http/200 response is usually suitable for web hook calls)
// The string and error parameters are logged if non-empty/non-nil
func Finished(msg string, e error) (events.APIGatewayProxyResponse, error) {
	if e != nil {
		log.Printf("ERROR: %s %v\n", msg, e)
	} else if msg != "" {
		log.Printf("OK: %s\n", msg)
	}
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}
