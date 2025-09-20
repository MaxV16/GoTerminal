package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

type command struct {
	name        string
	description string
	handler     func(args []string)
}

var commands = make(map[string]command)
var history []string
var arithmeticHistory []string // New global variable to store arithmetic operations

type Config struct {
	Prompt      string
	PromptColor string
	OutputColor string
	ErrorColor  string
	HistoryFile string
	Theme       string
}

var config = Config{
	Prompt:      "> ",
	PromptColor: "cyan",
	OutputColor: "blue",
	ErrorColor:  "red",
	HistoryFile: ".gosh_history",
	Theme:       "default",
}

func main() {

	// Register commands
	registerCommand("help", "Show help", cmdHelp)
	registerCommand("exit", "Exit the program", cmdExit)
	registerCommand("history", "Show command history", cmdHistory)
	registerCommand("clear", "Clear command history", cmdClear)
	registerCommand("echo", "Echo the input", cmdEcho)
	registerCommand("arithmetic", "Perform arithmetic operations", cmdArithmeticCommands)
	registerCommand("add", "Add two numbers", cmdAdd)
	registerCommand("multiply", "Multiply two numbers", cmdMultiply)
	registerCommand("subtract", "Subtract two numbers", cmdSubtract)
	registerCommand("division", "Divide two numbers", cmdDivision)
	registerCommand("modulus", "Calculate modulus of two numbers", cmdModulus)
	registerCommand("pipe", "Pipe output between commands", cmdPipe)
	registerCommand("cat", "Display file contents", cmdCat)
	registerCommand("write", "Write to a file", cmdWrite)
	registerCommand("append", "Append to a file", cmdAppend)
	registerCommand("config", "Configure terminal settings", cmdConfig)
	registerCommand("set", "Set a configuration value", cmdSet)
	registerCommand("theme", "Change color theme", cmdTheme)
	registerCommand("save", "Save history or arithmetic operations", cmdSave) // Register the new save command

	// Set up readline with command completion
	completer := readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("history"),
		readline.PcItem("clear"),
		readline.PcItem("echo"),
		readline.PcItem("arithmetic"),
		readline.PcItem("add"),
		readline.PcItem("multiply"),
		readline.PcItem("subtract"),
		readline.PcItem("division"),
		readline.PcItem("modulus"),
	)
	// Initialize readline with the completer
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: completer,
	})
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		os.Exit(1)
	}
	defer rl.Close()

	color.Cyan("Welcome to goTerminal!!!!!!!!!!!")
	color.Yellow("Type 'help' to see available commands. for arithmetic commands type 'arithmetic'")

	for {
		input, err := rl.Readline()
		if err != nil { // Handles Ctrl+D (EOF)
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		history = append(history, input)
		parts := strings.Fields(input)
		cmdName := parts[0]
		args := parts[1:]

		if cmd, exists := commands[cmdName]; exists {
			cmd.handler(args)
		} else {
			color.Red("Unknown command: %s", cmdName)
		}
	}

}

func registerCommand(name, description string, handler func(args []string)) {
	commands[name] = command{
		name:        name,
		description: description,
		handler:     handler,
	}
}

func cmdHelp(args []string) {
	color.Green("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
}

func cmdExit(args []string) {
	color.Yellow("Exiting the program. Goodbye!")
	os.Exit(0)
}
func cmdHistory(args []string) {
	if len(history) == 0 {
		color.Yellow("No command history available.")
		return
	}
	color.Green("Command History:")
	for i, cmd := range history {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}
}
func cmdClear(args []string) {
	history = []string{}
	color.Green("Command history cleared.")
}

func cmdEcho(args []string) {
	if len(args) == 0 {
		color.Red("Usage: echo <message>")
		return
	}
	message := strings.Join(args, " ")
	color.Blue("Echo: %s", message)
}

func cmdAdd(args []string) {
	if len(args) != 2 {
		color.Red("Usage: add <num1> <num2>")
		return
	}
	var num1, num2 float64
	_, err1 := fmt.Sscanf(args[0], "%f", &num1)
	_, err2 := fmt.Sscanf(args[1], "%f", &num2)
	if err1 != nil || err2 != nil {
		color.Red("Both arguments must be numbers.")
		return
	}
	result := num1 + num2
	arithmeticHistory = append(arithmeticHistory, fmt.Sprintf("%f + %f = %f", num1, num2, result)) // Store operation
	color.Blue("Result: %f", result)
}

func cmdMultiply(args []string) {
	if len(args) != 2 {
		color.Red("Usage: multiply <num1> <num2>")
		return
	}
	var num1, num2 float64
	_, err1 := fmt.Sscanf(args[0], "%f", &num1)
	_, err2 := fmt.Sscanf(args[1], "%f", &num2)
	if err1 != nil || err2 != nil {
		color.Red("Both arguments must be numbers.")
		return
	}
	result := num1 * num2
	arithmeticHistory = append(arithmeticHistory, fmt.Sprintf("%f * %f = %f", num1, num2, result)) // Store operation
	color.Blue("Result: %f", result)
}
func cmdSubtract(args []string) {
	if len(args) != 2 {
		color.Red("Usage: subtract <num1> <num2>")
		return
	}
	var num1, num2 float64
	_, err1 := fmt.Sscanf(args[0], "%f", &num1)
	_, err2 := fmt.Sscanf(args[1], "%f", &num2)
	if err1 != nil || err2 != nil {
		color.Red("Both arguments must be numbers.")
		return
	}
	result := num1 - num2
	arithmeticHistory = append(arithmeticHistory, fmt.Sprintf("%f - %f = %f", num1, num2, result)) // Store operation
	color.Blue("Result: %f", result)
}

// List of arithmetic commands
func cmdArithmeticCommands(args []string) {
	color.Green("Arithmetic Commands:")
	color.Yellow("  add <num1> <num2>       - Add two numbers")
	color.Yellow("  subtract <num1> <num2>  - Subtract two numbers")
	color.Yellow("  multiply <num1> <num2>  - Multiply two numbers")
	color.Yellow("  division <num1> <num2>   - Divide two numbers")
	color.Yellow("  modulus <num1> <num2>    - Calculate modulus of two numbers")
}

func cmdDivision(args []string) {
	if len(args) != 2 {
		color.Red("Usage: division <num1> <num2>")
		return
	}
	var num1, num2 float64
	_, err1 := fmt.Sscanf(args[0], "%f", &num1)
	_, err2 := fmt.Sscanf(args[1], "%f", &num2)
	if err1 != nil || err2 != nil {
		color.Red("Both arguments must be numbers.")
		return
	}
	if num2 == 0 {
		color.Red("Cannot divide by zero.")
		return
	}
	result := num1 / num2
	arithmeticHistory = append(arithmeticHistory, fmt.Sprintf("%f / %f = %f", num1, num2, result)) // Store operation
	color.Blue("Result: %f", result)
}
func cmdModulus(args []string) {
	if len(args) != 2 {
		color.Red("Usage: modulus <num1> <num2>")
		return
	}
	var num1, num2 float64
	_, err1 := fmt.Sscanf(args[0], "%f", &num1)
	_, err2 := fmt.Sscanf(args[1], "%f", &num2)
	if err1 != nil || err2 != nil {
		color.Red("Both arguments must be numbers.")
		return
	}
	if num2 == 0 {
		color.Red("Cannot perform modulus with zero.")
		return
	}
	result := int(num1) % int(num2)
	arithmeticHistory = append(arithmeticHistory, fmt.Sprintf("%d %% %d = %d", int(num1), int(num2), result)) // Store operation
	color.Blue("Result: %d", result)
}

// Pipe command handler
func cmdPipe(args []string) {
	if len(args) < 2 {
		color.Red("Usage: pipe <command1> | <command2> [| <command3> ...]")
		return
	}

	// Check if a pipe symbol exists
	hasPipe := false
	for _, arg := range args {
		if arg == "|" {
			hasPipe = true
			break
		}
	}
	if !hasPipe {
		color.Red("Usage: pipe <command1> | <command2> [| <command3> ...]")
		return
	}

	// Split commands by pipe character
	var cmdChain [][]string
	currentCmd := []string{}
	for _, arg := range args {
		if arg == "|" {
			if len(currentCmd) > 0 {
				cmdChain = append(cmdChain, currentCmd)
				currentCmd = []string{}
			}
		} else {
			currentCmd = append(currentCmd, arg)
		}
	}
	if len(currentCmd) > 0 {
		cmdChain = append(cmdChain, currentCmd)
	}

	// Execute commands in chain with piping
	var lastOutput string
	for i, cmdParts := range cmdChain {
		if len(cmdParts) == 0 {
			continue
		}

		cmdName := cmdParts[0]
		cmdArgs := cmdParts[1:]

		if i > 0 {
			// Prepend last output to arguments
			cmdArgs = append([]string{lastOutput}, cmdArgs...)
		}

		if cmd, exists := commands[cmdName]; exists {
			// Create a custom handler to capture output
			output := ""
			originalHandler := cmd.handler
			cmd.handler = func(args []string) {
				buf := new(bytes.Buffer)
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				originalHandler(args)

				w.Close()
				os.Stdout = old
				io.Copy(buf, r)
				output = buf.String()
			}

			cmd.handler(cmdArgs)
			lastOutput = output
			cmd.handler = originalHandler // Restore original handler
		} else {
			color.Red("Unknown command in pipe: %s", cmdName)
			return
		}
	}

	color.Blue(lastOutput)
}

// File I/O commands
func cmdCat(args []string) {
	if len(args) != 1 {
		color.Red("Usage: cat <filename>")
		return
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		color.Red("Error reading file: %v", err)
		return
	}

	color.Blue(string(content))
}

func cmdWrite(args []string) {
	if len(args) < 2 {
		color.Red("Usage: write <filename> <content>")
		return
	}

	content := strings.Join(args[1:], " ")
	err := os.WriteFile(args[0], []byte(content), 0644)
	if err != nil {
		color.Red("Error writing file: %v", err)
		return
	}

	color.Green("File written successfully")
}

func cmdAppend(args []string) {
	if len(args) < 2 {
		color.Red("Usage: append <filename> <content>")
		return
	}

	f, err := os.OpenFile(args[0], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		color.Red("Error opening file: %v", err)
		return
	}
	defer f.Close()

	content := strings.Join(args[1:], " ")
	if _, err = f.WriteString(content + "\n"); err != nil {
		color.Red("Error appending to file: %v", err)
		return
	}

	color.Green("Content appended successfully")
}

// Configuration commands
func cmdConfig(args []string) {
	color.Green("Current Configuration:")
	fmt.Printf("Prompt: %s\n", config.Prompt)
	fmt.Printf("Prompt Color: %s\n", config.PromptColor)
	fmt.Printf("Output Color: %s\n", config.OutputColor)
	fmt.Printf("Error Color: %s\n", config.ErrorColor)
	fmt.Printf("History File: %s\n", config.HistoryFile)
	fmt.Printf("Theme: %s\n", config.Theme)
}

// Set configuration property
func cmdSet(args []string) {
	if len(args) != 2 {
		color.Red("Usage: set <property> <value>")
		return
	}

	property := args[0]
	value := args[1]

	switch property {
	case "prompt":
		config.Prompt = value
	case "prompt_color":
		config.PromptColor = value
	case "output_color":
		config.OutputColor = value
	case "error_color":
		config.ErrorColor = value
	case "history_file":
		config.HistoryFile = value
	default:
		color.Red("Unknown configuration property: %s", property)
		return
	}

	color.Green("Configuration updated")
}

// Theme command handler
func cmdTheme(args []string) {
	if len(args) != 1 {
		color.Red("Available themes: default, dark, light, solarized")
		return
	}

	theme := args[0]
	switch theme {
	case "default":
		config.PromptColor = "cyan"
		config.OutputColor = "blue"
		config.ErrorColor = "red"
	case "dark":
		config.PromptColor = "magenta"
		config.OutputColor = "green"
		config.ErrorColor = "yellow"
	case "light":
		config.PromptColor = "blue"
		config.OutputColor = "black"
		config.ErrorColor = "red"
	case "solarized":
		config.PromptColor = "yellow"
		config.OutputColor = "cyan"
		config.ErrorColor = "red"
	default:
		color.Red("Unknown theme: %s", theme)
		return
	}

	config.Theme = theme
	color.Green("Theme set to %s", theme)
}

// New save command handler
func cmdSave(args []string) {
	if len(args) == 0 {
		color.Green("Save options:")
		color.Yellow("  save history - Save command history to a file")
		color.Yellow("  save arithmetic - Save past arithmetic operations to a file")
		return
	}

	switch args[0] { //
	case "history":
		saveHistoryToFile()
	case "arithmetic":
		saveArithmeticToFile()
	default:
		color.Red("Unknown save option: %s", args[0])
		color.Yellow("Please choose 'history' or 'arithmetic'.")
	}
}

func saveHistoryToFile() {
	if len(history) == 0 {
		color.Yellow("No command history to save.")
		return
	}

	filename := "command_history.txt"
	file, err := os.Create(filename)
	if err != nil {
		color.Red("Error creating file %s: %v", filename, err)
		return
	}
	defer file.Close()

	for _, cmd := range history {
		_, err := file.WriteString(cmd + "\n")
		if err != nil {
			color.Red("Error writing to file %s: %v", filename, err)
			return
		}
	}
	color.Green("Command history saved to %s", filename)
}

func saveArithmeticToFile() {
	if len(arithmeticHistory) == 0 {
		color.Yellow("No arithmetic operations to save.")
		return
	}

	filename := "arithmetic_operations.txt"
	file, err := os.Create(filename)
	if err != nil {
		color.Red("Error creating file %s: %v", filename, err)
		return
	}
	defer file.Close()

	for _, op := range arithmeticHistory {
		_, err := file.WriteString(op + "\n")
		if err != nil {
			color.Red("Error writing to file %s: %v", filename, err)
			return
		}
	}
	color.Green("Arithmetic operations saved to %s", filename)
}
