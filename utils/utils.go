package utils

import (
	"strings"
	"unicode/utf8"
)

func RightPaddedString(value string, maxValueLength, padAmount int) string {
	valueLength := utf8.RuneCountInString(value)
	if maxValueLength-padAmount >= valueLength {
		return strings.Repeat(" ", padAmount) + value + strings.Repeat(" ", maxValueLength-valueLength-padAmount)
	} else if maxValueLength-padAmount < valueLength {
		if maxValueLength-padAmount-3 < 0 {
			return "..."
		}

		tmp := strings.Trim(value[0:maxValueLength-padAmount-3], " ") + "..."
		tmpLength := utf8.RuneCountInString(tmp)
		return strings.Repeat(" ", padAmount) + tmp + strings.Repeat(" ", maxValueLength-tmpLength-padAmount)
	}

	return value
}

func LeftPaddedString(value string, maxValueLength, padAmount int) string {
	valueLength := utf8.RuneCountInString(value)
	if maxValueLength-padAmount >= valueLength {
		return strings.Repeat(" ", maxValueLength-valueLength-padAmount) + value + strings.Repeat(" ", padAmount)
	} else if maxValueLength-padAmount < valueLength {
		if maxValueLength-padAmount-3 < 0 {
			return "..."
		}

		tmp := strings.Trim(value[0:maxValueLength-padAmount-3], " ") + "..."
		tmpLength := utf8.RuneCountInString(tmp)
		return strings.Repeat(" ", maxValueLength-tmpLength-padAmount) + tmp + strings.Repeat(" ", padAmount)
	}

	return value
}
