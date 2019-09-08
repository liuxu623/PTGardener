package main

import "fmt"

//ErrConfig 用来控制一些参数
type ErrConfig struct {
	PrintTag bool
	TagNum   int
}

//Tag 大写比较好看
func (e ErrConfig) Tag(n int) ErrConfig {
	e.TagNum = n
	return e
}
func (e *ErrConfig) init() {
	e.PrintTag = true
}
func errPrint(err error, E ...ErrConfig) bool {
	if err != nil {
		if len(E) > 0 && E[0].PrintTag {
			if E[0].TagNum != 0 {
				fmt.Println("↘↘ ============>", E[0].TagNum, "<=========== ↙↙")
			}
		}
		fmt.Println(err)
		return true
	}
	return false
}
func errPainc(err error) {
	if err != nil {
		panic(err)
	}
}
func doNotPanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

var e ErrConfig
