package models

type ActivationMail struct {
	Username        string   `json:"username"`
	ActivationToken string   `json:"activationToken"`
	Subject         string   `json:"subject"`
	Recipients      []string `json:"recipients"`
}

type ConfirmationMail struct {
	Username   string   `json:"username"`
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
}
