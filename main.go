package main

import (
	"GeeCache/geecache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom": "630",
	"Jack": "678",
	"Sam": "567",
}

func main()  {
	//geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
	//	func(key string) ([]byte, error) {
	//		log.Println("[SlowDB] search key", key)
	//		if v, ok := db[key];ok {
	//			return []byte(v), nil
	//		}
	//		return nil, fmt.Errorf("%s not exist", key)
	//	}))
	//addr := "localhost:9999"
	//peers := geecache.NewHTTPPool(addr)
	//log.Println("geecache is running at", addr)
	//log.Fatal(http.ListenAndServe(addr, peers))
	var port int
	var api bool

	flag.IntVar(&port, "port",8001, "GeeCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001 : "http://localhost:8001",
		8002 : "http://localhost:8002",
		8003 : "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap{
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}

//创建新的分组
func createGroup() *geecache.Group {
	return geecache.NewGroup("score", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key];ok{
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

//启动缓存服务
func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	//将peers注册到group中
	gee.RegisterPeers(peers)
	log.Println("geeecache is running at", addr)
	//todo
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			key := request.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/octet-stream")
			writer.Write(view.ByteSlice())
		}))
	log.Println("fronted server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}