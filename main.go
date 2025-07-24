package main

import (
	"fmt"
	"todoCLI/sjcli"
)

func addCommand(data sjcli.CLIData) error {
	some, ok := data.GetInt("price")
	fmt.Printf("%v\t%v", some, ok)
	return nil
}

func main() {
	app := sjcli.NewCLIProgramm()
	app.CommandHandler("print", addCommand)
	app.Run()
}
