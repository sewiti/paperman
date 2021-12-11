package atomicfs

import "os"

func WriteFile(name string, data []byte, perm os.FileMode) error {
	f, err := NewWriter(name, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
