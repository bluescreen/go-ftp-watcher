package main

import (
	"fmt"
)

type Event struct {
	Type        string
	Description string
	Payload     string
	Date        string
}

func (Event) String() string {
	return ""
}

type MatchInterface interface {
	Table() string
	Description() string
}

type Match struct {
	Shiaijo         string
	Pool            int
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
	Updated         string
}

func (m Match) Table() string {
	return "matches"
}

func (m Match) Description() string {
	return fmt.Sprintf("Pool %d, Fight %s on Shiaijo %s. %s vs %s",
		m.Pool, m.Fight, m.Shiaijo, m.NameTareRed,
		m.NumberTareWhite)
}

type TeamMatch struct {
	Shiaijo             string
	Pool                string
	Fight               string
	NumberTareWhite     string
	NameTareWhite       string
	HansokuWhite        string
	IpponWhite2         string
	IpponWhite1         string
	EnchoOrHikiwake     string
	IpponRed1           string
	IpponRed2           string
	HansokuRed          string
	NumberTareRed       string
	NameTareRed         string
	FightNumber         string
	TeamWhite           string
	WinsWhite           string
	SetWhite            string
	TeamEnchoOrHikiwake string
	SetRed              string
	WinsRed             string
	TeamRed             string
	Updated             string
}

func (m TeamMatch) Table() string {
	return "matches"
}

func (m TeamMatch) Description() string {
	return fmt.Sprintf("Pool %s, Fight %s on Shiaijo %s. %s vs %s",
		m.Pool, m.Fight, m.Shiaijo, m.NameTareRed,
		m.NumberTareWhite)
}
