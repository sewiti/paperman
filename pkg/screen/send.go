package screen

import (
	"context"
	"os/exec"
)

func (s Screen) SendStuff(stuff string) error {
	return s.SendStuffContext(context.Background(), stuff)
}

func (s Screen) SendStuffContext(ctx context.Context, stuff string) error {
	return exec.CommandContext(ctx, "screen", "-S", string(s), "-p", "0", "-X", "stuff", stuff).Run()
}
