package listening

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ykaseng/hlwm/pkg/logging"
)

type Service interface {
	Start()
}

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Start() {
	cmd := exec.Command("herbstclient", "--idle", "tag_*")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Logger.Errorf("api: could not create stdout pipe: %v\n", err)
	}

	stdin := bufio.NewScanner(stdout)

	done := make(chan interface{})
	defer close(done)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for stdin.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				cmd := exec.Command("herbstclient", "tag_status")
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					logging.Logger.Errorf("api: could not create stdout pipe: %v\n", err)
				}

				if err = cmd.Start(); err != nil {
					logging.Logger.Errorf("api: could not execute command: %v\n", err)
				}

				stdin := bufio.NewScanner(stdout)
				if !stdin.Scan() {
					logging.Logger.Errorf("api: could not read status: %v\n", err)
				}

				rt := strings.Split(strings.TrimSpace(stdin.Text()), "\t")
				ts := generator(ctx, rt)

				numProcessors := runtime.NumCPU()
				processors := make([]<-chan result, numProcessors)
				for i := 0; i < numProcessors; i++ {
					processors[i] = processor(ctx, ts)
				}

				sM := make(StatusMap)
				for r := range fanIn(ctx, processors...) {
					if r.Error != nil {
						logging.Logger.Errorf("api: could not get tag state: %v\n", r.Error)
						continue
					}

					sM[r.Tag] = r.View
				}

				b, err := json.Marshal(sM)
				if err != nil {
					logging.Logger.Errorf("api: could not marshal status map: %v\n", err)
				}

				fmt.Println(string(b))
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		logging.Logger.Errorf("api: could not execute command: %v\n", err)
	}

	if err = cmd.Wait(); err != nil {
		logging.Logger.Errorf("api: could not wait: %v\n", err)
	}
}

func generator(ctx context.Context, tags []string) <-chan string {
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

func processor(ctx context.Context, ts <-chan string) <-chan result {
	process := func(s string) result {
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
			return result{
				Error: err,
			}
		}

		return result{
			Status: Status{
				Tag:  t,
				View: v,
			},
		}
	}

	rs := make(chan result)
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

func fanIn(ctx context.Context, channels ...<-chan result) <-chan result {
	var wg sync.WaitGroup
	multiplexedSteam := make(chan result)

	multiplex := func(c <-chan result) {
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
