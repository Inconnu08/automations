package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"os"
	"strconv"
	"sync"

	"github.com/gocarina/gocsv"
)

// "Title,Message 1,Message 2,Stream Delay,Run Times\nCLI Invoker Name,First Message,Second Msg,2,10"
type CliStreamerRecord struct {
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

type CliRunnerRecord struct {
	// How many streamer will run.
	Run         string `csv:"Run"`
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

func (cliRunnerRecord CliRunnerRecord) CliStreamerRecord() CliStreamerRecord {
	return CliStreamerRecord{
		Title:       cliRunnerRecord.Title,
		Message1:    cliRunnerRecord.Message1,
		Message2:    cliRunnerRecord.Message2,
		StreamDelay: cliRunnerRecord.StreamDelay,
		RunTimes:    cliRunnerRecord.RunTimes,
	}
}

func (cliRunnerRecord CliRunnerRecord) CliStreamerRecordCsv() string {
	cliStreamerRecords := []CliStreamerRecord{cliRunnerRecord.CliStreamerRecord()}

	out, err := gocsv.MarshalString(cliStreamerRecords)

	if err != nil {
		panic(err)
	}

	return out
}

func Csv(cliRunners *[]CliRunnerRecord) string {
	out, err := gocsv.MarshalString(cliRunners)

	if err != nil {
		panic(err)
	}

	return out
}

func main() {
	args := "Run,Title,Message 1,Message 2,Stream Delay,Run Times\n2,CLI Invoke1,First Msg 1,Second Msg 2,2,500\n2,CLI Invoke2,First Msg 1,Second Msg 2,2,500"
	var fileMutex sync.Mutex
	//gocsv.UnmarshalString(
	//	args,
	//	&cliStreamers)
	//
	////fmt.Print(cliStreamers[1].StreamDelay)
	////fmt.Println("\n")
	////fmt.Print(gocsv.MarshalString(cliStreamers))
	//fmt.Println(time.Duration(cliStreamers[0].StreamDelay) * time.Second)
	var cliRunners []CliRunnerRecord
	gocsv.UnmarshalString(
		args,
		&cliRunners)

	fmt.Println(cliRunners[0].Run)
	cmd := exec.Command("CLIStreamer.exe", "Title,Message 1,Message 2,Stream Delay,Run Times\nCLI Invoker Name,First Message,Second Msg,2,10")

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		writeToFile(m, &fileMutex)
	}

	scanner = bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		writeToFile(m, &fileMutex)
	}

	cmd.Wait()
}

func writeToFile(text string, fileMutex *sync.Mutex) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	f, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Println(err)
	}
}

// run runs the CLI task
func run(records []CliRunnerRecord) {
	println(records)
	//var wg sync.WaitGroup
	//
	//for _, record := range records {
	//	count := 0
	//	for count < record.RunTimes {
	//		consolePrint(record.Title, record.Message1, &wg)
	//		time.Sleep(time.Duration(record.StreamDelay) * time.Second)
	//		consolePrint(record.Title, record.Message2, &wg)
	//		count++
	//	}
	//	count = 0
	//}
	//wg.Wait()
}

// consolePrint prints csv messages concurrently
func consolePrint(title, message string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("%s->%s\n", title, message)
	}()
}

// cli entry
func execCli() (string, error) {
	if len(os.Args) != 2 {
		return "", errors.New("no CSV string provided")
	}

	flag.Parse()

	return os.Args[1], nil
}

// getCSVRecords constructs the CSV records and returns `CliStreamerRecord`s
func getCSVRecords(csv string) ([]CliRunnerRecord, error) {
	csv, _ = strconv.Unquote(`"` + csv + `"`)

	var cliRunner []CliRunnerRecord
	if err := gocsv.UnmarshalString(csv, &cliRunner); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return cliRunner, nil
}
