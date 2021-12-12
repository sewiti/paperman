package screen

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"unicode"
)

type Screen string

func List(match string) ([]Screen, error) {
	return ListContext(context.Background(), match)
}

func ListContext(ctx context.Context, match string) ([]Screen, error) {
	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx, "screen", "-ls", match)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		if bytes.HasPrefix(buf.Bytes(), []byte("No Sockets found")) {
			return nil, nil
		}
		return nil, err
	}
	return parseScreens(&buf)
}

func parseScreens(r io.Reader) ([]Screen, error) {
	var screens []Screen
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if !bytes.HasPrefix(sc.Bytes(), []byte("\t")) {
			continue
		}
		start := bytes.IndexByte(sc.Bytes(), '.')
		if start < 0 {
			continue
		}
		start++ // period
		end := bytes.LastIndexByte(sc.Bytes(), '(')
		if end < 0 {
			continue
		}
		end = bytes.LastIndexByte(sc.Bytes()[:end], '(')
		if end < 0 {
			continue
		}
		end = bytes.LastIndexFunc(sc.Bytes()[:end], notSpace)
		if end < start {
			continue
		}
		screen := sc.Text()[start : end+1]
		screens = append(screens, Screen(screen))
	}
	if sc.Err() != nil {
		return nil, sc.Err()
	}
	return screens, nil
}

func notSpace(r rune) bool {
	return !unicode.IsSpace(r)
}
