package listener

import (
	"3cx-listener/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"time"
)

func createListener() (listener *pq.Listener) {
	_, err := sql.Open("postgres", viper.GetString("pbx.database.string"))
	if err != nil {
		panic(err)
	}

	listener = pq.NewListener(viper.GetString("pbx.database.string"), 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("recordings")
	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	log.Println("Start monitoring PostgreSQL...")
	return
}

func reportProblem(ev pq.ListenerEventType, err error) {
	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}
}

func WaitForNotification() (payload string) {
	l := createListener()
	for {
		select {
		case n := <-l.Notify:
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				utils.CaptureSentryException(err)
				log.Fatalln(err)
				return
			}
			payload = string(prettyJSON.Bytes())
			return
		case <-time.After(120 * time.Second):
			log.Println("Received no events for 120 seconds, checking connection")
			go func() {
				err := l.Ping()
				if err != nil {
					utils.CaptureSentryException(err)
					log.Fatalln(err)
				}
			}()
			return
		}
	}
}
