package main

import (
	"3cx-listener/config"
	"3cx-listener/converter"
	"3cx-listener/listener"
	"3cx-listener/pbxdb"
	"3cx-listener/rt"
	"3cx-listener/types"
	"3cx-listener/utils"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)


func main() {
	fmt.Println("Build")
	config.ReadConfig()
	utils.SentryInit()
	for {
		payload := listener.WaitForNotification()
		var recordingId = gjson.Get(payload, "id_recording").String()
		if len(recordingId) > 0 {
			handleNotificationCallback(recordingId)
		}
	}
}

func handleNotificationCallback(recordingId string) {
	if len(recordingId) > 0 {
		log.Println("New record:", recordingId)
		record := pbxdb.FindRecordById(recordingId)
		numbers := pbxdb.FindParticipantsByRecordingId(recordingId)
		payload := types.GetRtPayload(record, numbers)

		_, findIncoming := FindInNumbers("numbers.support", payload.IncomingDN)
		_, findOutbound := FindInNumbers("numbers.support", payload.OutboundDN)

		if findIncoming || findOutbound {
			fmt.Println("Is support call")
			recordPath := converter.ConvertRecord(payload)
			ticket := string(rt.CreateTicket(payload))
			ticketId := gjson.Get(ticket, "id").Str
			rt.SentAttachment(payload, recordPath, ticketId)
			converter.DeleteRecord(recordPath)
		}
	}
}

func FindInNumbers(path string, val string) (int, bool) {
	slice := strings.Split(viper.GetString(path), ",")

	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}