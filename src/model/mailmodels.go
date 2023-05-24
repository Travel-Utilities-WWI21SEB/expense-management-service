package model

type ActivationMail struct {
	Username      string   `json:"username"`
	ActivationUrl string   `json:"activationUrl"`
	Subject       string   `json:"subject"`
	Recipients    []string `json:"recipients"`
}

type ConfirmationMail struct {
	Username   string   `json:"username"`
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
}
