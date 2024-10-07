package docprepsrv

import (
	"net/http"
	client "simrep/pkg/docprepsrv"
)

type (
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

type (
	apiClient = client.ClientWithResponsesInterface
)
