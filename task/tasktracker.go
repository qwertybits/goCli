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
	err = os.WriteFile(path, bytes, os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

func Run() {

	tasks := make([]TaskObj, 0) //storage of tasks

	loadFromJson(defaultJsonPath, &tasks)

	tasks = append(tasks, NewTask(0, "hello world"))
	exportToJson(defaultJsonPath, tasks)

	fmt.Printf("%v", tasks)
}
