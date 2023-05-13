package model

// LifeCheckResponse is the response for the lifecheck endpoint
type LifeCheckResponse struct {
	Alive   bool   `json:"alive"`
	Version string `json:"version"`
}
