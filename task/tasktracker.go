package task

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultJsonPath = "task.json"

func loadFromJson(path string, storage *[]TaskObj) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, storage)
	if err != nil {
		return err
	}
	return nil
}

func exportToJson(path string, storage []TaskObj) error {
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

func addTask(id int, content string, storage *[]TaskObj) {
	*storage = append(*storage, NewTask(id, content))
}

func printTasks(storage *[]TaskObj) {
	for _, task := range *storage {
		fmt.Printf("%v\n", task)
	}
}

func Run() {

	tasks := make([]TaskObj, 0) //storage of tasks

	err := loadFromJson(defaultJsonPath, &tasks)

	if err != nil {
		fmt.Printf("%v", err)
	}

	printTasks(&tasks)

	err = exportToJson(defaultJsonPath, tasks)
	if err != nil {
		fmt.Printf("%v", err)
	}

}
