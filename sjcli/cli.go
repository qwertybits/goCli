package sjcli

import (
	"fmt"
	"os"
)

type CLIHandler func(CLIData) error

type CLI struct {
	commands map[string]CLIHandler
}

func NewCLIProgramm() CLI {
	return CLI{commands: make(map[string]CLIHandler)}
}

func (cli *CLI) CommandHandler(name string, handler CLIHandler) {
	cli.commands[name] = handler
}

func (cli *CLI) Run() {
	command, data := parseInput()
	exec, exist := cli.commands[command]
	if !exist {
		fmt.Printf("Unknow command. Enter help to see avialable commands")
		return
	}
	if err := exec(data); err != nil {
		fmt.Printf("%v", err)
	}
}

func parseInput() (string, CLIData) {
	args := os.Args
	if len(args) <= 1 {
		return "", CLIData{}
	}
	if len(args) == 2 {
		return args[1], CLIData{}
	}
	command := args[1]
	args = args[2:]
	return command, CLIData{arguments: args}
}
