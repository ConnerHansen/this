package this

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func error(msg ...string) {
	write(os.Stderr, "\033[31m"+strings.Join(msg, " ")+"\033[0m")
}

func errorln(msg ...string) {
	write(os.Stderr, "\033[31m"+strings.Join(msg, " ")+"\033[0m\n")
}

func info(msgs ...string) {
	write(os.Stdout, "\033[39m", strings.Join(msgs, " "), "\033[0m")
}

func infoln(msgs ...string) {
	write(os.Stdout, "\033[39m", strings.Join(msgs, " "), "\033[0m\n")
}

func success(msgs ...string) {
	write(os.Stdout, "\033[32m", strings.Join(msgs, " "), "\033[0m")
}

func successln(msgs ...string) {
	write(os.Stdout, "\033[32m", strings.Join(msgs, " "), "\033[0m\n")
}

func warn(msgs ...string) {
	write(os.Stdout, "\033[33m", strings.Join(msgs, " "), "\033[0m")
}

func warnln(msgs ...string) {
	write(os.Stdout, "\033[33m", strings.Join(msgs, " "), "\033[0m\n")
}

func write(writer *os.File, s ...string) {
	writer.WriteString(strings.Join(s, ""))
}

// Fail fails the currently running test immediately
func Fail() {
	panic("Test Failure")
}

// Skip halts the current test and skips it
func Skip() {
	panic("Skip this.Should Test")
}

// GomegaFailHandler the fail handler to assign to Gomega for tests
func GomegaFailHandler(message string, callerSkip ...int) {
	messageLines := strings.Split(message, "\n")
	for i, line := range messageLines {
		messageLines[i] = "\t" + line
	}

	errorln("\n", strings.Join(messageLines, "\n"))
	panic("Test Failure")
}

// Should uses descriptive naming to run the tests
func Should(description string, t *testing.T, do func()) {
	startAt := time.Now()

	defer func() {
		if r := recover(); r != nil {
			lineSheer := 6
			stack := string(debug.Stack())
			stackLines := strings.Split(stack, "\n")
			newpacked := make([]string, len(stackLines)-lineSheer)

			if r == "Skip this.Should Test" {
				if testing.Verbose() {
					warnln(" Skipped")
				} else {
					warn("\u2022")
				}
			} else {
				// This wasn't a skip -- fail it!
				t.Fail()

				if testing.Verbose() {
					error(fmt.Sprintf(" Failed (%.3fs)", time.Now().Sub(startAt).Seconds()))
				} else {
					errorln("["+t.Name()+"]", description, fmt.Sprintf("-- Failed (%.3fs)", time.Now().Sub(startAt).Seconds()))
				}

				if r != "Test Failure" {
					error("\n\t"+fmt.Sprint(r), "\n\t")
				}

				newpacked[0] = "\n\t" + "Test failed at:\n\n\t" + stackLines[0]
				for i := lineSheer + 1; i < len(stackLines); i++ {
					newpacked[i-lineSheer] = "\t" + stackLines[i]
				}

				errorln(strings.Join(newpacked, "\n"))
			}
		}
	}()

	if testing.Verbose() {
		info("  ", description+":")
	}

	do()

	if testing.Verbose() {
		// if !t.Failed() {
		successln(fmt.Sprintf(" Passed (%.3fs)", time.Now().Sub(startAt).Seconds()))
		// } else {
		// 	errorln(" Failed")
		// }
	} else {
		success("\u2022")
		// Non-verbose logging, make sure we detail the test failures
		// if t.Failed() {
		// 	errorln("["+t.Name()+"]", description, fmt.Sprintf("-- Failed (%.3fs)", time.Now().Sub(startAt).Seconds()))
		// }
	}
}
