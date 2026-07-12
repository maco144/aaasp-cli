package output

import (
	"encoding/json"
	"fmt"
	"os"
)

var jsonMode bool

func SetJSON(v bool) { jsonMode = v }

func JSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func Line(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

func Fatal(format string, args ...any) {
	Error(format, args...)
	os.Exit(1)
}

func IsJSON() bool { return jsonMode }

// Table prints a simple aligned key/value list
func KV(pairs [][2]string) {
	maxKey := 0
	for _, p := range pairs {
		if len(p[0]) > maxKey {
			maxKey = len(p[0])
		}
	}
	for _, p := range pairs {
		fmt.Printf("  %-*s  %s\n", maxKey, p[0], p[1])
	}
}
