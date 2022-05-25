package heartbeat

import (
	"dsxy/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

func ChooseRandomDataServer(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	use := make(map[string]bool)
	for _, addr := range exclude {
		use[addr] = true
	}
	servers := GetDataServers()
	for _, s := range servers {
		_, ok := use[s]
		if !ok {
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		return nil
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])

	}
	return ds
}

//func ChooseRandomDataServer() string {
//	ds := GetDataServers()
//	n := len(ds)
//	if n == 0 {
//		return ""
//	}
//	return ds[rand.Intn(n)]
//}
