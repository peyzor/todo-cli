package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/mergestat/timediff"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	IsDoneNo  string = "N"
	IsDoneYes string = "Y"
)

const CSVStorageFilename = "storage.csv"

const TimeLayout = "2006-01-02 15:04:05.999999 -0700"

func GetOrCreateCSVStorage() (*os.File, error) {
	var f *os.File
	filename := CSVStorageFilename
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

func DeleteCSVRecord(f io.ReadWriter, ID int) error {
	rows, err := GetRows(f)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return errors.New("malformed data: header not found")
	}

	header := rows[0]

	var newRows [][]string
	for _, row := range rows[1:] {
		rowID, err := getRowID(row, header)
		if err != nil {
			return err
		}

		if ID == rowID {
			continue
		}

		newRows = append(newRows, row)
	}

	tmpFile, err := os.CreateTemp("", "temp_csv")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for _, newRow := range newRows {
		err = writer.Write(newRow)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), CSVStorageFilename)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCSVRecord(f io.ReadWriter, ID int) error {
	rows, err := GetRows(f)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return errors.New("malformed data: header not found")
	}

	header := rows[0]

	var newRows [][]string
	for _, row := range rows[1:] {
		rowID, err := getRowID(row, header)
		if err != nil {
			return err
		}

		if ID == rowID {
			updatedRow, err := markRecordCompleted(row, header)
			if err != nil {
				return err
			}
			newRows = append(newRows, updatedRow)
			continue
		}

		newRows = append(newRows, row)
	}

	tmpFile, err := os.CreateTemp("", "temp_csv")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for _, newRow := range newRows {
		err = writer.Write(newRow)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), CSVStorageFilename)
	if err != nil {
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

func replaceWithHumanReadableTime(rows [][]string) ([][]string, error) {
	if len(rows) == 0 {
		return nil, errors.New("malformed data: header not found")
	}

	header := rows[0]

	var newRows [][]string
	newRows = append(newRows, header)

	for _, row := range rows[1:] {
		createdTime, err := getCreatedTime(row, header)
		if err != nil {
			return nil, err
		}

		for i, h := range header {
			if h == "Created" {
				row[i] = timediff.TimeDiff(createdTime)
			}
		}
		newRows = append(newRows, row)
	}

	return newRows, nil
}

func GetRowsTabular(f io.Reader, w io.Writer) error {
	writer := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	rows, err := GetRows(f)
	if err != nil {
		return err
	}
	rows, err = replaceWithHumanReadableTime(rows)
	if err != nil {
		return err
	}

	for _, row := range rows {
		_, err = fmt.Fprintln(writer, strings.Join(row, "\t"))
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
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

func getRowID(row []string, header []string) (int, error) {
	var rowID int
	if !slices.Contains(header, "ID") {
		return rowID, errors.New("ID header not found")
	}

	rowMapped := make(map[string]string)
	for i, h := range header {
		rowMapped[h] = row[i]
	}

	rowIDStr := rowMapped["ID"]
	rowID, err := strconv.Atoi(rowIDStr)
	if err != nil {
		return rowID, fmt.Errorf("malformed data: %v", err)
	}

	return rowID, nil
}

func getCreatedTime(row []string, header []string) (time.Time, error) {
	var createdTime time.Time
	if !slices.Contains(header, "Created") {
		return createdTime, errors.New("created header not found")
	}

	rowMapped := make(map[string]string)
	for i, h := range header {
		rowMapped[h] = row[i]
	}

	createdStr := rowMapped["Created"]
	createdTime, err := time.Parse(TimeLayout, createdStr)
	if err != nil {
		return createdTime, err
	}

	return createdTime, nil
}

func markRecordCompleted(row, header []string) ([]string, error) {
	var updatedRecord []string
	for i, h := range header {
		if h == "Done" {
			updatedRecord = append(updatedRecord, IsDoneYes)
			continue
		}
		updatedRecord = append(updatedRecord, row[i])
	}

	return updatedRecord, nil
}
