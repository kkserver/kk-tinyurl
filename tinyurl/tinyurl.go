package tinyurl

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

type Tinyurl struct {
	Id    int64  `json:"id,string"`
	Key   string `json:"key"`
	Url   string `json:"url"`
	Ctime int64  `json:"ctime"` //创建时间
}

type TinyurlApp struct {
	app.App
	DB           *app.DBConfig
	Tinyurl      *TinyurlService
	Remote       *remote.Service
	Client       *client.Service
	ClientCache  *client.WithService
	TinyurlTable kk.DBTable
}

func (C *TinyurlApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}
