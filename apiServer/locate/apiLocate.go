package locate

import (
	"dsxy/rabbitmq"
	"dsxy/rs"
	"dsxy/types"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])

	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}
func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	var info types.LocalMessage
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return locateInfo
}

//func Locate(name string) string {
//	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
//	q.Publish("dataServers", name)
//	c := q.Consume()
//	go func() {
//		time.Sleep(time.Second)
//		q.Close()
//	}()
//	msg := <-c
//	s, _ := strconv.Unquote(string(msg.Body))
//	return s
//}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}

//func Exist(name string) bool {
//	return Locate(name) != ""
//}
