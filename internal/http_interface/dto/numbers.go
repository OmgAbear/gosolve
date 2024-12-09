package dto

// NumbersResult is the response structure of the API.
// While it is also used inside the infrastructure layer, in my "best practices" opinion, this should not be the case
// and each layer should generally have their respective objects (at least DTOs)
// make use of common "domain" objects (as available in DDD)
// as well as have mappings to and from those objects as required
type NumbersResult struct {
	Index   int     `json:"index"`
	Value   int     `json:"value"`
	Message *string `json:"message,omitempty"`
}
