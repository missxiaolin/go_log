package main

import (
	"fmt"
	"unsafe"
)

type TestStructTobytes struct {
	data int64
	typeTest int64
}

//其中addr为数值的地址，len为当地数值的长度，cap为数值的容量。
//转换的时候，需要定义一个和[]byte底层结构一致的struct（如例子中的SliceMock），然后把结构体的地址赋给addr，
type SliceMock1 struct {
	addr uintptr
	len int
	cap int
}

// struct 和 byte互转
func main()  {
	var testStruct = &TestStructTobytes{100,200}
	fmt.Println(testStruct)
	Len := unsafe.Sizeof(*testStruct)
	fmt.Println(Len)
	testBytes := &SliceMock1{
		addr: uintptr(unsafe.Pointer(testStruct)),
		cap: int(Len),
		len: int(Len),
	}
	data := *(*[]byte)(unsafe.Pointer(testBytes))
	fmt.Println("[]byte is : ", data)

	var ptestStruct *TestStructTobytes = *(**TestStructTobytes)(unsafe.Pointer(&data))
	fmt.Println("ptestStruct.data is : ", ptestStruct)
}
