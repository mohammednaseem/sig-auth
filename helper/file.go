package helper

import (
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("")
		return err
	}
	if isFileExisting {
		err = os.Remove(fileName)
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}
	}
	return nil
}
