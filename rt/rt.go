package rt

import (
	"3cx-listener/types"
	"3cx-listener/utils"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateTicket(p types.RtPayload) (body []byte) {
	var debug bool = viper.GetBool("debug")
	var rtUrl string = viper.GetString("rt.api.host")
	var rtToken string = viper.GetString("rt.api.token")
	var rtTicketsQueue = viper.GetString("rt.api.tickets_queue")
	var subject string = "Call from " + p.OutboundDN + " to " + p.IncomingDN + " [" + p.StartTime + "]"
	if debug {
		subject = "TEST_" + subject
	}

	fmt.Println("Creating ticket: " + subject)

	url := rtUrl + "/ticket?token=" + rtToken
	var jsonStr = []byte(`{ "Subject": "`+subject+`", "Queue": "`+rtTicketsQueue+`" }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-type", "application/json")

	client := getHttpClient()

	resp, err := client.Do(req)

	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)
	body, _ = ioutil.ReadAll(resp.Body)
	log.Println("Response body:", string(body))
	return
}

func SentAttachment(payload types.RtPayload, recordPath string, ticketId string) {
	f, err := os.Open(recordPath)
	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	client := getHttpClient()
	var rtUrl string = viper.GetString("rt.api.host")
	var rtToken string = viper.GetString("rt.api.token")
	var subject string = viper.GetString("rt.api.comment_subject")
	var recordName = strings.Split(recordPath, "/")

	url := rtUrl + "/ticket/"+ticketId+"/comment?token=" + rtToken
	var jsonStr = []byte(`{ "Subject": "`+subject+`", "Attachments": [{ "FileName": "`+recordName[len(recordName) - 1]+`", "FileType": "audio/mpeg", "FileContent": "`+encoded+`" }] }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-type", "application/json")

	log.Println("Sending attachment in ticket: " + ticketId)

	resp, err := client.Do(req)

	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response body:", string(body))
}

func getHttpClient() (client *http.Client) {
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return
}