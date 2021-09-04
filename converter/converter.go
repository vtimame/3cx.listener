package converter

import (
	"3cx-listener/types"
	"3cx-listener/utils"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
)

const ShellToUse = "bash"
func shellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func DeleteRecord(path string) {
	log.Println("Deleting temp record: " + path)
	err := os.Remove(path)
	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}
}

func ConvertRecord(payload types.RtPayload) (outPath string) {
	dir, err := os.Getwd()
	if err != nil {
		utils.CaptureSentryException(err)
		log.Fatalln(err)
	}

	var url string
	debug := viper.GetBool("debug")
	recordsPath := viper.GetString("pbx.records.path")

	if debug {
		url = dir + "/test/test.wav"
	} else {
		url = recordsPath + "/" + payload.RecordingUrl
	}

	log.Println("WAV path: " + url)

	fileName := "Call_from_"+payload.OutboundDN+"_to_"+payload.IncomingDN+"_["+payload.StartTime+"].mp3"
	outPath = dir + "/temp/"+fileName

	log.Println("Convert to: " + fileName)
	log.Println("Converted path: " + outPath)

	//wavFile, _ := os.OpenFile(url, os.O_RDONLY, 0555)
	//mp3File, _ := os.OpenFile(outPath, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0755)
	//defer mp3File.Close()
	//
	//wavHdr, err := lame.ReadWavHeader(wavFile)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//	os.Exit(1)
	//}
	//
	//wr, _ := lame.NewWriter(mp3File)
	//wr.EncodeOptions = wavHdr.ToEncodeOptions()
	//io.Copy(wr, wavFile)
	//wr.Close()

	err, out, errOut := shellOut(`lame --preset insane "`+ url +`" ` +  outPath )
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	//if len(errOut) > 0 {
	//	utils.CaptureSentryException(errors.New(errOut))
	//	log.Fatalln("Lame error: " + errOut)
	//}

	fmt.Println("Lame out: " + out, errOut)
	return
}
