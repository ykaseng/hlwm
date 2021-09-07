package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Listener struct {
	quit chan bool
}

func NewListener(quit chan bool) *Listener {
	return &Listener{quit}
}

func (l *Listener) Start() {
	cmd := exec.Command("herbstclient", "--idle", "tag_*")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("[ hlwm ] api: could not create stdout pipe: %v\n", err)
	}

	stdin := bufio.NewScanner(stdout)
	go func() {
		for stdin.Scan() {
			s, err := getStatus()
			if err != nil {
				log.Fatalf("[ hlwm ] api: could not get status: %v\n", err)
			}

			b, err := json.Marshal(s)
			if err != nil {
				log.Fatalf("[ hlwm ] api: could not marshal JSON: %v\n", err)
			}

			fmt.Println(string(b))
		}
	}()

	if err = cmd.Start(); err != nil {
		log.Fatalf("[ hlwm ] api: could not execute command: %v\n", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatalf("[ hlwm ] api: could not wait: %v\n", err)
	}
}

func getStatus() (Status, error) {
	cmd := exec.Command("herbstclient", "tag_status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Status{}, err
	}

	if err = cmd.Start(); err != nil {
		return Status{}, err
	}

	stdin := bufio.NewScanner(stdout)
	if !stdin.Scan() {
		return Status{}, fmt.Errorf("could not read status")
	}

	ts := strings.Split(strings.TrimSpace(stdin.Text()), "\t")

	tags := make([]int, len(ts))
	views := make([]State, len(ts))
	for i := range ts {
		arr := strings.Split(ts[i], "")
		switch arr[0] {
		case "#":
			views[i] = Focused
		case ":":
			views[i] = Occupied
		case "!":
			views[i] = Urgent
		case "-":
			views[i] = Viewed
		default:
			views[i] = Empty
		}

		tag, err := strconv.Atoi(arr[1])
		if err != nil {
			return Status{}, err
		}

		tags[i] = tag
	}

	return Status{
		Tags:  tags,
		Views: views,
	}, nil
}
