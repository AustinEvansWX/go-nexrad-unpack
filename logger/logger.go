package logger

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

func log(c color.Attribute, format string, msg ...interface{}) {
	bold := color.New(c, color.Bold)
	bold.Printf(format+"\n", msg...)
}

func Success(format string, msg ...interface{}) {
	log(color.FgGreen, format, msg...)
}

func Error(format string, msg ...interface{}) {
	log(color.FgRed, format, msg...)
}

func Warning(format string, msg ...interface{}) {
	log(color.FgYellow, format, msg...)
}

func Info(format string, msg ...interface{}) {
	log(color.FgCyan, format, msg...)
}

func PrettyPrintJSON(data interface{}) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "  ")
	err := enc.Encode(data)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buffer.String())
}
