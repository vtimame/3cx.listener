package types

import (
	"fmt"
	"github.com/jackc/pgtype"
	"strconv"
	"time"
)

type Record struct {
	RecordingUrl string
	StartTime    pgtype.Timestamp
	EndTime      pgtype.Timestamp
}

type RtPayload struct {
	RecordingUrl string
	StartTime    string
	EndTime      string
	IncomingDN string
	OutboundDN string
}

func GetRtPayload(record Record, numbers [] string) (payload RtPayload) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	s := record.StartTime.Time.In(loc)
	e := record.StartTime.Time.In(loc)

	fmt.Println(string(s.Hour()))

	payload = RtPayload{
		RecordingUrl: record.RecordingUrl,
		StartTime:    strconv.Itoa(s.Hour())+":"+strconv.Itoa(s.Minute())+"_"+strconv.Itoa(s.Day())+"."+strconv.Itoa(int(s.Month()))+"."+strconv.Itoa(s.Year()),
		EndTime:      strconv.Itoa(e.Hour())+":"+strconv.Itoa(e.Minute())+"_"+strconv.Itoa(e.Day())+"."+strconv.Itoa(int(e.Month()))+"."+strconv.Itoa(e.Year()),
		IncomingDN:   numbers[0],
		OutboundDN:   numbers[1],
	}

	return
}