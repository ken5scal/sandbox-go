package examples

import (
	// "net/http"
	"net/http"
	"time"

	"github.com/ken5scal/go-httpclient/gohttp"
	"github.com/ken5scal/go-httpclient/gomime"
)

var httpClient = getHttpClient()

func getHttpClient() gohttp.Client {
	headers := make(http.Header)
	headers.Set(gomime.HeaderContentType, gomime.ContentTypeJson)
	// currentClient := http.Client{}
	client := gohttp.NewBuilder().
		SetHeaders(headers).
		SetConnectionTimeout(2 * time.Second).
		SetResponseTimeout(5 * time.Second).
		SetUserAgent("Fedes-Computer").
		// SetHttpClient(&currentClient).
		Build()
	return client
}
