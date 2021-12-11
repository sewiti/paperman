package instances

import (
	"context"
	"strings"

	"github.com/sewiti/paperman/pkg/tmux"
)

const instancePrefix = "paperman-"

func ListAll(ctx context.Context) ([]tmux.Session, error) {
	sessions, err := tmux.ListSessionsContext(ctx)
	if err != nil {
		return nil, err
	}
	var instances []tmux.Session
	for _, session := range sessions {
		if !strings.HasPrefix(string(session), instancePrefix) {
			continue
		}
		instances = append(instances, session)
	}
	return instances, nil
}
