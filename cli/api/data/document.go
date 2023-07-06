package data

import (
	"io/ioutil"
	"os"
)

type Document struct {
	ID       int    `json:"related_req_id"`
	Filename string `json:"title"`
	Content  string `json:"content"`
}

func NewDocument(id int, filename string) (doc Document, err error) {
	content, err := readFile(filename)
	doc = Document{
		ID:       id,
		Filename: filename,
		Content:  content,
	}
	return
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func readFile(filename string) (content string, err error) {
	if !fileExists(filename) {
		return "", os.ErrNotExist
	}
	data, err := ioutil.ReadFile(filename)
	content = string(data)

	return
}
