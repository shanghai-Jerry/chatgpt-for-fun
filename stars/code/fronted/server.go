package main

import (
	"log"
	"net/http"
)

func main() {
	// 设置静态文件服务
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// 启动服务器
	log.Println("服务器启动在端口 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
