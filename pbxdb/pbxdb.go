package pbxdb

import (
	"3cx-listener/types"
	"3cx-listener/utils"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var conn *sqlx.DB
var once sync.Once

func GetConnection() *sqlx.DB {
	connectionStr := viper.GetString("pbx.database.string")

	once.Do(func() {
		var err error
		if conn, err = sqlx.Open("postgres", connectionStr); err != nil {
			log.Panic(err)
		}
		conn.SetMaxOpenConns(20)
		conn.SetMaxIdleConns(0)
		conn.SetConnMaxLifetime(time.Nanosecond)
	})
	return conn

	//db, err := sqlx.Connect("postgres",  viper.GetString("pbx.database.string"))
	//if err != nil {
	//	utils.CaptureSentryException(err)
	//	log.Fatalln(err)
	//}
	//
	//db.SetMaxOpenConns(16)
	//db.SetMaxIdleConns(16)
	//db.SetConnMaxLifetime(5*time.Minute)
	//database = db
}

//func CloseDatabaseConnection(conn *sqlx.DB) {
//	err := conn.Close()
//	if err != nil {
//		utils.CaptureSentryException(err)
//		log.Fatalln(err)
//	}
//}

func FindParticipantsByRecordingId(recordingId string) (numbers []string) {
	conn := GetConnection()
	rows, err := conn.Queryx(
		"SELECT caller_number FROM recording_participant WHERE fkid_recordings = $1",
		recordingId)

	defer rows.Close()

	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	for rows.Next() {
		var caller_number string
		err := rows.Scan(&caller_number)
		if err != nil {
			utils.CaptureSentryException(err)
			log.Fatalln(err)
		}

		numbers = append(numbers, strings.Replace(caller_number, "Ext.", "", 1))
	}

	return
}

func FindRecordById(recordingId string) (record types.Record) {
	conn := GetConnection()
	rows, err := conn.Queryx(
		"SELECT recording_url, start_time, end_time FROM recordings WHERE id_recording = $1",
		recordingId)

	defer rows.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		var recording_url string
		var start_time pgtype.Timestamp
		var end_time pgtype.Timestamp
		err := rows.Scan(&recording_url, &start_time, &end_time)
		if err != nil {
			utils.CaptureSentryException(err)
			log.Fatalln(err)
		}

		recordInstance := types.Record{
			RecordingUrl: recording_url,
			StartTime: start_time,
			EndTime: end_time,
		}

		record = recordInstance
	}

	return
}