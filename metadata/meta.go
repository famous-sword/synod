package metadata

type (
	// Meta the struct of meta data
	Meta struct {
		Version int    `json:"version"`
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Hash    string `json:"hash"`
	}
)
