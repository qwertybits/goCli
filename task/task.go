package task

import (
	"errors"
	"fmt"
	"time"
)

const printFormat string = "%v: %v | %v\nCreated: %v\tUpdated: %v" //id, content, status, cdate, udate

type TaskObj struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewTask(id int, content string) TaskObj {
	status := "todo"
	createdAt := time.Now()
	updatedAt := createdAt
	return TaskObj{id, content, status, createdAt, updatedAt}
}

func (t *TaskObj) update() {
	t.UpdatedAt = time.Now()
}

func (t *TaskObj) SetDescription(content string) error {
	if len(content) == 0 {
		return errors.New("description cant be empty")
	}
	t.Description = content
	t.update()
	return nil
}

func (t *TaskObj) SetStatus(status string) error {
	if status != "todo" && status != "in-progress" && status != "done" {
		msg := fmt.Sprintf("unknow status: %v", status)
		return errors.New(msg)
	}
	t.update()
	t.Status = status
	return nil
}

func (t TaskObj) String() string {
	return fmt.Sprintf(printFormat, t.Id, t.Description, t.Status,
		t.CreatedAt.Format("2006-01-02 15:04"), t.UpdatedAt.Format("2006-01-02 15:04"))
}

func (t *TaskObj) SetId(id int) {
	t.Id = id
}
