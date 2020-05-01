package serverRoom

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/unknwon/com"
)

//配置文件操作

type SingleConfig struct {
	Debug bool              `json:"debug"`
	Tcp   map[string]string `json:"tcp"`
}

var v SingleConfig
var once sync.Once

func GetConfigInstance() SingleConfig {
	once.Do(load)
	return v
}

func (conf SingleConfig) GetStringMapString() map[string]string {
	load()
	return v.Tcp
}

//配置文件初始化
func load() {
	confPath := Arg.configfile
	if !com.IsFile(confPath) {
		log.Fatalln(confPath + " not exists")
	}
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if err = json.Unmarshal(data, &v); err != nil {
		log.Fatalln(err)
	}
}

type ConfResponse struct {
	Action string
	Key    string
	Value  interface{}
	Error  error
}

func ConfWatch(stop chan struct{}) <-chan *ConfResponse {
	respChan := make(chan *ConfResponse, 10)

	go func() {
		//inode
		watcher, err := fsnotify.NewWatcher()
		//监视配置文件inode 出错了,退出程序
		if err != nil {
			panic(err)
		}

		watcher.Add(Arg.configfile)

		go func() {
			<-stop
			watcher.Close()
		}()

		respdata := &ConfResponse{
			Error: nil,
		}

		for {
			select {
			case event := <-watcher.Events:
				//fmt.Println(event)
				if event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename ||
					event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					watcher.Remove(Arg.configfile)
					watcher.Add(Arg.configfile)

					//需要读取配置文件
					//通过chan通知
					respChan <- respdata
				}

			case err := <-watcher.Errors:
				respdata.Error = err
				respChan <- respdata
			}

		}
	}()

	return respChan
}
