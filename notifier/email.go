package notifier

import (
	"log"
	"net/mail"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

//EmailAlert - Send email alerts
type EmailAlert struct {
	Body         string
	Subject      string
	ProviderName string
}

//XMLEmailAlert - Send email alerts
type XMLEmailAlert struct {
	Body         string
	Subject      string
	ProviderName string
}

//SendViaChannel - For sending email alerts
func (a EmailAlert) SendViaChannel(eCounter chan int) {

	err := godotenv.Load("/home/cloudiff/.env") //Load Environmental Variables
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		gmailAddress  = os.Getenv("GMAIL_ADDRESS")
		gmailPassword = os.Getenv("GMAIL_APP_PASSWORD")
	)

	if (a.ProviderName == "Localhost") || (a.ProviderName == "LocalHostNmapResults") { //Don't send localhost alerts
		log.Printf("No changes to SEND")
	} else {
		m := gomail.NewMessage()
		m.SetHeader("From", gmailAddress)
		m.SetHeader("To", gmailAddress)
		m.SetHeader("Subject", encodeRFC2047("Cloudiff "+a.Subject))
		m.SetBody("text/html", a.Body)
		d := gomail.NewDialer("smtp.gmail.com", 587, gmailAddress, gmailPassword)
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		log.Printf("SENT %s changes", a.ProviderName)
	}

	eCounter <- 1
}

//SendViaChannel - For sending xml email alerts
func (a XMLEmailAlert) SendViaChannel(eCounter chan int) {

	err := godotenv.Load("/home/cloudiff/.env") //Load Environmental Variables
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		gmailAddress  = os.Getenv("GMAIL_ADDRESS")
		gmailPassword = os.Getenv("GMAIL_APP_PASSWORD")
	)

	if (a.ProviderName == "Localhost") || (a.ProviderName == "LocalHostNmapResults") { //Don't send localhost alerts
		log.Printf("No changes to SEND")
	} else {
		m := gomail.NewMessage()
		m.SetHeader("From", gmailAddress)
		m.SetHeader("To", gmailAddress)
		m.SetHeader("Subject", encodeRFC2047("Cloudiff "+a.Subject))
		m.SetBody("text/html", a.Body)
		d := gomail.NewDialer("smtp.gmail.com", 587, gmailAddress, gmailPassword)
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		log.Printf("SENT %s changes", a.ProviderName)
	}

	eCounter <- 1
}

func encodeRFC2047(subject string) string {
	//mail.Address use mail's rfc2047 to encode any string
	addr := mail.Address{subject, ""}
	return strings.Trim(addr.String(), " <>")
}

//SendIPChangeAlerts - Opens a channel to send IP change alerts
func SendIPChangeAlerts(wg *sync.WaitGroup, alerts chan EmailAlert, aCounter chan int) {
	defer wg.Done()
	for alert := range alerts {
		go alert.SendViaChannel(aCounter)
	}
}

//SendServiceChangeAlerts - Open channel to send service change alerts
func SendServiceChangeAlerts(wg *sync.WaitGroup, alerts chan XMLEmailAlert, aCounter chan int) {
	defer wg.Done()
	for alert := range alerts {
		go alert.SendViaChannel(aCounter)
	}
}

//TrackIPChangeAlerts - numberOfAlerts is equal numberOfProviders
func TrackIPChangeAlerts(numberOfAlerts int, alerts chan EmailAlert, counter chan int) {
	c := 0
	for {
		select {
		case i := <-counter:
			c = c + i
			log.Printf("IP Change Email #%d sent...", c)
			if c == numberOfAlerts {
				close(alerts)
				break
			}
		}
	}
}

//TrackServiceChangeAlerts - numberOfAlerts is equal numberOfProviders
func TrackServiceChangeAlerts(numberOfAlerts int, alerts chan XMLEmailAlert, counter chan int) {
	c := 0
	for {
		select {
		case i := <-counter:
			c = c + i
			log.Printf("Service Change Email #%d sent...", c)
			if c == numberOfAlerts {
				close(alerts)
				break
			}
		}
	}
}
