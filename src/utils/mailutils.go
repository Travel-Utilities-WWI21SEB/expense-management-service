package utils

import (
	"fmt"
	"log"

	"github.com/matcornic/hermes/v2"
)

var h = hermes.Hermes{
	Product: hermes.Product{
		Name:        "Costventures",
		Link:        "https://expenseui.c930.net",
		Logo:        "https://raw.githubusercontent.com/Travel-Utilities-WWI21SEB/expense-management-ui/main/static/logo.jpeg",
		TroubleText: "If the {ACTION}-button is not working for you, just copy and paste the URL below into your web browser.",
		Copyright:   "Copyright © 2023 Travel-Utilities-WWI21SEB",
	},
	Theme: new(hermes.Default),
}

// PrepareActivationMailBody prepares the body of the activation mail
func PrepareActivationMailBody(activationLink string, username string) string {
	hermesMail := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				fmt.Sprintf("Welcome to Costventures, %v! We're very excited to have you on board.", username),
			},
			Actions: []hermes.Action{
				{
					Instructions: "To confirm your account, please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Confirm your account",
						Link:  activationLink,
					},
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
