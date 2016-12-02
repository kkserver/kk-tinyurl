package tinyurl

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TinyurlCreateTaskResult struct {
	app.Result
	Hash    string   `json:"hash,omitempty"`
	Tinyurl *Tinyurl `json:"tinyurl,omitempty"`
}

type TinyurlCreateTask struct {
	app.Task
	Url    string `json:"url"`
	Result TinyurlCreateTaskResult
}

func (task *TinyurlCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *TinyurlCreateTask) GetInhertType() string {
	return "tinyurl"
}

func (task *TinyurlCreateTask) GetClientName() string {
	return "Tinyurl.Create"
}
