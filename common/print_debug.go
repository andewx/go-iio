package common

import (
	"fmt"
	"reflect"
)

// DebugMode is a boolean that determines if debug messages should be printed
var DebugMode bool

// VerboseMode is a boolean that determines if verbose messages should be printed
var VerboseMode bool

const (
	CONSOLE_GREEN   = "\033[32m" // Green
	CONSOLE_RED     = "\033[31m" // Red
	CONSOLE_YELLOW  = "\033[33m" // Yellow
	CONSOLE_BLUE    = "\033[34m" // Blue
	CONSOLE_MAGENTA = "\033[35m" // Magenta
	CONSOLE_CYAN    = "\033[36m" // Cyan
	CONSOLE_RESET   = "\033[0m"  // Reset
)

// PrintDebug prints the debug message
func PrintDebug(message string, v interface{}) {

	// Use reflection to get the type and value of the variable
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Print the type and value
	if DebugMode {
		fmt.Printf("Message: %s, Type: %s, Value: %v\n", message, typ, val)
	}
}

// PrintError prints the error message
func PrintError(message string) {
	fmt.Printf("%sERROR:%s %s \n", CONSOLE_RED, CONSOLE_RESET, message)
}

// PrintWarning prints the warning message
func PrintWarning(message string) {
	if DebugMode {
		fmt.Printf("%sWARNING:%s %s \n", CONSOLE_YELLOW, CONSOLE_RESET, message)
	}
}

// PrintInfo prints the info message
func PrintInfo(message string) {
	if VerboseMode {
		fmt.Printf("%sINFO:%s %s \n", CONSOLE_GREEN, CONSOLE_RESET, message)
	}
}

// PrintValid prints the valid message
func PrintValid(message string) {
	fmt.Printf("%sVALID:%s %s \n", CONSOLE_GREEN, CONSOLE_RESET, message)
}
