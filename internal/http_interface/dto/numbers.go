package dto

type NumbersResult struct {
	Index   int     `json:"index"`
	Value   int     `json:"value"`
	Message *string `json:"message,omitempty"`
}
