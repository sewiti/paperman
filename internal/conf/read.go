package conf

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const commentPrefix = "#"

func Read(path string) (Values, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f)
}

func parse(r io.Reader) (Values, error) {
	values := make(Values)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if strings.HasPrefix(sc.Text(), commentPrefix) {
			continue
		}
		parts := strings.SplitN(sc.Text(), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		values.Add(key, value)
	}
	return values, sc.Err()
}
