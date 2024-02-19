package utils

import "fmt"

// PrintBuildTag - Method for printing build tags into terminal during start application
func PrintBuildTag(tagName, tagValue string) {
	fmt.Printf("Build %s: %s\n", tagName, tagValue)
}
