package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/jlaffaye/ftp"
	"google.golang.org/api/option"
)

var firebaseConfig *firebase.Config = &firebase.Config{ProjectID: firestoreProjectId}
var firestoreClient *firestore.Client
var ctx = context.Background()
var ftpClient *ftp.ServerConn

var matches []*Match
var entries []*ftp.Entry

func main() {
	fmt.Println("Starting ...")
	initEventLog()
	lastUpdated := readLastUpdated()
	for {
		run(lastUpdated)
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func initEventLog() {
	file, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	log.SetOutput(file)
}

func run(lastUpdated []time.Time) {
	clearTerminal()
	fmt.Println("Polling Data", time.Now().Format(dateLayout))

	initFirestore()
	defer firestoreClient.Close()

	connectFTP()
	data := readFileList()
	for i, entry := range data {

		if strings.Contains(entry.Name, "txt") {
			number := toCharStrConst(i)

			previousUpdate := lastUpdated[i-1]
			if entry.Time.After(previousUpdate) {
				match := download(entry)
				event := makeEvent(number, match)
				lastUpdated[i-1] = entry.Time
				fmt.Println(number, "Update", entry.Time.Format(dateLayout), match)
				log.Println(event.Description)
				// trackEvent(event)
				storeMatch(match, i)
			} else {
				fmt.Println(number, "Last update", previousUpdate.Format(dateLayout))
			}
		}
	}

	storeLastUpdated(lastUpdated)
	closeFTP()

}

func trackEvent(event Event) {
	/*
		err := client.Track("13793", "Signed Up", map[string]interface{}{
			"Referred By": "Friend",
		})*/
}

func makeEvent(number string, match Match) Event {

	//pointsRed := []string{"foo", "bar", "hello"}
	//ointsWhite := []string{"foo", "bar", "world"}

	// type = Start, ScoreRed, ScoreWhite, Hansoku, Encho, Finished

	//Pool 8, Fight 2 on Shiaijo A. Furl vs. Somame. Furl could achieve two men (M),
	//Somame could achieve a tsuki (T). Furl has already has a hikiwake as well (h)
	description := fmt.Sprintf("Pool %d, Fight %s on Shiaijo %s. %s vs %s",
		match.Pool, match.Fight, match.Shiaijo, match.NameTareRed,
		match.NumberTareWhite)

	return Event{
		Type:        "Start",
		Description: description,
	}

}

func difference(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func readLastUpdated() []time.Time {
	data, err := readLines("lastUpdate.txt")
	if err != nil {
		fmt.Println("No last updated found")
	}
	mapped := make([]time.Time, 4, 4)
	for i, updatedString := range data {
		updatedTime, _ := time.Parse(dateLayout, updatedString)
		mapped[i] = updatedTime
	}
	return mapped
}

func storeLastUpdated(lastUpdated []time.Time) {
	mapped := make([]string, 0, 4)
	for _, updated := range lastUpdated {
		mapped = append(mapped, updated.Format(dateLayout))
	}
	writeLines(mapped, "lastUpdate.txt")
}

func storeMatch(match Match, n int) {
	// TODO:  https://github.com/dukex/mixpanel
	colref := firestoreClient.Collection("matches")
	result, err := colref.Doc("Shiaijo"+toCharStrConst(n+1)).Set(ctx, match)
	check(err)
	fmt.Println(n, result, match)
}

func toCharStrConst(i int) string {
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return abc[i-1 : i]
}

func connectFTP() {
	c, err := ftp.Dial(ftpConfig, ftp.DialWithTimeout(5*time.Second))
	check(err)

	ftpClient = c

	err = ftpClient.Login(ftpUser, ftpPassword)
	check(err)
}

func closeFTP() {
	if err := ftpClient.Quit(); err != nil {
		log.Fatal(err)
	}
}

func readFileList() []*ftp.Entry {
	entries, err := ftpClient.List("/")
	check(err)
	return entries
}

func initFirestore() {
	opt := option.WithCredentialsFile(firestoreCredentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	check(err)

	firestoreClient, err = app.Firestore(ctx)
	check(err)
}

func download(entry *ftp.Entry) Match {
	r, err := ftpClient.Retr(entry.Name)
	defer r.Close()
	check(err)

	buf, err := ioutil.ReadAll(r)
	return ParseMatch(string(buf))
}

func ParseMatch(data string) Match {
	reader := csv.NewReader(strings.NewReader(data))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	check(err)

	match := MakeMatch(records[0])
	matches = append(matches, &match)
	return match
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
