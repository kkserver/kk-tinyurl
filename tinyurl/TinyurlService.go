package tinyurl

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kkserver/kk-cache/cache"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/json"
	"time"
)

type TinyurlService struct {
	app.Service
	Get    *TinyurlTask
	Create *TinyurlCreateTask
}

func (S *TinyurlService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func IdToHash(id int64) string {

	var b = bytes.NewBuffer(nil)

	for id != 0 {

		var v = id % 62

		if v < 26 {
			b.Write([]byte{byte(v + 65)})
		} else if v < 36 {
			b.Write([]byte{byte(v - 26 + 48)})
		} else {
			b.Write([]byte{byte(v - 36 + 97)})
		}

		id = id / 62
	}

	return b.String()
}

func HashToId(hash string) int64 {

	var b = []byte(hash)
	var id int64 = 0
	var count = len(b)

	for i := count - 1; i >= 0; i-- {
		var v = b[i]
		if v >= 65 && v < 65+26 {
			id = id*62 + int64(v-65)
		} else if v >= 48 && v < 48+10 {
			id = id*62 + int64(v-48+26)
		} else if v >= 97 && v < 97+26 {
			id = id*62 + int64(v-97+36)
		}
	}

	return id
}

func URLToKey(url string) string {
	var v = md5.New()
	v.Write([]byte(url))
	return hex.EncodeToString(v.Sum(nil))
}

func (S *TinyurlService) HandleTinyurlCreateTask(a *TinyurlApp, task *TinyurlCreateTask) error {

	if task.Url == "" {
		task.Result.Errno = ERROR_TINYURL_NOT_FOUND_URL
		task.Result.Errmsg = "Not found url"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TINYURL
		task.Result.Errmsg = err.Error()
		return nil
	}

	var prefix = a.DB.Prefix

	var key = URLToKey(task.Url)

	var get = func() (*Tinyurl, error) {

		var v = Tinyurl{}
		var scaner = kk.NewDBScaner(&v)

		rows, err := kk.DBQuery(db, &a.TinyurlTable, prefix, " WHERE `key`=?", key)

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		if rows.Next() {

			err = scaner.Scan(rows)

			if err != nil {
				return nil, err
			}

			return &v, nil

		}

		return nil, nil
	}

	v, err := get()

	if err != nil {
		task.Result.Errno = ERROR_TINYURL
		task.Result.Errmsg = err.Error()
		return nil
	}

	if v == nil {

		var vv = Tinyurl{}
		vv.Url = task.Url
		vv.Key = key
		vv.Ctime = time.Now().Unix()

		_, err := kk.DBInsert(db, &a.TinyurlTable, prefix, &vv)

		if err != nil {

			v, err = get()

			if err != nil {
				task.Result.Errno = ERROR_TINYURL
				task.Result.Errmsg = err.Error()
				return nil
			} else if v != nil {
				task.Result.Tinyurl = &vv
				task.Result.Hash = IdToHash(vv.Id)
			} else {
				task.Result.Errno = ERROR_TINYURL_NOT_FOUND_URL
				task.Result.Errmsg = "Not found url"
			}

		} else {

			task.Result.Tinyurl = &vv
			task.Result.Hash = IdToHash(vv.Id)

			{
				var cache = cache.CacheSetTask{}
				cache.Key = fmt.Sprintf("tinyurl.%d", vv.Id)
				cache.Expires = 6
				b, _ := json.Encode(&vv)
				cache.Value = string(b)

				app.Handle(a, &cache)
			}
		}

	} else {
		task.Result.Tinyurl = v
		task.Result.Hash = IdToHash(v.Id)
	}

	return nil
}

func (S *TinyurlService) HandleTinyurlTask(a *TinyurlApp, task *TinyurlTask) error {

	if task.Hash == "" && task.Id == 0 {
		task.Result.Errno = ERROR_TINYURL_NOT_FOUND_HASH
		task.Result.Errmsg = "Not found hash"
		return nil
	}

	if task.Hash != "" {
		task.Id = HashToId(task.Hash)
	}

	{
		var cache = cache.CacheTask{}
		cache.Key = fmt.Sprintf("tinyurl.%d", task.Id)
		var err = app.Handle(a, &cache)
		if err == nil && cache.Result.Errno == 0 && cache.Result.Value != "" {
			var vv = Tinyurl{}
			err = json.Decode([]byte(cache.Result.Value), &vv)
			if err == nil {
				task.Result.Tinyurl = &vv
				task.Result.Hash = IdToHash(vv.Id)
				return nil
			}
		}
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TINYURL
		task.Result.Errmsg = err.Error()
		return nil
	}

	var prefix = a.DB.Prefix

	if task.Id != 0 {

		var v = Tinyurl{}
		var scaner = kk.NewDBScaner(&v)

		rows, err := kk.DBQuery(db, &a.TinyurlTable, prefix, " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_TINYURL
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		if rows.Next() {

			err = scaner.Scan(rows)

			if err != nil {
				task.Result.Errno = ERROR_TINYURL
				task.Result.Errmsg = err.Error()
				return nil
			}

			task.Result.Tinyurl = &v
			task.Result.Hash = IdToHash(v.Id)

			{
				var cache = cache.CacheSetTask{}
				cache.Key = fmt.Sprintf("tinyurl.%d", v.Id)
				cache.Expires = 6
				b, _ := json.Encode(&v)
				cache.Value = string(b)
				app.Handle(a, &cache)
			}

		} else {
			task.Result.Errno = ERROR_TINYURL_NOT_FOUND_URL
			task.Result.Errmsg = "Not found url"
		}

	} else {
		task.Result.Errno = ERROR_TINYURL_NOT_FOUND_URL
		task.Result.Errmsg = "Not found url"
	}

	return nil
}
