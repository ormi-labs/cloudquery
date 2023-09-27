package client

import (
	"strconv"
	"strings"
	"time"
)

// FmtElapsedTime calculates elapsed time from the arguments and formats it to a human-readable representation.
func FmtElapsedTime(then time.Time, now time.Time, sep string, end bool) string {
	var (
		s = func(x int) string {
			if int(x) == 1 {
				return ""
			}
			return "s"
		}

		elapsed  = int(now.Sub(then).Seconds())
		remained int
		day      = 24 * 60 * 60
		hour     = 60 * 60
		minute   = 60

		parts []string
		text  string
	)

	days := elapsed / day
	remained = elapsed % day
	hours := remained / hour
	remained = remained % hour
	minutes := remained / minute
	remained = remained % minute
	seconds := remained

	if days > 0 {
		parts = append(parts, strconv.Itoa(days)+" day"+s(days))
	}

	if hours > 0 {
		parts = append(parts, strconv.Itoa(hours)+" hour"+s(hours))
	}

	if minutes > 0 {
		parts = append(parts, strconv.Itoa(minutes)+" minute"+s(minutes))
	}

	if seconds > 0 {
		parts = append(parts, strconv.Itoa(seconds)+" second"+s(seconds))
	}

	if end {
		if now.After(then) {
			text = " ago"
		} else {
			text = " after"
		}
	}

	if len(parts) == 0 {
		return "just now"
	}

	return strings.Join(parts, sep) + text
}
