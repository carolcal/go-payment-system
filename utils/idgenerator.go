package utils

import (
	"fmt"
	"time"
)

func GenerateID() string {
	return "pay_" + fmt.Sprintf("%d", time.Now().UnixNano())
}