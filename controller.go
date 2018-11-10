package gozw

// Controller contains information for the controller.
type Controller struct {
	APIVersion          string `json:"apiversion"`
	APILibraryType      string `json:"apilibrary_type"`
	HomeID              uint32 `json:"home_id"`
	NodeID              byte   `json:"node_id"`
	Version             byte   `json:"version"`
	APIType             string `json:"apitype"`
	IsPrimaryController bool   `json:"is_primary_controller"`
	ApplicationVersion  byte   `json:"application_version"`
	ApplicationRevision byte   `json:"application_revision"`
	SupportedFunctions  []byte `json:"supported_functions"`
	NodeList            []byte `json:"node_list"`
}
