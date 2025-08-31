package infrastructure

import (
	"time"
)

// timeToUnix converte time.Time para timestamp Unix
func timeToUnix(t time.Time) int64 {
	return t.Unix()
}

// timeFromUnix converte timestamp Unix para time.Time
func timeFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// timeNow retorna o timestamp atual
func timeNow() time.Time {
	return time.Now()
}
