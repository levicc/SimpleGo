package main

import (
	"SimpleGo/controllers"
	"SimpleGo/simplego"
	"fmt"
)

type testInterface interface {
	TestA(*testReflect, string)
}

type testReflect struct {
	A string
	b string
}

func (test *testReflect) TestA(s *testReflect, a string) {
	fmt.Println(s.A, a)
}

func (test testReflect) TestB() {
	fmt.Println("B")
}

func main() {
	simplego.Add("/", &controllers.LYLoginController{})
	// simplego.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "hellor world! Get")
	// })
	simplego.Add("/aaa/bbb/:id([\\w]+)/:username([1-9]+)", &controllers.LYLoginController{})
	simplego.Run()

	// indirect1 := reflect.Indirect(reflect.ValueOf(testReflect{}))
	// indirect2 := reflect.Indirect(reflect.ValueOf(&testReflect{}))

	// fmt.Println(indirect1, indirect2, indirect1.Type(), indirect2.Type(), indirect1.NumMethod(), indirect2.NumMethod())
	// fmt.Println(reflect.New(indirect1.Type()).NumMethod(), reflect.New(indirect2.Type()).NumMethod())

	// loginType := reflect.TypeOf(&testReflect{})
	// fmt.Println(loginType, loginType.NumMethod())

	// loginValue := reflect.ValueOf(&testReflect{})
	// fmt.Println(loginValue, loginValue.NumMethod())

	// fmt.Println(loginValue.Type(), loginValue.Type().NumMethod())

	// //new之后在原来的基础上返回指针，eLem()就是取指针指向的值
	// newLogin1 := reflect.New(loginType)
	// fmt.Println(newLogin1, newLogin1.NumMethod())

	// newLogin2 := reflect.New(loginType).Elem()
	// fmt.Println(newLogin2, newLogin2.NumMethod())

	// method := newLogin2.MethodByName("TestA")
	// fmt.Println(method)

	// test := &testReflect{A: "luyang"}
	// method.Call([]reflect.Value{reflect.ValueOf(test), reflect.ValueOf("hello")})

	//testTwoMethods(&testReflect{})
}

// func testTwoMethods(a testInterface) {
// 	typea := reflect.TypeOf(a)
// 	value := reflect.New(typea).Elem()
// 	method := value.MethodByName("TestA")
// 	test := &testReflect{A: "luyang"}
// 	method.Call([]reflect.Value{reflect.ValueOf(test), reflect.ValueOf("hello")})
// }
