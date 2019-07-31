// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command server is a test server for the Autobahn WebSockets Test Suite.
package sdk

import (
	"bytes"
	"encoding/base64"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var dll, err = syscall.LoadLibrary("Sdtapi.dll")
var Proc1, _ = syscall.GetProcAddress(dll, "InitComm")      //InitComm
var Proc2, _ = syscall.GetProcAddress(dll, "Authenticate")  //Authenticate
var Proc3, _ = syscall.GetProcAddress(dll, "CardOn")        //CardOn
var Proc4, _ = syscall.GetProcAddress(dll, "ReadBaseInfos") //ReadBaseInfos原型3
var Proc5, _ = syscall.GetProcAddress(dll, "CloseComm")     //CloseComm
//端口初始化
func InitComm() uintptr {
	InitComm, _, _ := syscall.Syscall(Proc1, 1, 1001, 0, 0)
	return InitComm
}

//端口初始化
func Authenticate() uintptr {
	Authenticate, _, _ := syscall.Syscall(Proc2, 0, 0, 0, 0)
	return Authenticate
}

//判断身份证是否在设备上
func CardOn() uintptr {
	CardOn, _, _ := syscall.Syscall(Proc3, 0, 0, 0, 0)
	return CardOn
}

var Name = make([]byte, 31)       //姓名
var Gender = make([]byte, 3)      //性别
var Folk = make([]byte, 10)       //民族
var BirthDay = make([]byte, 9)    //出生日期
var Code = make([]byte, 19)       //身份证号码
var Address = make([]byte, 71)    //地址
var Agency = make([]byte, 31)     //签证机关
var ExpireStart = make([]byte, 9) //有效期起始日期
var ExpireEnd = make([]byte, 9)   //有效期截至日期

//读卡信息
func ReadBaseInfos() uintptr {
	ReadBaseInfos, _, _ := syscall.Syscall12(Proc4, 9, uintptr(unsafe.Pointer(&Name[0])), uintptr(unsafe.Pointer(&Gender[0])), uintptr(unsafe.Pointer(&Folk[0])), uintptr(unsafe.Pointer(&BirthDay[0])), uintptr(unsafe.Pointer(&Code[0])), uintptr(unsafe.Pointer(&Address[0])), uintptr(unsafe.Pointer(&Agency[0])), uintptr(unsafe.Pointer(&ExpireStart[0])), uintptr(unsafe.Pointer(&ExpireEnd[0])), 0, 0, 0)
	return ReadBaseInfos
}

//端口关闭
func CloseComm() uintptr {
	CloseComm, _, _ := syscall.Syscall(Proc5, 0, 0, 0, 0)
	return CloseComm
}

//读身份证信息
func ReadCard() []byte {
	var ret uintptr
	if err != nil {
		return []byte(`{"ret":0, "msg":"` + err.Error() + `", "data": ""}`)
	}

	ret = InitComm()
	if ret != uintptr(1) {
		return []byte(`{"ret":0, "msg":"端口初始化失败", "data": ""}`)
	}

	ret = Authenticate()
	if ret != uintptr(1) {
		//return []byte(`{"ret":0, "msg":"卡认证失败", "data": ""}`)
	}

	ret = CardOn()
	if ret != uintptr(1) {
		return []byte(`{"ret":0, "msg":"无身份证", "data": ""}`)
	}

	ret = ReadBaseInfos()
	if ret == uintptr(0) {
		return []byte(`{"ret":0, "msg":"错误", "data": ""}`)
	} else if -ret == uintptr(4) {
		return []byte(`{"ret":0, "msg":"缺少dll", "data": ""}`)
	}

	ret = CloseComm()
	if ret != uintptr(1) {
		return []byte(`{"ret":0, "msg":"端口关闭失败", "data": ""}`)
	}
	//处理身份证信息
	name := Conversion(Name)
	gender := Conversion(Gender)
	folk := Conversion(Folk)
	birthDay := Conversion(BirthDay)
	code := Conversion(Code)
	address := Conversion(Address)
	agency := Conversion(Agency)
	expireStart := Conversion(ExpireStart)
	eexpireEnd := Conversion(ExpireEnd)
	photo := ImageBase64("photo.bmp")
	data := []byte(`{"name":"` + string(name) + `","gender":"` + string(gender) + `","folk":"` + string(folk) + `","birthDay":"` + string(birthDay) + `","code":"` + string(code) + `","address":"` + string(address) + `","agency":"` + string(agency) + `","expireStart":"` + string(expireStart) + `","eexpireEnd":"` + string(eexpireEnd) + `","photo":"` + photo + `"}`)
	return []byte(`{"ret":1, "msg":"--ReadCard", "data": ` + string(data) + `}`)
}

//gbk编码转utf-8编码并去掉\0
func Conversion(b []byte) []byte {
	ret, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	return bytes.Trim(ret, "\x00")
}

//图片转Base64
func ImageBase64(path string) string {
	//获取文件绝对路径
	s, _ := exec.LookPath(os.Args[0])
	i := strings.LastIndex(s, "\\")
	filepath := string(s[0:i+1] + path)

	image, _ := ioutil.ReadFile(filepath)
	ImageBase64 := base64.StdEncoding.EncodeToString(image)
	return "data:image/bmp;base64," + ImageBase64
}
