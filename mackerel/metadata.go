package mackerel

// PutMetadata creates or updates metadata of a host.
func (api *API) PutMetadata(hostID string, namespace string, metadata interface{}) error {
	resp, err := api.putJSON("/api/v0/hosts/"+hostID+"/metadata/"+namespace, metadata)
	defer closeResp(resp)

	return err
}
