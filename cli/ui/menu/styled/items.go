package styled

type index int

func (s index) getNext(last index) index {
	if s == last {
		return first
	}
	return s + 1
}

func (s index) getPrev(last index) index {
	if s == first {
		return last
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	first index = iota
	second
	third
	forth
	fifth
	sixth
)

type Item struct {
	idx         index
	title       string
	description string
	max_idx     index
}

func NewItem(index, max index, title, description string) Item {
	return Item{
		idx:         index,
		title:       title,
		description: description,
		max_idx:     max,
	}
}

func (t *Item) Next() {
	if t.idx == t.max_idx {
		t.idx = first
	} else {
		t.idx++
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
