package email

import (
	"fmt"

	"github.com/kubil6y/dukkan-go/internal/data"
)

// TODO
func ActivationEmail(user *data.User, code string) {
	/*
		from := mail.NewEmail("Dukkan", "lieqb2@gmail.com")
		subject := "Dukkan - Account Activation Code"
		to := mail.NewEmail(user.FullName(), user.Email)
		plainTextContent := fmt.Sprintf("Activation Code: %s", code)

		htmlContent := fmt.Sprintf(`
				<div>
					<h4>Dukkan!</h4>
					<p><strong>Activation Code: %s</strong></p>
					<small>
						Please visit <a href="%s/tokens/activation">activation page.</a>
					</small>
				</div>`, code, os.Getenv("DOMAIN"))

		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	*/

	fmt.Println("email sent in background...")
}
