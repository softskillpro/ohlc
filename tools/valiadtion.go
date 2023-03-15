package tools

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

// ValidateCSVHeader gives csv.Reader and try to fetch header of csv file.
func ValidateCSVHeader(reader *csv.Reader) (bool, error) {

	header, err := reader.Read()
	if err != nil {
		return false, fmt.Errorf("failed to read header: %v", err) // return error if failed to read header
	}

	if len(header) != 6 {
		// return error if length of header is not 6
		return false, fmt.Errorf("header has wrong format: expected=6 fields, actual=%d fields", len(header))
	}

	// return error if header fields don't match expected fields
	if header[0] != "UNIX" || header[1] != "SYMBOL" || header[2] != "OPEN" || header[3] != "HIGH" || header[4] != "LOW" || header[5] != "CLOSE" {
		return false, fmt.Errorf("header has wrong fields: expected=[UNIX,SYMBOL,OPEN,HIGH,LOW,CLOSE], actual=%v", header)
	}

	return true, nil
}

func ValidateOneRow(record []string) error {

	// Parse the string values from the CSV file to the necessary data types.
	if _, err := strconv.Atoi(record[0]); err != nil {
		return errors.New("invalid unix")
	}

	if _, err := strconv.ParseFloat(record[2], 64); err != nil {
		return errors.New("invalid open")
	}

	if _, err := strconv.ParseFloat(record[3], 64); err != nil {
		return errors.New("invalid high")
	}

	if _, err := strconv.ParseFloat(record[4], 64); err != nil {
		return errors.New("invalid low")
	}

	if _, err := strconv.ParseFloat(record[5], 64); err != nil {
		return errors.New("invalid close")
	}

	return nil
}
