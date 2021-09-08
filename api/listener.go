package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Listener struct{}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Start() {
	cmd := exec.Command("herbstclient", "--idle", "tag_*")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("[ hlwm ] api: could not create stdout pipe: %v\n", err)
	}

	stdin := bufio.NewScanner(stdout)

	done := make(chan interface{})
	defer close(done)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := func(ctx context.Context, tags []string) <-chan string {
		ts := make(chan string)
		go func() {
			defer close(ts)
			for i := range tags {
				select {
				case <-ctx.Done():
					return
				case ts <- tags[i]:
				}
			}
		}()

		return ts
	}

	type Result struct {
		Status
		Error error
	}

	process := func(s string) Result {
		var v State
		switch string(s[0]) {
		case "#":
			v = Focused
		case ":":
			v = Occupied
		case "!":
			v = Urgent
		case "-":
			v = Viewed
		default:
			v = Empty
		}

		t, err := strconv.Atoi(string(s[1]))
		if err != nil {
			return Result{
				Error: err,
			}
		}

		return Result{
			Status: Status{
				Tag:  t,
				View: v,
			},
		}
	}

	processor := func(ctx context.Context, ts <-chan string) <-chan Result {
		rs := make(chan Result)
		go func() {
			defer close(rs)
			for t := range ts {
				r := process(t)
				select {
				case <-ctx.Done():
					return
				case rs <- r:
				}
			}
		}()

		return rs
	}

	fanIn := func(ctx context.Context, channels ...<-chan Result) <-chan Result {
		var wg sync.WaitGroup
		multiplexedSteam := make(chan Result)

		multiplex := func(c <-chan Result) {
			defer wg.Done()
			for r := range c {
				select {
				case <-ctx.Done():
					return
				case multiplexedSteam <- r:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedSteam)
		}()

		return multiplexedSteam
	}

	go func() {
		for stdin.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				cmd := exec.Command("herbstclient", "tag_status")
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Fatalf("[ hlwm ] api: could not create stdout pipe: %v\n", err)
				}

				if err = cmd.Start(); err != nil {
					log.Fatalf("[ hlwm ] api: could not execute command: %v\n", err)
				}

				stdin := bufio.NewScanner(stdout)
				if !stdin.Scan() {
					log.Fatalf("[ hlwm ] api: could not read status: %v\n", err)
				}

				rt := strings.Split(strings.TrimSpace(stdin.Text()), "\t")
				ts := generator(ctx, rt)

				numProcessors := runtime.NumCPU()
				processors := make([]<-chan Result, numProcessors)
				for i := 0; i < numProcessors; i++ {
					processors[i] = processor(ctx, ts)
				}

				sM := make(StatusMap)
				for r := range fanIn(ctx, processors...) {
					if r.Error != nil {
						log.Fatalf("[ hlwm ] api: could not get tag state: %v\n", r.Error)
						continue
					}

					sM[r.Tag] = r.View
				}

				b, err := json.Marshal(sM)
				if err != nil {
					log.Fatalf("[ hlwm ] api: could not marshal status map: %v\n", err)
				}

				fmt.Println(string(b))
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		log.Fatalf("[ hlwm ] api: could not execute command: %v\n", err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatalf("[ hlwm ] api: could not wait: %v\n", err)
	}
}
