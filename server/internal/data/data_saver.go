package data

import (
	"fmt"
	"os"
)

var LogFile, err = os.OpenFile("snap.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

func init() {
	if err != nil {
		panic(err)
	}
}

func LogCommand(op, key, value string) {
	line := fmt.Sprintf("%s %s %s\n", op, key, value)
	_, err := LogFile.WriteString(line)
	if err != nil {
		fmt.Println("Failed to write log:", err)
	}
}
