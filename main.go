package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zan8in/masscan"
)

func main() {
	context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("60.10.116.10"),
		masscan.SetParamPorts("80"),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(50),
		masscan.WithContext(context),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		panic(err)
	}

	stdout := scanner.GetStdout()

	stderr := scanner.GetStderr()

	go func() {
		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			fmt.Println(srs.IP, srs.Port)
			scannerResult = append(scannerResult, srs)
		}
	}()

	go func() {
		for stderr.Scan() {
			fmt.Println("err: ", stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult), " PID : ", scanner.GetPid())

}