package main

import (
	"fmt"
	"todoCLI/sjcli"
)

func addCommand(data sjcli.CLIData) error {
	fmt.Printf("hello world!")
	return nil
}

func main() {
	cli := sjcli.NewCLIProgramm()
	cli.CommandHandler("print", addCommand)
	cli.Run()
}
