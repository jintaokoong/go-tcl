package utils

import (
	"fmt"
	"os"
	"time"
)

func GetFile(prefix string) (*os.File, error) {
	ct := time.Now()
	fn := fmt.Sprintf("%s_%s.log", prefix, ct.Format("20060102"))
	return os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}
