package exhibiting

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ykaseng/hlwm/pkg/logging"
)

type Service interface {
	ShowWidget(string)
	HideWidget(string)
	FlashWidget(context.Context, string, time.Duration) chan<- interface{}
}

type service struct{}

func NewService() *service { return &service{} }

func (s *service) ShowWidget(w string) { showWidget(context.Background(), w) }

func (s *service) HideWidget(w string) { hideWidget(context.Background(), w) }

func (s *service) FlashWidget(ctx context.Context, w string, d time.Duration) chan<- interface{} {
	t := time.NewTimer(d)
	fs := make(chan interface{})
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-fs:
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				if err := showWidget(ctx, w); err != nil {
					logging.Logger.Errorf("exhibiting: could not show widget: %v\n", err)
					cancel()
				}

				t = time.NewTimer(d)
			case <-t.C:
				hideWidget(ctx, w)
			}
		}
	}()

	return fs
}

func showWidget(ctx context.Context, w string) error {
	cmd := exec.Command("eww", "update", fmt.Sprintf("%s=true", w))
	return cmd.Start()
}

func hideWidget(ctx context.Context, w string) error {
	cmd := exec.Command("eww", "update", fmt.Sprintf("%s=false", w))
	return cmd.Start()
}
