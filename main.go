package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/jlaffaye/ftp"
	"google.golang.org/api/option"
)

var ftpConfig = "127.0.0.1:8082"
var ftpUser = "anonymous"
var ftpPassword = "anonymous"
var dateLayout = "2006-01-02 15:04:05"
var pollInterval time.Duration = 3 * time.Second

var firebaseConfig *firebase.Config = &firebase.Config{ProjectID: "project-4117648448"}
var firestoreClient *firestore.Client
var firestoreCredentialsFile = "ekc-stream-firebase-adminsdk-ses37-27da0fc036.json"
var ctx = context.Background()
var ftpClient *ftp.ServerConn

var matches []*Match
var entries []*ftp.Entry

func init() {
	fmt.Println("Init")

}

func main() {
	fmt.Println("Starting ...")

	lastUpdated := readLastUpdated()
	for {
		fmt.Println("Polling Data", time.Now().Format(dateLayout))
		run(lastUpdated)
		time.Sleep(pollInterval)
	}
}

func run(lastUpdated []time.Time) {
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
				lastUpdated[i-1] = entry.Time
				fmt.Println(number, "Update", entry.Time.Format(dateLayout), match)
				storeMatch(match, i)
			} else {
				fmt.Println(number, "Last updated", previousUpdate.Format(dateLayout))
			}
		}
	}

	storeLastUpdated(lastUpdated)
	closeFTP()
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

	fmt.Println("firebase app is initialized.")
}

func download(entry *ftp.Entry) Match {
	r, err := ftpClient.Retr(entry.Name)
	defer r.Close()
	check(err)

	buf, err := ioutil.ReadAll(r)

	reader := csv.NewReader(strings.NewReader(string(buf)))
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
