package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

//store configuration data containing user details
type configuration struct {
	Email emailConfig
	RT    rtConfig
}

type emailConfig struct {
	Server    string
	Sender    string
	Password  string
	Recipient string
}

type rtConfig struct {
	URL      string
	Username string
	Password string
	Queue    string
}

var configFileName string
var config configuration

//initialize application flags, load user details from config file.
func init() {
	flag.StringVar(&configFileName, "config", "config.json", "Path to config file.")
	flag.Parse()

	f, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
	}

	//myRT := rtgo.NewRT("URL", "username", "password")
	// myRT.CreateTicket("Queue", "requestor", subject, text)
}

//Main runs http server if appropriate flags specified
func main() {
	if config.Email.Server == "" && config.RT.URL == "" {
		fmt.Println("Looks like your configuration file is incomplete.")
		fmt.Println("At least Email server or RT URL MUST be supplied.")
	}
	// if emailServer == "" || emailSender == "" || emailRecipient == "" {
	// 	fmt.Println("emailServer, emailSender and emailRecipient required. Run with --help for info.")
	// } else {
	// 	http.HandleFunc("/transcript", handler)
	// 	http.ListenAndServe(":8080", nil)
	// }
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
	Name      string `json:"name"`
	Email     string `json:"email"`
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()

	fmt.Println(s)

	// var t []transcript
	// dec := json.NewDecoder(r.Body)
	// err := dec.Decode(&t)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	//
	// //Loop through objects, if message is present email the transcript
	// for _, i := range t {
	// 	if i.Message != "" {
	// 		sendEmail(
	// 			config.Email.Server,
	// 			config.Email.Sender,
	// 			config.Email.Password,
	// 			config.Email.Recipient, i)
	// 		fmt.Println(i.toString())
	// 	}
	// }
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
