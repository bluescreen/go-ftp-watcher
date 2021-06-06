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

type Match struct {
	Shiaijo         string
	Pool            string
	Fight           string
	NumberTareWhite string
	NameTareWhite   string
	HansokuWhite    string
	IpponWhite2     string
	IpponWhite1     string
	EnchoOrHikiwake string
	IpponRed1       string
	IpponRed2       string
	HansokuRed      string
	NumberTareRed   string
	NameTareRed     string
}

var firebaseConfig *firebase.Config = &firebase.Config{ProjectID: "project-4117648448"}
var firestoreClient *firestore.Client
var ctx = context.Background()
var ftpClient *ftp.ServerConn

var matches []*Match
var entries []*ftp.Entry

func main() {
	fmt.Println("Starting ...")
	initFirestore()
	defer firestoreClient.Close()

	lastUpdated := make([]time.Time, 4, 4)
	for {
		fmt.Println("Polling Data")
		run(lastUpdated)
		time.Sleep(1e10)
	}
}

func run(lastUpdated []time.Time) {
	connectFTP()
	data := readFileList()
	for i, entry := range data {
		if strings.Contains(entry.Name, "txt") {
			previousUpdate := lastUpdated[i-1]
			if entry.Time.After(previousUpdate) {
				match := download(entry)
				lastUpdated[i-1] = entry.Time
				fmt.Println("Update", entry.Time, match)
				storeMatch(match, i)
			} else {
				fmt.Println("Last updated", previousUpdate)
			}
		}
	}
	closeFTP()
}

func storeMatch(match Match, n int) {
	colref := firestoreClient.Collection("matches")
	result, err := colref.Doc("Shiaijo"+toCharStrConst(n+1)).Set(ctx, match)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n, result, match)
}

func toCharStrConst(i int) string {
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return abc[i-1 : i]
}

func connectFTP() {
	c, err := ftp.Dial("127.0.0.1:8082", ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		log.Fatal(err)
	}
	ftpClient = c

	err = ftpClient.Login("anonymous", "anonymous")
	if err != nil {
		log.Fatal(err)
	}
}

func closeFTP() {
	if err := ftpClient.Quit(); err != nil {
		log.Fatal(err)
	}
}

func readFileList() []*ftp.Entry {
	// Do something with the FTP conn
	entries, err := ftpClient.List("/")
	if err != nil {
		log.Fatal(err)
	}
	return entries
}

func initFirestore() {

	opt := option.WithCredentialsFile("ekc-stream-firebase-adminsdk-ses37-27da0fc036.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("firebase app is initialized.")
}

func download(entry *ftp.Entry) Match {
	r, err := ftpClient.Retr(entry.Name)
	defer r.Close()
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadAll(r)

	reader := csv.NewReader(strings.NewReader(string(buf)))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	match := makeMatch(records[0])
	matches = append(matches, &match)
	return match
}

func makeMatch(record []string) Match {
	return Match{
		Shiaijo:         record[0],
		Pool:            record[1],
		Fight:           record[2],
		NumberTareWhite: record[3],
		NameTareWhite:   record[4],
		HansokuWhite:    record[5],
		IpponWhite2:     record[6],
		IpponWhite1:     record[7],
		EnchoOrHikiwake: record[8],
		IpponRed1:       record[9],
		IpponRed2:       record[10],
		HansokuRed:      record[11],
		NumberTareRed:   record[12],
		NameTareRed:     record[13],
	}
}
