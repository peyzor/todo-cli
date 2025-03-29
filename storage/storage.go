package storage

import (
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

func GetRows(f io.Reader) ([][]string, error) {
	reader := csv.NewReader(f)

	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		rows = append(rows, record)
	}

	return rows, nil
}

func GetRowsMapped(f io.Reader) ([]map[string]string, error) {
	rows, err := GetRows(f)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, errors.New("malformed data: header not found")
	}

	header := rows[0]

	var rowsMapped []map[string]string
	for _, row := range rows[1:] {
		rowMapped := make(map[string]string)
		for i, h := range header {
			rowMapped[h] = row[i]
		}

		rowsMapped = append(rowsMapped, rowMapped)
	}

	return rowsMapped, nil
}

func GetNextID(f io.Reader) (int, error) {
	maxID := 0

	rows, err := GetRowsMapped(f)
	if err != nil {
		return maxID, err
	}

	for _, row := range rows {
		IDStr, ok := row["ID"]
		if !ok {
			return maxID, errors.New("malformed data: ID value does not exist")
		}

		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			return maxID, fmt.Errorf("malformed data: could not convert %s to integer", IDStr)
		}

		if ID > maxID {
			maxID = ID
		}
	}

	return maxID + 1, nil
}
