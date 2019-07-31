// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command server is a test server for the Autobahn WebSockets Test Suite.
package main

import (
	"./sdk"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Status int

func echoReadAll(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	action := r.Form.Get("action")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	if action == "" {
		fmt.Fprintf(w, `{"ret":0, "msg":"无效参数", "data":""}`)
	} else {
		ret := handleMessage([]byte(action))
		fmt.Fprintf(w, string(ret))
	}

}
func echoCaptureAll(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Println("Upgrade:", err)
		return
	} else {
		if sdk.Camera == -1 {
			sdk.Camera = 0
		}
		go sdk.OpenCamera()
		//发送视频图片
		for {
			if sdk.Image == "" {
				err = conn.WriteMessage(1, []byte(`{"ret":0, "msg":"视频捕获失败", "data":""}`))
				if err != nil {
					log.Println("VideoCapture:", err)
					return
				}
				continue
			}
			err = conn.WriteMessage(1, []byte(`{"ret":1, "msg":"--VideoCapture", "data":"data:image/jpg;base64, `+sdk.Image+`"}`))
			if err != nil {
				log.Println("VideoCapture:", err)
				return
			}
			if sdk.Camera == -1 {
				return
			}
		}
	}

	for {
		mt, b, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			return
		}
		if mt == websocket.TextMessage {
			if !utf8.Valid(b) {
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
					time.Time{})
				log.Println("ReadAll: invalid utf8")
			}
		}
	}

}

//消息处理
func handleMessage(message []byte) []byte {
	if string(message) == "--ReadCard" {
		return sdk.ReadCard()
	} else if string(message) == "--CloseCamera" {
		go sdk.CloseCamera()
		return []byte(`{"ret":1, "msg":"请求关闭摄像头", "data": ""}`)
	} else if string(message) == "--Close" {
		//服务关闭
		Status = -1
		//摄像头关闭
		sdk.Camera = -1
		//清空摄像头内容
		sdk.Image = ""
		os.Exit(0)
		s, _ := exec.LookPath(os.Args[0])
		i := strings.LastIndex(s, "\\")
		path := string(s[i+1:])
		c := exec.Command("taskkill.exe", "/f", "/im", path)
		c.Start()
		return []byte(`{"ret":0, "msg":"服务关闭", "data": ""}`)
	} else {
		return []byte(`{"ret":0, "msg":"未知指令", "data": ""}`)
	}
}

var addr = flag.String("addr", ":9000", "http service address")

func main() {
	flag.Parse()
	http.HandleFunc("/", echoReadAll)
	http.HandleFunc("/capture", echoCaptureAll)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
