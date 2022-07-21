package helper

import (
	"fmt"
	"os"
)

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func DeleteFileIfExists(fileName string) error {
	isFileExisting, err := fileExists(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if isFileExisting {
		err = os.Remove(fileName)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
