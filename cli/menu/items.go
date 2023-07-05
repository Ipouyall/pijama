package menu

type status int

func (s status) getNext() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	todo status = iota
	inProgress
	done
)

type Item struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Item {
	return Item{status: status, title: title, description: description}
}

func (t *Item) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Item) FilterValue() string {
	return t.title
}

func (t Item) Title() string {
	return t.title
}

func (t Item) Description() string {
	return t.description
}
