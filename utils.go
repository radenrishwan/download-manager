package downloadmanager

import (
	"os"
)

func CreateDummyFile(filename string, size int64) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	err = f.Truncate(size)
	if err != nil {
		return err
	}

	defer f.Close()
	return nil
}
