package main

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"
)

func ParseMatch(data string) MatchInterface {
	reader := csv.NewReader(strings.NewReader(data))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	check(err)

	if len(records[0]) > 15 {
		return MakeTeamMatch(records[0])
	} else {
		return MakeMatch(records[0])
	}
}

func MakeMatch(record []string) Match {
	poolNumber, _ := strconv.Atoi(record[1])
	return Match{
		Shiaijo:         record[0],
		Pool:            poolNumber,
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
