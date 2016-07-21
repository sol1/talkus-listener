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
	flag.StringVar(&emailRecipient, "emailRecipient", "", "Email address used to recieve the transcript")
	flag.StringVar(&emailSender, "emailSender", "", "Email address used to send tickets.")
	flag.StringVar(&emailPassword, "emailPassword", "", "(OPTIONAL) Password associated with the sending email address. If none supplied there will be no authentication")
	flag.Parse()
}

//Main runs http server if appropriate flags specified
func main() {
	if emailServer == "" || emailSender == "" || emailRecipient == "" {
		fmt.Println("emailServer, emailSender and emailRecipient required. Run with --help for info.")
	} else {
		http.HandleFunc("/transcript", handler)
		http.ListenAndServe(":8080", nil)
	}
}

//Stores message data, includes requried JSON tags
type transcript struct {
	Event       string   `json:"event"`
	CreatedAt   string   `json:"createdAt"`
	ChannelName string   `json:"channelName"`
	VisitorName string   `json:"visitorName"`
	VisitorID   string   `json:"visitorId"`
	Identity    identity `json:"identity"`
	Message     string   `json:"message"`
}

type identity struct {
	UserAgent string `json:"userAgent"`
	Location  string `json:"location"`
	Title     string `json:"title"`
	IP        string `json:"ip"`
	Languages string `json:"languages"`
}

//Gets transcript data in a readable format
func (t transcript) toString() string {
	createdTime, err := time.Parse(time.RFC3339, t.CreatedAt)
	if err != nil {
		fmt.Printf("Error parsing time: %s\n", err)
	}

	return fmt.Sprintf("Visitor name: %s\nChannel name: %s\nVisitor ID: %s\nArchived on: %s\nUser agent: %s\nUser IP: %s\n\n%s",
		t.VisitorName,
		t.ChannelName,
		t.VisitorID,
		createdTime.Local().Format(time.Stamp),
		t.Identity.UserAgent,
		t.Identity.IP,
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
			sendEmail(emailServer, emailSender, emailPassword, emailRecipient, i)
			fmt.Println(i.toString())
		}
	}
}

//Send message to email
func sendEmail(emailServer string, emailSender string, emailPassword string, emailRecipient string, trans transcript) {
	// Set up authentication information.
	var auth smtp.Auth
	if emailPassword != "" {
		auth = smtp.PlainAuth("", emailSender, emailPassword, emailServer)
	}

	// Set up mail headers and content
	to := []string{emailRecipient}
	msg := []byte("To: " + emailRecipient + "\r\n" +
		"Subject: Live support with " + trans.VisitorName + "\r\n" +
		"\r\n" +
		trans.toString() + "\r\n")

	//Send mail, if password supplied use auth.
	err := smtp.SendMail(emailServer+":25", auth, emailSender, to, msg)

	if err != nil {
		fmt.Printf("Failed to send mail: %s\n", err)
	} else {
		fmt.Println("Message sent..")
	}
}
