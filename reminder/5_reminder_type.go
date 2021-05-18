package reminder

import "strings"

type reminderType string

const OIL reminderType = "OIL"
const TUNE_UP reminderType = "TUNEUP"

func ReminderType(x string) reminderType {
	switch strings.ToLower(x) {
	case "oil":
		return OIL
	case "tune_up":
		return TUNE_UP
	}
	return reminderType("UNKNOWN")
}
