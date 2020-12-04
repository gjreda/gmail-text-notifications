package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Config contains necessary app configuration data
type Config struct {
	Twilio struct {
		AccountSID  string
		AuthToken   string
		PhoneNumber string
		BaseURL     string
	}
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config := &Config{}
	err = json.NewDecoder(f).Decode(config)
	return config, err
}

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	token, err := readTokenFile(tokenFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(tokenFile, token)
	}
	return config.Client(context.Background(), token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to fetch token from web: %v", err)
	}
	return token
}

func readTokenFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func queryMessages(service *gmail.Service, user string, q string) []*gmail.Message {
	log.Printf("Searching for messages containing: %v", q)
	response, err := service.Users.Messages.List(user).Q(q).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}
	if response.HTTPStatusCode != 200 {
		log.Printf("Request returned status code: %v\n", response.HTTPStatusCode)
	}
	log.Printf("Number of messages found: %v\n", len(response.Messages))
	return response.Messages
}

func buildSMS(service *gmail.Service, user string, messages []*gmail.Message, q string, includeSnippets bool) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Hi! You have %v emails matching your search of \"%v\".", len(messages), q)
	if len(messages) == 0 {
		return ""
	}
	if includeSnippets == true {
		fmt.Fprintf(&sb, " Here's what they look like.\n")
		for i, m := range messages {
			fmt.Printf("(%v) Fetching message %v\n", i+1, m.Id)
			m, err := service.Users.Messages.Get(user, m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message ID %v: %v", m.Id, err)
			}
			fmt.Fprintf(&sb, "(%v) - %v\n", i+1, m.Snippet)
		}
	}
	return sb.String()
}

func sendSMS(phoneNumber string, message string, config *Config) {
	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", config.Twilio.PhoneNumber)
	msgData.Set("Body", message)
	reader := *strings.NewReader(msgData.Encode())

	reqURL := config.Twilio.BaseURL + "/Accounts/" + config.Twilio.AccountSID + "/Messages.json"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", reqURL, &reader)
	req.SetBasicAuth(config.Twilio.AccountSID, config.Twilio.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, _ := client.Do(req)
	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&data)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if err == nil {
			log.Printf("Twilio message SID: %v", data["sid"])
		}
	} else {
		log.Printf("Twilio returned status: %v", response.Status)
	}
}

func main() {
	var searchQ string
	var phoneNumber string
	user := "me"

	flag.StringVar(&searchQ, "q", "", "Keyword(s) to search for")
	flag.StringVar(&phoneNumber, "phone", "", "Phone number to text")
	flag.Parse()

	cfg, err := getConfig("config.json")

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}
	client := getClient(config)

	service, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	messages := queryMessages(service, user, searchQ)
	if len(messages) == 0 {
		os.Exit(0)
	}
	msg := buildSMS(service, user, messages, searchQ, false)
	sendSMS(phoneNumber, msg, cfg)
}
