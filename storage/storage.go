package storage

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	IsDoneNo  string = "N"
	IsDoneYes string = "Y"
)

func GetOrCreateCSVStorage() (*os.File, error) {
	var f *os.File
	filename := "storage.csv"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(filename)
		if err != nil {
			return nil, err
		}

		header := []string{"ID", "Todo", "Created", "Done"}
		err = AddNewCSVRecord(f, header)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0644)
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

func GetNextID(f io.Reader) (int, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}

	reader := csv.NewReader(bytes.NewReader(data))
	header, err := reader.Read()
	if err != nil {
		return 0, err
	}

	var rows []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		row := make(map[string]string)
		for i, h := range header {
			row[h] = record[i]
		}

		rows = append(rows, row)
	}

	maxID := 0
	for _, row := range rows {
		IDStr, ok := row["ID"]
		if !ok {
			return 0, errors.New("malformed data: ID value does not exist")
		}

		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			return 0, fmt.Errorf("malformed data: could not convert %s to integer", IDStr)
		}

		if ID > maxID {
			maxID = ID
		}
	}

	return maxID + 1, nil
}
