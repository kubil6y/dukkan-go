package email

import (
	"fmt"
	"log"
	"os"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func YarakMail() {
	from := mail.NewEmail("Dukkan", "lieqb2@gmail.com")
	subject := "Dukkan - Account Activation Code"
	to := mail.NewEmail("Kubilay", "lieqb2@gmail.com")
	plainTextContent := "Activation Code: HASLDFJHALSFHJOIWUEH"
	htmlContent := "<strong>Activation Code: HASLDFJHALSFHJOIWUEH</strong>"
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
}

func ActivationEmail(user *data.User, code string) {
	from := mail.NewEmail("Dukkan", "lieqb2@gmail.com")
	subject := "Dukkan - Account Activation Code"
	to := mail.NewEmail(user.FullName(), user.Email)
	plainTextContent := fmt.Sprintf("Activation Code: %s", code)

	htmlContent := fmt.Sprintf(`
		<div>
			<h4>Dukkan!</h4>
			<p><strong>Activation Code: HASLDFJHALSFHJOIWUEH</strong></p>
			<small>
				Please visit <a href="%s/tokens/activation">activation page.</a>
			</small>
		</div>`, os.Getenv("DOMAIN"))

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
}
