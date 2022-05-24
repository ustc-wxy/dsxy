package objects

import (
	"dsxy/apiServer/heartbeat"
	"dsxy/apiServer/locate"
	"dsxy/es"
	"dsxy/objectstream"
	"dsxy/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//func put(w http.ResponseWriter, r *http.Request) {
//	object := strings.Split(r.URL.EscapedPath(), "/")[2]
//	c, e := storeObject(r.Body, object)
//	if e != nil {
//		log.Println(e)
//	}
//	w.WriteHeader(c)
//}

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header) //从头部digest获得hash

	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header) //从头部content-length获得size
	c, e := storeObject(r.Body, hash, size)

	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	rTemp := io.TeeReader(r, stream)
	d := utils.CalculateHash(rTemp)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

//func storeObject(r io.Reader, object string) (int, error) {
//	stream, e := putStream(object)
//	if e != nil {
//		return http.StatusServiceUnavailable, e
//	}
//	io.Copy(stream, r)
//	e = stream.Close()
//	if e != nil {
//		return http.StatusInternalServerError, e
//	}
//	return http.StatusOK, nil
//}

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectstream.NewTempPutStream(server, hash, size)
}
