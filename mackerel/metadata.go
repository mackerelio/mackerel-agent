package mackerel

import (
	"net/http"
)

// PutMetadata creates or updates metadata of a host.
func (api *API) PutMetadata(hostID string, namespace string, metadata interface{}) (*http.Response, error) {
	resp, err := api.putJSON("/api/v0/hosts/"+hostID+"/metadata/"+namespace, metadata)
	defer closeResp(resp)

	return resp, err
}
