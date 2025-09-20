package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/fatih/color"
)

// Helper to capture stdout
func captureStdout(f func()) string {
	oldStdout := os.Stdout
	oldColorOutput := color.Output // Save original color.Output

	r, w, _ := os.Pipe()
	os.Stdout = w // Redirect os.Stdout

	var buf bytes.Buffer
	color.Output = &buf // Redirect color.Output to our buffer

	f()

	w.Close()        // Close the pipe writer to signal EOF
	io.Copy(&buf, r) // Copy anything written to the pipe to the buffer

	os.Stdout = oldStdout         // Restore original os.Stdout
	color.Output = oldColorOutput // Restore original color.Output

	return buf.String()
}

// Helper to reset global state for tests
var mu sync.Mutex // Mutex to protect global state during tests

func resetGlobalState() {
	mu.Lock()
	defer mu.Unlock()
	commands = make(map[string]command)
	history = []string{}
	arithmeticHistory = []string{}
	config = Config{
		Prompt:      "> ",
		PromptColor: "cyan",
		OutputColor: "blue",
		ErrorColor:  "red",
		HistoryFile: ".gosh_history",
		Theme:       "default",
	}
}

func TestRegisterCommand(t *testing.T) {
	resetGlobalState()
	var handlerCalled bool
	testHandler := func(args []string) { handlerCalled = true }

	registerCommand("testcmd", "A test command", testHandler)

	if _, exists := commands["testcmd"]; !exists {
		t.Errorf("Command 'testcmd' was not registered")
	}

	cmd := commands["testcmd"]
	if cmd.name != "testcmd" {
		t.Errorf("Expected command name 'testcmd', got '%s'", cmd.name)
	}
	if cmd.description != "A test command" {
		t.Errorf("Expected description 'A test command', got '%s'", cmd.description)
	}

	// Test if the handler can be called
	cmd.handler([]string{})
	if !handlerCalled {
		t.Errorf("Registered command handler was not called")
	}
}

func TestCmdHelp(t *testing.T) {
	resetGlobalState()
	registerCommand("test1", "Description 1", func(args []string) {})
	registerCommand("test2", "Description 2", func(args []string) {})

	output := captureStdout(func() {
		cmdHelp([]string{})
	})

	if !strings.Contains(output, "Available commands:") {
		t.Errorf("Output missing 'Available commands:'")
	}
	if !strings.Contains(output, "test1: Description 1") {
		t.Errorf("Output missing 'test1: Description 1'")
	}
	if !strings.Contains(output, "test2: Description 2") {
		t.Errorf("Output missing 'test2: Description 2'")
	}
}

func TestCmdHistory(t *testing.T) {
	resetGlobalState()
	history = []string{"cmd1", "cmd2"}

	output := captureStdout(func() {
		cmdHistory([]string{})
	})

	if !strings.Contains(output, "Command History:") {
		t.Errorf("Output missing 'Command History:'")
	}
	if !strings.Contains(output, "1: cmd1") {
		t.Errorf("Output missing '1: cmd1'")
	}
	if !strings.Contains(output, "2: cmd2") {
		t.Errorf("Output missing '2: cmd2'")
	}

	// Test empty history
	resetGlobalState()
	output = captureStdout(func() {
		cmdHistory([]string{})
	})
	if !strings.Contains(output, "No command history available.") {
		t.Errorf("Output missing 'No command history available.' for empty history")
	}
}

func TestCmdClear(t *testing.T) {
	resetGlobalState()
	history = []string{"cmd1", "cmd2"}

	output := captureStdout(func() {
		cmdClear([]string{})
	})

	if len(history) != 0 {
		t.Errorf("History not cleared, length is %d", len(history))
	}
	if !strings.Contains(output, "Command history cleared.") {
		t.Errorf("Output missing 'Command history cleared.'")
	}
}

func TestCmdEcho(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdEcho([]string{"hello", "world"})
	})
	if !strings.Contains(output, "Echo: hello world") {
		t.Errorf("Expected 'Echo: hello world', got '%s'", output)
	}

	output = captureStdout(func() {
		cmdEcho([]string{})
	})
	if !strings.Contains(output, "Usage: echo <message>") {
		t.Errorf("Expected usage message for no arguments, got '%s'", output)
	}
}

func TestCmdAdd(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdAdd([]string{"10", "5"})
	})
	if !strings.Contains(output, "Result: 15.000000") {
		t.Errorf("Expected 'Result: 15.000000', got '%s'", output)
	}
	if len(arithmeticHistory) != 1 || !strings.Contains(arithmeticHistory[0], "10.000000 + 5.000000 = 15.000000") {
		t.Errorf("Arithmetic history not updated correctly: %v", arithmeticHistory)
	}

	output = captureStdout(func() {
		cmdAdd([]string{"abc", "5"})
	})
	if !strings.Contains(output, "Both arguments must be numbers.") {
		t.Errorf("Expected error for non-numeric input, got '%s'", output)
	}

	output = captureStdout(func() {
		cmdAdd([]string{"10"})
	})
	if !strings.Contains(output, "Usage: add <num1> <num2>") {
		t.Errorf("Expected usage message for wrong number of arguments, got '%s'", output)
	}
}

func TestCmdMultiply(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdMultiply([]string{"10", "5"})
	})
	if !strings.Contains(output, "Result: 50.000000") {
		t.Errorf("Expected 'Result: 50.000000', got '%s'", output)
	}
	if len(arithmeticHistory) != 1 || !strings.Contains(arithmeticHistory[0], "10.000000 * 5.000000 = 50.000000") {
		t.Errorf("Arithmetic history not updated correctly: %v", arithmeticHistory)
	}
}

func TestCmdSubtract(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdSubtract([]string{"10", "5"})
	})
	if !strings.Contains(output, "Result: 5.000000") {
		t.Errorf("Expected 'Result: 5.000000', got '%s'", output)
	}
	if len(arithmeticHistory) != 1 || !strings.Contains(arithmeticHistory[0], "10.000000 - 5.000000 = 5.000000") {
		t.Errorf("Arithmetic history not updated correctly: %v", arithmeticHistory)
	}
}

func TestCmdDivision(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdDivision([]string{"10", "5"})
	})
	if !strings.Contains(output, "Result: 2.000000") {
		t.Errorf("Expected 'Result: 2.000000', got '%s'", output)
	}
	if len(arithmeticHistory) != 1 || !strings.Contains(arithmeticHistory[0], "10.000000 / 5.000000 = 2.000000") {
		t.Errorf("Arithmetic history not updated correctly: %v", arithmeticHistory)
	}

	output = captureStdout(func() {
		cmdDivision([]string{"10", "0"})
	})
	if !strings.Contains(output, "Cannot divide by zero.") {
		t.Errorf("Expected 'Cannot divide by zero.', got '%s'", output)
	}
}

func TestCmdModulus(t *testing.T) {
	resetGlobalState()
	output := captureStdout(func() {
		cmdModulus([]string{"10", "3"})
	})
	if !strings.Contains(output, "Result: 1") {
		t.Errorf("Expected 'Result: 1', got '%s'", output)
	}
	if len(arithmeticHistory) != 1 || !strings.Contains(arithmeticHistory[0], "10 % 3 = 1") {
		t.Errorf("Arithmetic history not updated correctly: %v", arithmeticHistory)
	}

	output = captureStdout(func() {
		cmdModulus([]string{"10", "0"})
	})
	if !strings.Contains(output, "Cannot perform modulus with zero.") {
		t.Errorf("Expected 'Cannot perform modulus with zero.', got '%s'", output)
	}
}

func TestCmdCat(t *testing.T) {
	resetGlobalState()
	testFilename := "test_cat_file.txt"
	testContent := "Hello, Cat!"
	os.WriteFile(testFilename, []byte(testContent), 0644)
	defer os.Remove(testFilename)

	output := captureStdout(func() {
		cmdCat([]string{testFilename})
	})
	if !strings.Contains(output, testContent) {
		t.Errorf("Expected file content '%s', got '%s'", testContent, output)
	}

	output = captureStdout(func() {
		cmdCat([]string{"non_existent_file.txt"})
	})
	if !strings.Contains(output, "Error reading file:") {
		t.Errorf("Expected error for non-existent file, got '%s'", output)
	}
}

func TestCmdWrite(t *testing.T) {
	resetGlobalState()
	testFilename := "test_write_file.txt"
	testContent := "Hello, Write!"
	defer os.Remove(testFilename)

	output := captureStdout(func() {
		cmdWrite([]string{testFilename, testContent})
	})
	if !strings.Contains(output, "File written successfully") {
		t.Errorf("Expected success message, got '%s'", output)
	}

	content, err := os.ReadFile(testFilename)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("Expected file content '%s', got '%s'", testContent, string(content))
	}
}

func TestCmdAppend(t *testing.T) {
	resetGlobalState()
	testFilename := "test_append_file.txt"
	initialContent := "Initial line\n"
	appendContent := "Appended line"
	defer os.Remove(testFilename)

	os.WriteFile(testFilename, []byte(initialContent), 0644)

	output := captureStdout(func() {
		cmdAppend([]string{testFilename, appendContent})
	})
	if !strings.Contains(output, "Content appended successfully") {
		t.Errorf("Expected success message, got '%s'", output)
	}

	content, err := os.ReadFile(testFilename)
	if err != nil {
		t.Fatalf("Failed to read appended file: %v", err)
	}
	expectedContent := initialContent + appendContent + "\n"
	if string(content) != expectedContent {
		t.Errorf("Expected file content '%s', got '%s'", expectedContent, string(content))
	}
}

func TestCmdConfig(t *testing.T) {
	resetGlobalState()
	config.Prompt = ">>> "
	config.Theme = "dark"

	output := captureStdout(func() {
		cmdConfig([]string{})
	})

	if !strings.Contains(output, "Prompt: >>> ") {
		t.Errorf("Output missing 'Prompt: >>> '")
	}
	if !strings.Contains(output, "Theme: dark") {
		t.Errorf("Output missing 'Theme: dark'")
	}
}

func TestCmdSet(t *testing.T) {
	resetGlobalState()

	// Test setting prompt
	output := captureStdout(func() {
		cmdSet([]string{"prompt", ">>> "})
	})
	if config.Prompt != ">>> " {
		t.Errorf("Prompt not updated, expected '>>> ', got '%s'", config.Prompt)
	}
	if !strings.Contains(output, "Configuration updated") {
		t.Errorf("Expected success message, got '%s'", output)
	}

	// Test setting unknown property
	output = captureStdout(func() {
		cmdSet([]string{"unknown_prop", "value"})
	})
	if !strings.Contains(output, "Unknown configuration property: unknown_prop") {
		t.Errorf("Expected error for unknown property, got '%s'", output)
	}
}

func TestCmdTheme(t *testing.T) {
	resetGlobalState()

	// Test setting dark theme
	output := captureStdout(func() {
		cmdTheme([]string{"dark"})
	})
	if config.Theme != "dark" {
		t.Errorf("Theme not updated, expected 'dark', got '%s'", config.Theme)
	}
	if config.PromptColor != "magenta" {
		t.Errorf("PromptColor not updated for dark theme")
	}
	if !strings.Contains(output, "Theme set to dark") {
		t.Errorf("Expected success message, got '%s'", output)
	}

	// Test setting unknown theme
	output = captureStdout(func() {
		cmdTheme([]string{"galaxy"})
	})
	if !strings.Contains(output, "Unknown theme: galaxy") {
		t.Errorf("Expected error for unknown theme, got '%s'", output)
	}
}

func TestCmdSave(t *testing.T) {
	resetGlobalState()
	history = []string{"ls", "pwd"}
	arithmeticHistory = []string{"1+1=2"}

	// Test save history
	testHistoryFilename := "command_history.txt" // Use the actual filename
	defer os.Remove(testHistoryFilename)
	output := captureStdout(func() {
		cmdSave([]string{"history"})
	})
	if !strings.Contains(output, fmt.Sprintf("Command history saved to %s", testHistoryFilename)) {
		t.Errorf("Expected history save message, got '%s'", output)
	}
	content, err := os.ReadFile(testHistoryFilename)
	if err != nil {
		t.Fatalf("Failed to read saved history file: %v", err)
	}
	if !strings.Contains(string(content), "ls\npwd\n") {
		t.Errorf("Saved history content incorrect: %s", string(content))
	}

	// Test save arithmetic
	testArithmeticFilename := "arithmetic_operations.txt" // Use the actual filename
	defer os.Remove(testArithmeticFilename)
	output = captureStdout(func() {
		cmdSave([]string{"arithmetic"})
	})
	if !strings.Contains(output, fmt.Sprintf("Arithmetic operations saved to %s", testArithmeticFilename)) {
		t.Errorf("Expected arithmetic save message, got '%s'", output)
	}
	content, err = os.ReadFile(testArithmeticFilename)
	if err != nil {
		t.Fatalf("Failed to read saved arithmetic file: %v", err)
	}
	if !strings.Contains(string(content), "1+1=2\n") {
		t.Errorf("Saved arithmetic content incorrect: %s", string(content))
	}

	// Test unknown save option
	output = captureStdout(func() {
		cmdSave([]string{"unknown"})
	})
	if !strings.Contains(output, "Unknown save option: unknown") {
		t.Errorf("Expected error for unknown save option, got '%s'", output)
	}

	// Test no arguments
	output = captureStdout(func() {
		cmdSave([]string{})
	})
	if !strings.Contains(output, "Save options:") {
		t.Errorf("Expected save options message, got '%s'", output)
	}
}

func TestSaveHistoryToFile(t *testing.T) {
	resetGlobalState()
	history = []string{"test_cmd_1", "test_cmd_2"}
	filename := "command_history.txt" // Use the actual filename
	defer os.Remove(filename)

	output := captureStdout(func() {
		saveHistoryToFile()
	})

	if !strings.Contains(output, fmt.Sprintf("Command history saved to %s", filename)) {
		t.Errorf("Expected success message, got '%s'", output)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read saved history file: %v", err)
	}
	expectedContent := "test_cmd_1\ntest_cmd_2\n"
	if string(content) != expectedContent {
		t.Errorf("Saved history content mismatch. Expected:\n%sGot:\n%s", expectedContent, string(content))
	}

	// Test with empty history
	resetGlobalState()
	output = captureStdout(func() {
		saveHistoryToFile()
	})
	if !strings.Contains(output, "No command history to save.") {
		t.Errorf("Expected 'No command history to save.' for empty history, got '%s'", output)
	}
}

func TestSaveArithmeticToFile(t *testing.T) {
	resetGlobalState()
	arithmeticHistory = []string{"10 + 5 = 15", "2 * 3 = 6"}
	filename := "arithmetic_operations.txt" // Use the actual filename
	defer os.Remove(filename)

	output := captureStdout(func() {
		saveArithmeticToFile()
	})

	if !strings.Contains(output, fmt.Sprintf("Arithmetic operations saved to %s", filename)) {
		t.Errorf("Expected success message, got '%s'", output)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read saved arithmetic file: %v", err)
	}
	expectedContent := "10 + 5 = 15\n2 * 3 = 6\n"
	if string(content) != expectedContent {
		t.Errorf("Saved arithmetic content mismatch. Expected:\n%sGot:\n%s", expectedContent, string(content))
	}

	// Test with empty arithmetic history
	resetGlobalState()
	output = captureStdout(func() {
		saveArithmeticToFile()
	})
	if !strings.Contains(output, "No arithmetic operations to save.") {
		t.Errorf("Expected 'No arithmetic operations to save.' for empty history, got '%s'", output)
	}
}

// Test for cmdPipe - this is a complex function due to stdout capture and command execution.
// A full test would involve mocking os.Stdout and the command handlers more thoroughly.
// For simplicity, we'll test a basic pipe scenario with existing commands.
func TestCmdPipe(t *testing.T) {
	resetGlobalState()
	registerCommand("echo", "Echo the input", cmdEcho)
	registerCommand("add", "Add two numbers", cmdAdd) // Add a simple command that takes args

	// Test a simple pipe: echo "hello" | echo
	output := captureStdout(func() {
		cmdPipe([]string{"echo", "hello", "|", "echo"})
	})

	// The pipe command itself prints the final output.
	// The inner echo will print "Echo: hello", and then the outer echo will print "Echo: Echo: hello"
	// This test is tricky because cmdPipe directly calls cmd.handler which prints to stdout.
	// The internal stdout redirection within cmdPipe is complex to test without deeper mocking.
	// For now, we'll check for the expected final output from the last command in the pipe.
	if !strings.Contains(output, "Echo: hello") {
		t.Errorf("Expected piped output to contain 'Echo: hello', got '%s'", output)
	}

	// Test a pipe with an unknown command
	output = captureStdout(func() {
		cmdPipe([]string{"echo", "test", "|", "unknowncmd"})
	})
	if !strings.Contains(output, "Unknown command in pipe: unknowncmd") {
		t.Errorf("Expected error for unknown command in pipe, got '%s'", output)
	}

	// Test insufficient arguments (empty args)
	output = captureStdout(func() {
		cmdPipe([]string{})
	})
	if !strings.Contains(output, "Usage: pipe <command1> | <command2> [| <command3> ...]") {
		t.Errorf("Expected usage message for insufficient arguments, got '%s'", output)
	}

	// Test insufficient arguments (only one command part)
	output = captureStdout(func() {
		cmdPipe([]string{"echo", "test"})
	})
	if !strings.Contains(output, "Usage: pipe <command1> | <command2> [| <command3> ...]") {
		t.Errorf("Expected usage message for insufficient arguments, got '%s'", output)
	}
}
