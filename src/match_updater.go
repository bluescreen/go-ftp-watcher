package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/jlaffaye/ftp"
	"google.golang.org/api/option"
)

var lastUpdated []time.Time

func run(forceUpdate bool) {
	lastUpdated = readLastUpdated(forceUpdate)

	//clearTerminal()
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
				processMatch(i, number, entry)
			} else {
				fmt.Println(number, "Last update", previousUpdate.Format(dateLayout))
			}
		}
	}

	storeLastUpdated(lastUpdated)
	closeFTP()
}

func readLastUpdated(forceUpdate bool) []time.Time {
	data, err := readLines("lastUpdate.txt")
	if err != nil {
		fmt.Println("No last updated found")
	}
	mapped := make([]time.Time, 4, 4)
	if forceUpdate {
		return mapped
	}
	for i, updatedString := range data {
		updatedTime, _ := time.Parse(dateLayout, updatedString)
		mapped[i] = updatedTime
	}
	return mapped
}

func initFirestore() {
	opt := option.WithCredentialsFile(firestoreCredentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	check(err)

	firestoreClient, err = app.Firestore(ctx)
	check(err)
}

func download(entry *ftp.Entry) MatchInterface {
	r, err := ftpClient.Retr(entry.Name)
	defer r.Close()
	check(err)

	buf, err := ioutil.ReadAll(r)
	return ParseMatch(string(buf))
}

func processMatch(i int, number string, entry *ftp.Entry) {
	match := download(entry)
	event := makeEvent(number, match)
	lastUpdated[i-1] = entry.Time
	fmt.Println(number, "Update", entry.Time.Format(dateLayout), match)
	log.Println(event.Description)
	trackEvent(event)
	storeMatch(match, i)
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

func storeMatch(match MatchInterface, n int) {
	colref := firestoreClient.Collection(match.Table())
	result, err := colref.Doc("Shiaijo"+toCharStrConst(n+1)).Set(ctx, match)
	check(err)
	fmt.Println(n, result, match)
}

func storeLastUpdated(lastUpdated []time.Time) {
	mapped := make([]string, 0, 4)
	for _, updated := range lastUpdated {
		mapped = append(mapped, updated.Format(dateLayout))
	}
	writeLines(mapped, "lastUpdate.txt")
}
