package locate

import (
	"dsxy/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

//func Locate(name string) bool {
//	_, err := os.Stat(name)
//	return !os.IsNotExist(err)
//}

func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	return ok
}

func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}
func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		//fmt.Println("data Locate debug,msg is", string(msg.Body))
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		//if Locate("." + os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
		//	fmt.Println("locate success!")
		//	q.Send(msg.ReplyTo,
		//		os.Getenv("LISTEN_ADDRESS"))
		//}
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/object/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
