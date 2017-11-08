package this

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

var (
	failPanicMessage = "Test Failure"
	skipPanicMessage = "Skip this.Should Test"

	debugColor   = "\033[36m"
	defaultColor = "\033[39m"

	failChar        = ""
	failColor       = "\033[31m"
	failVerboseChar = "\u2055"

	warningChar        = "|"
	warningColor       = "\033[33m"
	warningVerboseChar = "\u2192"

	successChar        = "\u2022"
	successColor       = "\033[32m"
	successVerboseChar = "\u2714"
)

func cyan(msgs ...string) {
	write(os.Stdout, debugColor, strings.Join(msgs, " "), "\033[0m")
}

func cyanln(msgs ...string) {
	write(os.Stdout, debugColor, strings.Join(msgs, " "), "\033[0m\n")
}

func green(msgs ...string) {
	write(os.Stdout, successColor, strings.Join(msgs, " "), "\033[0m")
}

func greenln(msgs ...string) {
	write(os.Stdout, successColor, strings.Join(msgs, " "), "\033[0m\n")
}

func red(msg ...string) {
	write(os.Stderr, failColor, strings.Join(msg, " "), "\033[0m")
}

func redln(msg ...string) {
	write(os.Stderr, failColor, strings.Join(msg, " "), "\033[0m\n")
}

func white(msgs ...string) {
	write(os.Stdout, defaultColor, strings.Join(msgs, " "), "\033[0m")
}

func whiteln(msgs ...string) {
	write(os.Stdout, defaultColor, strings.Join(msgs, " "), "\033[0m\n")
}

func write(writer *os.File, s ...string) {
	writer.WriteString(strings.Join(s, ""))
}

func yellow(msgs ...string) {
	write(os.Stdout, warningColor, strings.Join(msgs, " "), "\033[0m")
}

func yellowln(msgs ...string) {
	write(os.Stdout, warningColor, strings.Join(msgs, " "), "\033[0m\n")
}

// Fail fails the currently running test immediately
func Fail() {
	panic(failPanicMessage)
}

// Skip halts the current test and skips it
func Skip() {
	panic(skipPanicMessage)
}

// GomegaFailHandler the fail handler to assign to Gomega for tests
func GomegaFailHandler(message string, callerSkip ...int) {
	messageLines := strings.Split(message, "\n")
	for i, line := range messageLines {
		messageLines[i] = line
	}

	// Pass the assertion error up to our fail handler inside the Should call
	panic(failPanicMessage + strings.Join(messageLines, "\n"))
}

// Should uses descriptive naming to run the tests
func Should(description string, t *testing.T, do func()) {
	startAt := time.Now()
	verbose := testing.Verbose()

	defer func() {
		if r := recover(); r != nil {
			// Capture skip events and don't fail the test
			if r == skipPanicMessage {
				if verbose {
					yellowln(fmt.Sprintf(" %s (%.3fs)", warningVerboseChar, time.Now().Sub(startAt).Seconds()))
				} else {
					yellow(warningChar)
				}
			} else {
				// This wasn't a skip so it must have been something bad
				t.Fail()

				if verbose {
					red(fmt.Sprintf(" %s (%.3fs)", failVerboseChar, time.Now().Sub(startAt).Seconds()))
				} else {
					// Since this isn't verbose, we don't have any information about the
					// currently running test
					redln("\n["+t.Name()+"]", description,
						fmt.Sprintf("%s(%.3fs)", failChar, time.Now().Sub(startAt).Seconds()))
				}

				if r != failPanicMessage {
					str := fmt.Sprint(r)
					if strings.HasPrefix(str, failPanicMessage) {
						red("\n"+fmt.Sprint(strings.Replace(str, failPanicMessage, "", 1)), "\n\t")
					} else {
						red("\n\t"+fmt.Sprint(r), "\n\t")
					}
				}

				lineSheer := 6 /* How many lines do we want to remove from the stack */
				stack := string(debug.Stack())
				stackLines := strings.Split(stack, "\n")
				newpacked := make([]string, len(stackLines)-lineSheer)

				newpacked[0] = "\n\t" + "Test failed at:\n\n\t" + stackLines[0]
				for i := lineSheer + 1; i < len(stackLines); i++ {
					newpacked[i-lineSheer] = "\t" + stackLines[i]
				}

				redln(strings.Join(newpacked, "\n"))
			}
		}
	}()

	if verbose {
		cyan("["+t.Name()+"]", description+":")
	}

	// Since we allow tests to continue to execute, if someone has marked the
	// base test *itself* to be skipped we should skip all subsequent should
	// statements.
	if t.Skipped() {
		Skip()
	}

	do()

	if verbose {
		greenln(fmt.Sprintf(" %s (%.3fs)", successVerboseChar, time.Now().Sub(startAt).Seconds()))
	} else {
		green(successChar)
	}
}
