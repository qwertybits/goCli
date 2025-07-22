package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type CLICommand struct {
	cmd       string
	arguments []string
}

var aviableCommands = map[string]bool{
	"add":              true,
	"update":           true,
	"delete":           true,
	"mark-in-progress": true,
	"mark-done":        true,
	"list":             true,
}

const defaultJsonPath = "task.json"

func loadTasksFromJson(path string) ([]TaskObj, error) {
	var result = make([]TaskObj, 0)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func exportTasksToJson(path string, storage []TaskObj) error {
	bytes, err := json.MarshalIndent(storage, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func addTask(content string, storage *[]TaskObj) {
	*storage = append(*storage, NewTask(len(*storage), content))
}

func printTasks(storage *[]TaskObj) {
	for _, task := range *storage {
		fmt.Printf("%v\n", task)
	}
}

func updateStatus(id int, new string, storage *[]TaskObj) error {
	if !validateId(id, storage) {
		return errors.New("id range out")
	}
	return (*storage)[id].SetStatus(new)
}

func validateId(id int, storage *[]TaskObj) bool {
	return !(id >= len(*storage) || id < 0)
}

func updateDescription(id int, new string, storage *[]TaskObj) error {
	if !validateId(id, storage) {
		return errors.New("id range out")
	}
	return (*storage)[id].SetDescription(new)
}

func deleteTaskById(id int, storage *[]TaskObj) error {

	if !validateId(id, storage) {
		return errors.New("id range out")
	}

	for i := id; i < len(*storage)-1; i++ {
		(*storage)[i] = (*storage)[i+1]
		(*storage)[i].SetId(i)
	}
	(*storage) = (*storage)[:len((*storage))-1]
	return nil
}

func getCommand() (CLICommand, error) {
	args := os.Args
	if len(args) <= 1 {
		return CLICommand{}, errors.New("not enough arguments")
	}
	if len(args) == 2 {
		return CLICommand{args[1], make([]string, 0)}, nil
	}
	cmdArguments := args[2:]
	return CLICommand{args[1], cmdArguments}, nil
}

func valideArgumentCount(input *CLICommand, need int) bool {
	return len(input.arguments) >= need
}

func Run() {
	var tasks, err = loadTasksFromJson(defaultJsonPath)
	if err != nil {
		fmt.Printf("loadTaskFromJson: %v", err)
		return
	}

	command, err := getCommand()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	value := aviableCommands[command.cmd]
	fmt.Printf("%v\n%v", value, command)

	if err := exportTasksToJson(defaultJsonPath, tasks); err != nil {
		fmt.Printf("%v", err)
	}
}
