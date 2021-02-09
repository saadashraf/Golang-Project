package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CliStreamerRecord struct {
	Runner      int64
	Title       string
	Message1    string
	Message2    string
	StreamDelay int64
	RunTimes    int64
}

type CliRunnerRecord struct {
	StreamerRecord []CliStreamerRecord
}

func inputArgs() []CliStreamerRecord {

	var StreamerRecord []CliStreamerRecord
	var streamerRecord CliStreamerRecord
	input := os.Args[1:]
	args := strings.Join(input, " ")
	formattedArgs := strings.Split(args, "\\n")

	for i := 1; i < len(formattedArgs); i++ {
		individualInput := strings.Split(formattedArgs[i], ",")
		runner, _ := strconv.ParseInt(individualInput[0], 10, 64)
		delay, _ := strconv.ParseInt(individualInput[4], 10, 64)
		runTime, _ := strconv.ParseInt(individualInput[5], 10, 64)
		streamerRecord = CliStreamerRecord{
			Runner:      runner,
			Title:       individualInput[1],
			Message1:    individualInput[2],
			Message2:    individualInput[3],
			StreamDelay: delay,
			RunTimes:    runTime,
		}

		StreamerRecord = append(StreamerRecord, streamerRecord)

	}

	return StreamerRecord
}

func print(cliStreamer CliStreamerRecord, presentRunner int64, wg *sync.WaitGroup) {
	for i := 0; i < int(cliStreamer.RunTimes); i++ {
		fmt.Println(cliStreamer.Title + "->" + cliStreamer.Message1 + " " + "instance " + strconv.FormatInt(presentRunner+1, 10))
		fmt.Println(cliStreamer.Title + "->" + cliStreamer.Message2 + " " + "instance " + strconv.FormatInt(presentRunner+1, 10))
		delay := int(cliStreamer.StreamDelay) * 1000
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	wg.Done()
}

func write(cliStreamer CliStreamerRecord, presentRunner int64, wg *sync.WaitGroup, lock *sync.RWMutex) {
	for i := 0; i < int(cliStreamer.RunTimes); i++ {
		f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		output := cliStreamer.Title + "->" + cliStreamer.Message1 + " " + "instance " + strconv.FormatInt(presentRunner+1, 10) + "\n" + cliStreamer.Title + "->" + cliStreamer.Message2 + " " + "instance " + strconv.FormatInt(presentRunner+1, 10) + "\n"

		_, err = f.WriteString(output)

		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	lock.RUnlock()
	wg.Done()
}

func individualStream(cliStreamer CliStreamerRecord) {
	var wg = sync.WaitGroup{}
	var lock = sync.RWMutex{}
	for i := int64(0); i < cliStreamer.Runner; i++ {
		wg.Add(2)
		go print(cliStreamer, i, &wg)
		lock.RLock()
		go write(cliStreamer, i, &wg, &lock)
	}
	wg.Wait()
	//time.Sleep(2000 * time.Millisecond)
}

func (cliRunner CliRunnerRecord) streamAndWrite() {
	for i := 0; i < (len(cliRunner.StreamerRecord)); i++ {
		individualStream(cliRunner.StreamerRecord[i])
	}
}

func main() {

	StreamerRecord := inputArgs()
	fmt.Println(StreamerRecord[0].Runner)
	var runnerRecord CliRunnerRecord
	runnerRecord.StreamerRecord = StreamerRecord
	fmt.Println(runnerRecord.StreamerRecord[0].Message1)
	fmt.Println(len(runnerRecord.StreamerRecord))

	runnerRecord.streamAndWrite()
}
