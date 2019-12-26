package googleCalendar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	googleCredentials "xcal/config/google"

	"github.com/dustin/go-humanize"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func Init() {
	b, err := json.Marshal(googleCredentials.Get())

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	initClient(config)
}

func GetNextEvent(max int64, truncate *string) {
	b, err := json.Marshal(googleCredentials.Get())

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(max).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Printf("--:-- Nothing")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime

			if date != "" {
				parsed, _ := time.Parse(time.RFC3339, date)

				delay := fmt.Sprintf("%s", humanize.Time(parsed))
				delay = strings.Replace(delay, " from now", "", -1)
				title := fmt.Sprintf("%."+*truncate+"s", item.Summary)

				fmt.Printf("(in %s) %s", delay, title)

				break
			}
		}
	}
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	encodingError := json.NewEncoder(f).Encode(token)

	if encodingError != nil {
		panic("Error! Couldn't encode the token.")
	}
}

func getClient(config *oauth2.Config) *http.Client {
	homeDir, _ := os.UserHomeDir()
	tokFile := homeDir + "/.google-calendar-token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		panic("Error! Google calendar was not initialized.")
	}
	return config.Client(context.Background(), tok)
}

func initClient(config *oauth2.Config) *http.Client {
	homeDir, _ := os.UserHomeDir()
	tokFile := homeDir + "/.google-calendar-token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", authURL).Start()
	case "windows", "darwin":
		err = exec.Command("open", authURL).Start()
	default:
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	var authCode string
	fmt.Print("Google Calendar is being opened in your browser. Sign in with your account, and when you're done signing in, paste here the code that will be shown to you. \n\nCode: ")

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
