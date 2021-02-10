package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

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

func main() {
	var cliStreamers []CliStreamerRecord

	csvString, err := exec()
	if err != nil {
		log.Panicln(err)
	}

	if cliStreamers, err = getCSVRecords(csvString); err != nil {
		log.Panicln(err)
	}

	run(cliStreamers)
}

// run runs the CLI task
func run(records []CliStreamerRecord) {
	var wg sync.WaitGroup

	for _, record := range records {
		count := 0
		for count < record.RunTimes {
			consolePrint(record.Title, record.Message1, &wg)
			time.Sleep(time.Duration(record.StreamDelay) * time.Second)
			consolePrint(record.Title, record.Message2, &wg)
			count++
		}
		count = 0
	}
	wg.Wait()
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
func exec() (string, error) {
	if len(os.Args) != 2 {
		return "", errors.New("no CSV string provided")
	}

	flag.Parse()

	return os.Args[1], nil
}

// getCSVRecords constructs the CSV records and returns `CliStreamerRecord`s
func getCSVRecords(csv string) ([]CliStreamerRecord, error) {
	csv, _ = strconv.Unquote(`"` + csv + `"`)

	var cliStreamers []CliStreamerRecord
	if err := gocsv.UnmarshalString(csv, &cliStreamers); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return cliStreamers, nil
}
