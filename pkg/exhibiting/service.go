package exhibiting

import (
	"context"
	"os/exec"
	"time"

	"github.com/ykaseng/hlwm/pkg/logging"
)

type Service interface {
	ShowWidget(string)
	HideWidget(string)
	FlashWidget(string, time.Duration)
}

type service struct{}

func NewService() *service { return &service{} }

func (s *service) ShowWidget(w string) { showWidget(context.Background(), w) }

func (s *service) HideWidget(w string) { hideWidget(context.Background(), w) }

func (s *service) FlashWidget(w string, d time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := time.NewTimer(d)
	if err := showWidget(ctx, w); err != nil {
		logging.Logger.Errorf("exhibiting: could not show widget: %v\n", err)
		cancel()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			hideWidget(ctx, w)
			return
		}
	}
}

func showWidget(ctx context.Context, w string) error {
	cmd := exec.Command("eww", "open", w)
	return cmd.Start()
}

func hideWidget(ctx context.Context, w string) error {
	cmd := exec.Command("eww", "close", w)
	return cmd.Start()
}
