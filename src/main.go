package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/jlaffaye/ftp"
)

var firebaseConfig *firebase.Config = &firebase.Config{ProjectID: firestoreProjectId}
var firestoreClient *firestore.Client
var ctx = context.Background()
var ftpClient *ftp.ServerConn

var matches []*MatchInterface
var entries []*ftp.Entry

func main() {
	argsWithProg := os.Args
	forceUpdate := len(argsWithProg) > 1 && argsWithProg[1] == "-f"

	fmt.Println("Starting ...", forceUpdate)
	initEventLog()
	for {
		run(forceUpdate)
		time.Sleep(time.Duration(pollInterval) * time.Second)
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

func toCharStrConst(i int) string {
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return abc[i-1 : i]
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
