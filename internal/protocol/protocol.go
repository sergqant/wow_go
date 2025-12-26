package protocol

import "fmt"

func Challenge(ch string, diff int, timeoutMs int) string {
	return fmt.Sprintf("CHALLENGE %s %d %d\n", ch, diff, timeoutMs)
}

func OK(quote string) string {
	return fmt.Sprintf("OK %s\n", quote)
}

func Error(reason string) string {
	return fmt.Sprintf("ERROR %s\n", reason)
}
