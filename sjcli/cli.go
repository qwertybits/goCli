package sjcli

import (
	"flag"
	"fmt"
	"os"
)

type CLIHandler func(CLIData) error

type CLI struct {
	commands     map[string]CLIHandler
	parsingFlags map[string]any
}

func NewCLIProgramm() CLI {
	_ = flag.String("", "", "")
	return CLI{commands: make(map[string]CLIHandler), parsingFlags: make(map[string]any)}
}

func (cli *CLI) CommandHandler(name string, handler CLIHandler) {
	cli.commands[name] = handler
}

func (cli *CLI) Run() {
	command, data := cli.parseInput()
	exec, exist := cli.commands[command]
	if !exist {
		fmt.Printf("Unknow command. Enter help to see avialable commands")
		return
	}
	if err := exec(data); err != nil {
		fmt.Printf("%v", err)
	}
}

func (cli *CLI) parseInput() (string, CLIData) {
	args := os.Args
	if len(args) <= 1 {
		return "", CLIData{make([]string, 0), make(map[string]any)}
	}
	if len(args) == 2 {
		return args[1], CLIData{}
	}
	command := args[1]
	args = args[2:]
	flag.Parse()
	return command, CLIData{arguments: args, flags: cli.parsingFlags}
}

func (cli *CLI) FlagInt(name string, defaultValue int) {
	cli.parsingFlags[name] = flag.Int(name, defaultValue, "")
}
