package utils

import "time"

func UnixToDuration(expiredAt int64) int64 {
	return int64(time.Unix(expiredAt, 0).Second())
}
