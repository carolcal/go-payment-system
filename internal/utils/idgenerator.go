package utils

import (
	"fmt"
	"time"
)

func GenerateID(prefix string) string {
	return prefix + "_" + fmt.Sprintf("%d", time.Now().UnixNano())
}