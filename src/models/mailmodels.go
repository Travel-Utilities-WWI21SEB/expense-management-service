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

type PasswordResetMail struct {
	Username   string   `json:"username"`
	ResetToken string   `json:"resetToken"`
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
}

type ResetPasswordConfirmationMail struct {
	Username   string   `json:"username"`
	Subject    string   `json:"subject"`
	Recipients []string `json:"recipients"`
}

type SendContactMailRequest struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
