package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return result, err
	}
	err = json.Unmarshal(buffer[:n], &result)
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
		newFilter, exist := strToStatus[input.arguments[0]]
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

	newStatus, exist := strToStatus[input.arguments[0]]
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

func validateId(id int, storage *[]TaskObj) bool {
	return !(id >= len(*storage) || id < 0)
}

var availableCommands = map[string]func(CLICommand, *[]TaskObj) error{
	"add":    addTask,
	"update": updateDescription,
	"delete": deleteTaskById,
	"mark":   updateStatus,
	"list":   printTasks,
}

var strToStatus = map[string]StatusType{
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

	command, err := getCommand()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	exec, exist := availableCommands[command.cmd]
	if !exist {
		fmt.Printf("Unknow command")
		return
	}

	if err := exec(command, &tasks); err != nil {
		fmt.Printf("%v", err)
	}
}
