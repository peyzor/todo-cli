package storage

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

func GetOrCreateCSVStorage() (*os.File, error) {
	var f *os.File
	filename := "storage.csv"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(filename)
		if err != nil {
			return nil, err
		}

		headers := []string{"ID", "Todo", "Created"}
		err = AddNewCSVRecord(f, headers)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func AddNewCSVRecord(f io.Writer, record []string) error {
	writer := csv.NewWriter(f)
	err := writer.Write(record)
	if err != nil {
		return err
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return err
	}

	return nil
}
