package output

import (
	"encoding/json"
	"fmt"
	"io"
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

// KV prints a simple aligned key/value list to stdout
func KV(pairs [][2]string) { KVTo(os.Stdout, pairs) }

// KVTo prints a simple aligned key/value list to w
func KVTo(w io.Writer, pairs [][2]string) {
	maxKey := 0
	for _, p := range pairs {
		if len(p[0]) > maxKey {
			maxKey = len(p[0])
		}
	}
	for _, p := range pairs {
		fmt.Fprintf(w, "  %-*s  %s\n", maxKey, p[0], p[1])
	}
}
