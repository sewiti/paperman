package conf

import (
	"fmt"
	"sort"

	"github.com/sewiti/paperman/internal/atomicfs"
)

func Write(path string, values Values) error {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	f, err := atomicfs.NewWriter(path, 0640)
	if err != nil {
		return err
	}
	for _, k := range keys {
		for _, v := range values[k] {
			_, err = fmt.Fprintf(f, "%s=%s\n", k, v)
			if err != nil {
				_ = f.Close()
				return err
			}
		}
	}
	return f.Close()
}
