package task

import (
	"errors"
	"fmt"
	"time"
)

type StatusType int

const (
	ANY_STATUS StatusType = iota
	DONE_STATUS
	TODO_STATUS
	IN_PROGRESS_STATUS
)

func (s StatusType) String() string {
	switch s {
	case DONE_STATUS:
		return "done"
	case TODO_STATUS:
		return "todo"
	case IN_PROGRESS_STATUS:
		return "in-progress"
	}
	return "any"
}

const printFormat string = "%v: %v | %v\nCreated: %v\tUpdated: %v" //id, content, status, cdate, udate

type TaskObj struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      StatusType `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func NewTask(id int, content string) TaskObj {
	status := TODO_STATUS
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

func (t *TaskObj) SetStatus(status StatusType) {
	t.update()
	t.Status = status
}

func (t TaskObj) GetStatus() StatusType {
	return t.Status
}

func (t TaskObj) IsSameStatus(s StatusType) bool {
	if s == ANY_STATUS {
		return true
	}
	return t.GetStatus() == s
}

func (t TaskObj) String() string {
	return fmt.Sprintf(printFormat, t.Id, t.Description, t.Status,
		t.CreatedAt.Format("2006-01-02 15:04"), t.UpdatedAt.Format("2006-01-02 15:04"))
}

func (t *TaskObj) SetId(id int) {
	t.Id = id
}
