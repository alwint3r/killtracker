package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func grep(stdin io.Reader, pattern string) (bytes.Buffer, error) {
	cmd := exec.Command("grep", pattern)
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stdin = stdin
	err := cmd.Run()

	return out, err
}

func getProcess() (bytes.Buffer, error) {
	cmd := exec.Command("ps", "-e")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()

	return out, err
}

func filterEmptyString(arrayOfString []string) []string {
	filtered := make([]string, 0)
	for _, item := range arrayOfString {
		trimmed := strings.Trim(item, "\r\n")
		trimmed = strings.TrimSpace(trimmed)
		if len(trimmed) > 1 {
			filtered = append(filtered, trimmed)
		}
	}

	return filtered
}

func killProcess(pid string) error {
	cmd := exec.Command("kill", "-15", pid)

	return cmd.Run()
}

func killProcesses(pids []string) error {
	currentPid := strconv.Itoa(os.Getpid())
	for _, pid := range pids {
		if pid == currentPid {
			log.Printf("Skipping %s as it's my own PID\n", pid)
		} else {
			log.Println("Prepare to kill", pid)
			err := killProcess(pid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	log.Println("Starting...")

	ps, err := getProcess()
	if err != nil {
		log.Fatalln(err)
	}

	psBytes := ps.Bytes()

	grepped, err := grep(bytes.NewReader(psBytes), "tracker")
	if err != nil {
		log.Fatalln(err)
	}

	lines := filterEmptyString(strings.Split(grepped.String(), "\n"))
	pids := make([]string, len(lines))

	for idx, line := range lines {
		cols := filterEmptyString(strings.Split(line, " "))
		pids[idx] = cols[0]
	}

	if err := killProcesses(pids); err != nil {
		log.Fatalln(err)
	}
}
