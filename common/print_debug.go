package common

import (
	"fmt"
	"reflect"
)

const (
	CONSOLE_GREEN   = "\033[32m"
	CONSOLE_RED     = "\033[31m"
	CONSOLE_YELLOW  = "\033[33m"
	CONSOLE_BLUE    = "\033[34m"
	CONSOLE_MAGENTA = "\033[35m"
	CONSOLE_CYAN    = "\033[36m"
	CONSOLE_RESET   = "\033[0m"
)

func PrintDebug(message string, v interface{}) {

	// Use reflection to get the type and value of the variable
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Print the type and value
	fmt.Printf("Message: %s, Type: %s, Value: %v\n", message, typ, val)
}

func PrintError(message string) {
	fmt.Printf("%sERROR:%s %s \n", CONSOLE_RED, CONSOLE_RESET, message)
}

func PrintWarning(message string) {
	fmt.Printf("%sWARNING:%s %s \n", CONSOLE_YELLOW, CONSOLE_RESET, message)
}

func PrintInfo(message string) {
	fmt.Printf("%sINFO:%s %s \n", CONSOLE_GREEN, CONSOLE_RESET, message)
}

func PrintValid(message string) {
	fmt.Printf("%sVALID:%s %s \n", CONSOLE_GREEN, CONSOLE_RESET, message)
}

func PrintInvalid(message string) {
	fmt.Printf("%sINVALID:%s %s \n", CONSOLE_RED, CONSOLE_RESET, message)
}
