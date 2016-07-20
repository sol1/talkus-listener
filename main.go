package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/smtp"
	"time"
)

var emailServer string
var emailSender string
var emailPassword string
var emailRecipient string

//initialize application flags
func init() {
	flag.StringVar(&emailServer, "emailServer", "", "SMTP server used to send mail from, requires associated username and password")
	flag.StringVar(&emailSender, "emailSender", "", "Email address used to send tickets.")
	flag.StringVar(&emailPassword, "emailPassword", "", "Password associated with the sending email address.")
	flag.StringVar(&emailRecipient, "emailRecipient", "", "Email address used to recieve the transcript")
	flag.Parse()
}

//Main runs http server if appropriate flags specified
func main() {
	if emailServer == "" || emailSender == "" || emailPassword == "" || emailRecipient == "" {
		fmt.Println("All flags are required. Run with --help for info.")
	} else {
		http.HandleFunc("/transcript", handler)
		http.ListenAndServe(":8080", nil)
	}
}

//Stores message data, includes requried JSON tags
type transcript struct {
	Event       string `json:"event"`
	CreatedAt   string `json:"createdAt"`
	ChannelName string `json:"channelName"`
	VisitorName string `json:"visitorName"`
	VisitorID   string `json:"visitorId"`
	Message     string `json:"message"`
}

//Gets transcript data in a readable format
func (t transcript) toString() string {
	createdTime, err := time.Parse(time.RFC3339, t.CreatedAt)
	if err != nil {
		fmt.Printf("Error parsing time: %s\n", err)
	}
	currentTime := time.Now()

	return fmt.Sprintf("Visitor name: %s\nCreated on: %s\nArchived on: %s\n\n%s",
		t.VisitorName,
		createdTime.Local().Format(time.Stamp),
		currentTime.Format(time.Stamp),
		t.Message)
}

//Handles incoming requests under the /transcript path
func handler(w http.ResponseWriter, r *http.Request) {
	var t []transcript
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&t)
	if err != nil {
		fmt.Println("error:", err)
	}

	//Loop through objects, if message is present email the transcript
	for _, i := range t {
		if i.Message != "" {
			sendEmail(emailServer, emailSender, emailPassword, emailRecipient, i.toString())
			fmt.Println(i.toString())
		}
	}
}

//Send message to email
func sendEmail(emailServer string, emailSender string, emailPassword string, emailRecipient string, message string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", emailSender, emailPassword, emailServer)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{emailRecipient}
	msg := []byte("To: " + emailRecipient + "\r\n" +
		"Subject: SMTP test\r\n" +
		"\r\n" +
		message + "\r\n")
	err := smtp.SendMail(emailServer+":25", auth, emailSender, to, msg)
	if err != nil {
		fmt.Printf("Failed to send mail: %s\n", err)
	} else {
		fmt.Println("Message sent..")
	}
}
