package value

import (
	"time"
)

func UnixToTime(unix int64) *time.Time {
	t := time.Unix(unix, 0)
	return &t
}
