package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var schema = []string{
	"Id",
	"Date",
	"Transaction",
}

type validRecord struct {
	Id    int
	Value int64
	Date  time.Time
}

func (vr *validRecord) validateRecord(record []string) error {
	date := record[1]
	txn := record[2]
	var parsedDate time.Time

	for i := 0; i < 3; i++ {
		//Validating non blank fields
		if record[i] == "" {
			return fmt.Errorf("%s should not be null", schema[i])
		}
	}
	// need to validate id
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return err
	}

	txnValue, err := validateTransactionValue(txn)
	if err != nil {
		return err
	}

	parsedDate, err = validateDate(date)
	if err != nil {
		return err
	}

	vr.Id = id
	vr.Date = parsedDate
	vr.Value = txnValue
	return nil
}

func validateTransactionValue(value string) (int64, error) {
	// Check if value matches the pattern: sign (optional), digits, dot (optional), 2 digits (optional)
	matched, err := regexp.MatchString(`^[+-]?\d+(\.\d{1,2})?$`, value)
	if err != nil {
		return 0, err
	}

	if !matched {
		return 0, fmt.Errorf("invalid transaction value: %s", value)
	}

	// Convert string value to float64 because of CSV format
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	// Converting float value to int, from dollars into cents, so I can process it later with no rounding errors.
	intValue := int64(amount * 100)

	return intValue, nil
}

func validateDate(date string) (time.Time, error) {
	// Get the current year
	currentYear := time.Now().Year()

	// Parse the date string into a time.Time value
	t, err := time.Parse("1/2", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %s", date)
	}

	// Set the year to the current year
	t = t.AddDate(currentYear-t.Year(), 0, 0)

	return t, nil
}
