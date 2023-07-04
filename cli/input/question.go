package input

type Question struct {
	question string
	answer   string
	input    Input
}

func newQuestion(q string) Question {
	return Question{question: q}
}

func NewShortQuestion(q string, p string) Question {
	question := newQuestion(q)
	model := NewShortAnswerField(p)
	question.input = model
	return question
}

func NewLongQuestion(q string, p string) Question {
	question := newQuestion(q)
	model := NewLongAnswerField(p)
	question.input = model
	return question
}
