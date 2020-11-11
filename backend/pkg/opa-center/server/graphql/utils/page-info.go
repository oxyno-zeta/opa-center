package utils

// Pagination information.
type PageInfo struct {
	// Has a next page ?
	HasNextPage bool `json:"hasNextPage"`
	// Has a previous page ?
	HasPreviousPage bool `json:"hasPreviousPage"`
	// Shortcut to first edge cursor in the result chunk.
	StartCursor *string `json:"startCursor"`
	// Shortcut to last edge cursor in the result chunk.
	EndCursor *string `json:"endCursor"`
}
