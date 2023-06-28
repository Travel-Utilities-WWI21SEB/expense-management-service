package utils

import (
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
)

var h = hermes.Hermes{
	Product: hermes.Product{
		Name:        "Costventures",
		Link:        "https://costventures.works",
		Logo:        "https://github.com/Travel-Utilities-WWI21SEB/expense-management-ui/blob/main/static/BannerLogo.png",
		TroubleText: "If the {ACTION}-button is not working for you, just copy and paste the URL below into your web browser.",
		Copyright:   "Copyright © 2023 Travel-Utilities-WWI21SEB",
	},
	Theme: new(hermes.Default),
}

// PrepareActivationMailBody prepares the body of the activation mail
func PrepareActivationMailBody(inviteCode string, username string) string {
	hermesMail := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				fmt.Sprintf("Welcome to Costventures, %v! We're very excited to have you on board.", username),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please copy your verification token:",
					InviteCode:   inviteCode,
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := h.GenerateHTML(hermesMail)
	if err != nil {
		log.Printf("Error in utils.prepareActivationMailBody().GenerateHTML(): %v", err.Error())
		return ""
	}

	return emailBody
}

func PrepareConfirmationMailBody(username string) string {
	hermesMail := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"Thank you for verifying your email! You can now use all features of Costventures.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Costventures please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Go to Costventures",
						Link:  "https://expenseui.c930.net",
					},
				},
			},
			Outros: []string{
				"Have fun and enjoy your trips!",
			},
		},
	}

	emailBody, err := h.GenerateHTML(hermesMail)
	if err != nil {
		log.Printf("Error in utils.prepareConfirmationMailBody().GenerateHTML(): %v", err.Error())
		return ""
	}

	return emailBody
}

func PreparePasswordResetMailBody(username, token string) string {
	hermesMail := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"You have requested a password reset. Please follow the instructions to use Costventures again.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please copy your reset token:",
					InviteCode:   token,
				},
			},
			Outros: []string{
				"If you did not request a password reset, please ignore this email.",
			},
		},
	}

	emailBody, err := h.GenerateHTML(hermesMail)
	if err != nil {
		log.Printf("Error in utils.preparePasswordResetMailBody().GenerateHTML(): %v", err.Error())
		return ""
	}

	return emailBody
}

func PreparePasswordResetConfirmationMailBody(email string) string {
	hermesMail := hermes.Email{
		Body: hermes.Body{
			Name: email,
			Intros: []string{
				"Your password has been successfully reset.",
				"Please notify our support team if you did not request this password reset.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Costventures please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Go to Costventures",
						Link:  "https://expenseui.c930.net",
					},
				},
			},
			Outros: []string{
				"Have fun and enjoy your trips!",
			},
		},
	}

	emailBody, err := h.GenerateHTML(hermesMail)
	if err != nil {
		log.Printf("Error in utils.preparePasswordResetConfirmationMailBody().GenerateHTML(): %v", err.Error())
		return ""
	}

	return emailBody
}
