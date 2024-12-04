package utils

import (
	"time"
)

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}
