package main

import "time"

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
	Updated         string
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

func MakeMatch(record []string) Match {
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
		Updated:         time.Now().Format(dateLayout),
	}
}

func MakeTeamMatch(record []string) TeamMatch {
	return TeamMatch{
		Shiaijo:             record[0],
		Pool:                record[1],
		Fight:               record[2],
		NumberTareWhite:     record[3],
		NameTareWhite:       record[4],
		HansokuWhite:        record[5],
		IpponWhite2:         record[6],
		IpponWhite1:         record[7],
		EnchoOrHikiwake:     record[8],
		IpponRed1:           record[9],
		IpponRed2:           record[10],
		HansokuRed:          record[11],
		NumberTareRed:       record[12],
		NameTareRed:         record[13],
		FightNumber:         record[14],
		TeamWhite:           record[15],
		WinsWhite:           record[16],
		SetWhite:            record[17],
		TeamEnchoOrHikiwake: record[18],
		SetRed:              record[19],
		WinsRed:             record[20],
		TeamRed:             record[21],
		Updated:             time.Now().Format(dateLayout),
	}
}
