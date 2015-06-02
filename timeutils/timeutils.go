package timeutils

import (
	"time"
)

func UtcNow() time.Time {
	return time.Now().UTC()
}
