package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func loadTasksFromJson(path string) ([]TaskObj, error) {
	var result = make([]TaskObj, 0)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return result, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&result)
	return result, err
}

func exportTasksToJson(path string, storage []TaskObj) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(storage)
	return err
}

func addTask(input CLICommand, storage *[]TaskObj) error {
	if !valideArgumentCount(&input, 1) {
		return errors.New("addTask: not enough arguments")
	}
	content := input.arguments[0]
	id := len(*storage)
	*storage = append(*storage, NewTask(id, content))
	return nil
}

func printTasks(input CLICommand, storage *[]TaskObj) error {
	filter := ANY_STATUS
	if valideArgumentCount(&input, 1) {
		newFilter, exist := stringToStatus[input.arguments[0]]
		if !exist {
			return errors.New("list: unknow status")
		}
		filter = newFilter
	}
	var atLeastOne = false
	for _, task := range *storage {
		if task.IsSameStatus(filter) {
			atLeastOne = true
			fmt.Printf("%v\n\n", task)
		}
	}
	if !atLeastOne {
		fmt.Printf("No tasks\n")
	}
	return nil
}

func updateStatus(input CLICommand, storage *[]TaskObj) error {

	if !valideArgumentCount(&input, 2) {
		return errors.New("updateStatus: not enough arguments")
	}

	newStatus, exist := stringToStatus[input.arguments[0]]
	if !exist {
		return errors.New("updateStatus: invalid status")
	}

	id, err := strconv.Atoi(input.arguments[1])
	if err != nil {
		return errors.New("updateStatus: invalid id")
	}
	if !validateId(id, storage) {
		return errors.New("id range out")
	}
	(*storage)[id].SetStatus(newStatus)
	return nil
}

func updateDescription(input CLICommand, storage *[]TaskObj) error {

	if !valideArgumentCount(&input, 2) {
		return errors.New("updateDescription: not enough arguments")
	}

	id, err := strconv.Atoi(input.arguments[0])
	new := input.arguments[1]

	if err != nil {
		return errors.New("updateDescription: invalid id")
	}
	if !validateId(id, storage) {
		return errors.New("id range out")
	}
	return (*storage)[id].SetDescription(new)
}

func deleteTaskById(input CLICommand, storage *[]TaskObj) error {
	if !valideArgumentCount(&input, 1) {
		return errors.New("deleteTask: not enough arguments")
	}
	id, err := strconv.Atoi(input.arguments[0])
	if err != nil {
		return errors.New("deleteTask: invalid id")
	}
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

func getCommand() CLICommand {
	args := os.Args
	if len(args) <= 1 {
		return CLICommand{}
	}
	if len(args) == 2 {
		return CLICommand{args[1], make([]string, 0)}
	}
	cmdArguments := args[2:]
	return CLICommand{args[1], cmdArguments}
}

func valideArgumentCount(input *CLICommand, need int) bool {
	return len(input.arguments) >= need
}

func validateId(id int, storage *[]TaskObj) bool {
	return !(id >= len(*storage) || id < 0)
}

func helpCommand(input CLICommand, storage *[]TaskObj) error {
	fmt.Printf("CLITodo Tracker by ngixx\n\n")
	format := "%v: %v\n" //name, description
	for name, description := range helpDescription {
		fmt.Printf(format, name, description)
	}
	return nil
}

var availableCommands = map[string]func(CLICommand, *[]TaskObj) error{
	"add":    addTask,
	"update": updateDescription,
	"delete": deleteTaskById,
	"mark":   updateStatus,
	"list":   printTasks,
	"help":   helpCommand,
}

var helpDescription = map[string]string{
	"add [content]": "creates a new task",
	"delete [id]":   "deletes tasks by index",
	"mark [done | in-progress | todo] [index]": "sets the status of the task",
	"update [index] [content]":                 "updates the description of task",
	"list [filter]":                            "displays all tasks with a certain status, by default all of them",
}

var stringToStatus = map[string]StatusType{
	"done":        DONE_STATUS,
	"todo":        TODO_STATUS,
	"in-progress": IN_PROGRESS_STATUS,
}

type CLICommand struct {
	cmd       string
	arguments []string
}

const defaultJsonPath = "task.json"

func Run() {
	tasks, loadedErr := loadTasksFromJson(defaultJsonPath)
	if loadedErr != nil {
		fmt.Printf("Creating new json storage\n")
	}

	defer func() {
		if err := exportTasksToJson(defaultJsonPath, tasks); err != nil {
			fmt.Printf("%v", err)
		}
	}()

	command := getCommand()

	exec, exist := availableCommands[command.cmd]
	if !exist {
		fmt.Printf("Unknow command. Enter help to see avialable commands")
		return
	}

	if err := exec(command, &tasks); err != nil {
		fmt.Printf("%v", err)
	}
}
