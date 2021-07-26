package common

import (
	"encoding/json"
	"flightchain/model"
	"reflect"
	"regexp"
	"unsafe"
)

type ByteStruct struct {
	addr uintptr
	len  int
	cap  int
}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func StructToByte(stru model.Aircraft) []byte {
	len := unsafe.Sizeof(stru)
	bytes := ByteStruct{
		addr: uintptr(unsafe.Pointer(&stru)),
		cap:  int(len),
		len:  int(len),
	}
	data := *(*[]byte)(unsafe.Pointer(&bytes))
	return data
}

func ByteToStruct(bt []byte) model.Aircraft {
	var aircraft *model.Aircraft = *(**model.Aircraft)(unsafe.Pointer(&bt))
	return *aircraft
}

// 结构体转为string
func StructToString(i interface{}) string {

	tmp := Struct2Map(i)
	str, err := json.Marshal(tmp)
	if err != nil {
		return err.Error()
	}
	return string(str)

}

// 结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
