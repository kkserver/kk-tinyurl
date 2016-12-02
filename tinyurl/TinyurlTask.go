package tinyurl

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TinyurlTaskResult struct {
	app.Result
	Hash    string   `json:"hash,omitempty"`
	Tinyurl *Tinyurl `json:"tinyurl,omitempty"`
}

type TinyurlTask struct {
	app.Task
	Hash   string `json:"hash,omitempty"`
	Id     int64  `json:"id,string,omitempty"`
	Result TinyurlTaskResult
}

func (task *TinyurlTask) API() string {
	return "tinyurl/get"
}

func (task *TinyurlTask) GetResult() interface{} {
	return &task.Result
}

func (task *TinyurlTask) GetInhertType() string {
	return "tinyurl"
}

func (task *TinyurlTask) GetClientName() string {
	return "Tinyurl.Get"
}
