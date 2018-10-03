package goelster

import (
	"fmt"
	"log"

	"github.com/brutella/can"
)

// logFrame logs a frame with the same format as candump from can-utils.
func logFrame(frm can.Frame) {
	data := trimSuffix(frm.Data[:], 0x00)
	length := fmt.Sprintf("[%x]", frm.Length)

	chars := fmt.Sprintf("'%s'", printableString(data[:]))
	rcvr := ReceiverId(frm.Data[:2])
	formatted := fmt.Sprintf("%-4x %-3s % -24X %-10s %6x ", frm.ID, length, data, chars, rcvr)

	reg, payload := Payload(data)
	formatted += fmt.Sprintf("%04X ", reg)

	if data[0]&Data != 0 {
		if r := Reading(reg); r != nil {
			val := DecodeValue(payload, r.Type)
			valStr := payloadString(val)

			formatted += fmt.Sprintf("%-20s %8s", left(r.Name, 20), valStr)
		}
	}

	log.Println(formatted)
}

func payloadString(val interface{}) string {
	if _, ok := val.(float64); ok {
		return fmt.Sprintf("%6.1f", val)
	}

	return fmt.Sprintf("%v", val)
}

func left(s string, chars int) string {
	l := len(s)
	if chars < l {
		l = chars
	}
	return s[:l]
}

// trim returns a subslice of s by slicing off all trailing b bytes.
func trimSuffix(s []byte, b byte) []byte {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != b {
			return s[:i+1]
		}
	}

	return []byte{}
}

// printableString creates a string from s and replaces non-printable bytes (i.e. 0-32, 127)
// with '.' – similar how candump from can-utils does it.
func printableString(s []byte) string {
	var ascii []byte
	for _, b := range s {
		if b < 32 || b > 126 {
			b = byte('.')

		}
		ascii = append(ascii, b)
	}

	return string(ascii)
}
