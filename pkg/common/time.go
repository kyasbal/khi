package common

import "time"

func ParseTime(input string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, input)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.RFC3339, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
