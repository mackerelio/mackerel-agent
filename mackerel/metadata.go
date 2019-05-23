package mackerel

// PutMetadata creates or updates metadata of a host.
func (api *API) PutMetadata(hostID string, namespace string, metadata interface{}) error {
	return api.c.PutHostMetaData(hostID, namespace, metadata)
}
