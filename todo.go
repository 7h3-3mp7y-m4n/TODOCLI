package todocli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type Item struct {
	Task      string     `json:"task"`
	Result    bool       `json:"result"`
	Created   time.Time  `json:"created"`
	Completed *time.Time `json:"completed,omitempty"`
}

type Todos []Item

func (t *Todos) Add(task string) {
	todo := Item{
		Task:      task,
		Result:    false,
		Created:   time.Now(),
		Completed: nil,
	}
	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	if err := t.validateIndex(index); err != nil {
		return err
	}
	now := time.Now()
	(*t)[index-1].Completed = &now
	(*t)[index-1].Result = true
	return nil
}

func (t *Todos) Delete(index int) error {
	if err := t.validateIndex(index); err != nil {
		return err
	}
	*t = append((*t)[:index-1], (*t)[index:]...)
	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, t)
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) validateIndex(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index")
	}
	return nil
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "TASK"},
			{Align: simpletable.AlignCenter, Text: "RESULT"},
			{Align: simpletable.AlignCenter, Text: "CREATED_AT"},
			{Align: simpletable.AlignCenter, Text: "COMPLETED_AT"},
		},
	}
	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		task := blue(item.Task)
		result := red("no")
		completedAt := ""
		if item.Result {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			result = green("yes")
			if item.Completed != nil {
				completedAt = item.Completed.Format(time.RFC822)
			}
		}
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: result},
			{Text: item.Created.Format(time.RFC822)},
			{Text: completedAt},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending Tasks", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Result {
			total++
		}
	}
	return total
}
