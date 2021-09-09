package observing

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ykaseng/hlwm/pkg/logging"
)

type Service interface {
	TagChangeEvent(context.Context) <-chan Event
	TagStatus() string
}

type service struct{}

func NewService() *service {
	return &service{}
}

type Event int

const (
	TagChange Event = iota
)

func (s *service) TagChangeEvent(ctx context.Context) <-chan Event {
	cmd := exec.Command("herbstclient", "--idle", "tag_*")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Logger.Errorf("api: could not create stdout pipe: %v\n", err)
	}

	stdin := bufio.NewScanner(stdout)
	es := make(chan Event)
	go func() {
		defer close(es)
		for stdin.Scan() {
			select {
			case <-ctx.Done():
				return
			case es <- TagChange:
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		logging.Logger.Errorf("api: could not execute command: %v\n", err)
	}

	go cmd.Wait()

	return es
}

func (s *service) TagStatus() string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := generator(ctx)

	numProcessors := runtime.NumCPU()
	processors := make([]<-chan result, numProcessors)
	for i := 0; i < numProcessors; i++ {
		processors[i] = processor(ctx, ts)
	}

	sM := make(StatusMap)
	for r := range take(ctx, fanIn(ctx, processors...), 9) {
		sM[r.Tag] = r.View
	}

	b, err := json.Marshal(sM)
	if err != nil {
		logging.Logger.Errorf("api: could not marshal status map: %v\n", err)
	}

	return string(b)
}

func generator(ctx context.Context) <-chan string {
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

	tags := strings.Split(strings.TrimSpace(stdin.Text()), "\t")
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

func batch(ctx context.Context, rs <-chan result) <-chan StatusMap {
	smS := make(chan StatusMap)
	go func() {
		defer close(smS)
		sM := make(StatusMap)
		for {
			select {
			case <-ctx.Done():
				return
			case r := <-rs:
				sM[r.Tag] = r.View
				if len(sM) == 9 {
					smS <- sM
					sM = make(StatusMap)
				}
			}
		}
	}()

	return smS
}

func take(ctx context.Context, rs <-chan result, num int) <-chan result {
	ts := make(chan result)
	go func() {
		defer close(ts)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case ts <- <-rs:
			}
		}
	}()

	return ts
}
