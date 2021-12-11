package tmux

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"strings"
)

type Session string

const tmux = "tmux"

func ListSessions() ([]Session, error) {
	return ListSessionsContext(context.Background())
}

func ListSessionsContext(ctx context.Context) ([]Session, error) {
	bs, err := exec.Command(tmux, "list-sessions").Output()
	if err != nil {
		return nil, err
	}
	var sessions []Session
	sc := bufio.NewScanner(bytes.NewReader(bs))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ": ", 2)
		if len(parts) != 2 {
			continue
		}
		sessions = append(sessions, Session(parts[0]))
	}
	return sessions, sc.Err()
}
