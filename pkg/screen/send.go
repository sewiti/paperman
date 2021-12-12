package screen

import (
	"context"
	"fmt"
	"os/exec"
)

func (s Screen) SendStuff(stuff string) error {
	return s.SendStuffContext(context.Background(), stuff)
}

func (s Screen) SendStuffContext(ctx context.Context, stuff string) error {
	err := exec.CommandContext(ctx, "screen", "-S", string(s), "-p", "0", "-X", "stuff", stuff).Run()
	if err != nil {
		return fmt.Errorf("send-stuff: %w", err)
	}
	return nil
}
