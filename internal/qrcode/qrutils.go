package qrcode

import (
	"fmt"
	"time"
    "unicode"

    "golang.org/x/text/runes"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

func TLV(id string, value string) string {
	l := len([]byte(value)) // bytes length (UTF-8)
	return fmt.Sprintf("%s%02d%s", id, l, value)
}

func generateTXID() string {
	now := time.Now()
	formattedTime := now.Format("20060102150405")
	return fmt.Sprintf("42%s", formattedTime)
}

func crc16(data []byte) uint16 {
    var crc uint16 = 0xFFFF
    for _, b := range data {
        crc ^= uint16(b) << 8
        for i := 0; i < 8; i++ {
            if (crc & 0x8000) != 0 {
                crc = (crc << 1) ^ 0x1021
            } else {
                crc <<= 1
            }
        }
    }
    return crc & 0xFFFF
}

func removeAccents(s string) string {
    t := transform.Chain(
        norm.NFD,
        runes.Remove(runes.In(unicode.Mn)),
        norm.NFC,
    )
    result, _, err := transform.String(t, s)
    if err != nil {
        return s
    }
    return result
}