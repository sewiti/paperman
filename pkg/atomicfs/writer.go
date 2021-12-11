package atomicfs

import (
	"io"
	"os"
	"time"
)

const tempSuffix = ".temp"

type writer struct {
	f    *os.File
	path string
	err  error
}

func NewWriter(name string, perm os.FileMode) (io.WriteCloser, error) {
	const (
		maxWait = 15 * time.Second
		sleep   = 50 * time.Millisecond
	)
	var err error
	start := time.Now()
	for time.Since(start) < maxWait {
		var f *os.File
		f, err = os.OpenFile(name+tempSuffix, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
		if err != nil {
			if os.IsExist(err) {
				time.Sleep(sleep)
				continue
			}
			return nil, err
		}
		return &writer{
			f:    f,
			path: name,
			err:  nil,
		}, nil
	}
	return nil, err
}

func (w *writer) Write(b []byte) (int, error) {
	n, err := w.f.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

func (w *writer) Close() error {
	if w.err != nil {
		_ = w.f.Close()
		_ = os.Remove(w.path + tempSuffix)
		return w.err
	}
	err := w.f.Close()
	if err != nil {
		_ = os.Remove(w.path + tempSuffix)
		return err
	}
	err = os.Rename(w.path+tempSuffix, w.path)
	if err != nil {
		_ = os.Remove(w.path + tempSuffix)
		return err
	}
	return nil
}
