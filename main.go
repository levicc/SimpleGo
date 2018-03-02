package main

import (
	"SimpleGo/controllers"
	"SimpleGo/simplego"
	"fmt"
	"html/template"
	"net/http"
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
	simplego.SetStaticPath("/asset", "uploads")

	simplego.Add("/", &controllers.LYLoginController{})
	simplego.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		if t, error := template.ParseFiles("views/login.html"); error == nil {
			t.Execute(w, nil)
		}
	})
	simplego.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hellor world!")
	})
	simplego.Add("/aaa/bbb/:id([\\w]+)/:username([1-9]+)", &controllers.LYLoginController{})
	simplego.Run()

	//指针初始化为零值指针，这个是问题的根源
	//https://stackoverflow.com/questions/43072451/reflect-new-returns-nil-instead-of-initialized-struct
	// a := new(*testReflect)
	// fmt.Println(reflect.ValueOf(a).Elem())

	// indirect1 := reflect.Indirect(reflect.ValueOf(testReflect{A: "luyang"}))
	// indirect2 := reflect.Indirect(reflect.ValueOf(&testReflect{A: "luyang"}))

	// fmt.Println(indirect1, indirect2, indirect1.Type(), indirect2.Type(), indirect1.NumMethod(), indirect2.NumMethod())
	// fmt.Println(reflect.New(indirect1.Type()).NumMethod(), reflect.New(indirect2.Type()).NumMethod())

	// loginType1 := reflect.TypeOf(&testReflect{A: "luyang"})
	// loginType2 := reflect.TypeOf(testReflect{A: "luyang"})
	// fmt.Println(loginType1, loginType2, loginType1.NumMethod(), loginType2.NumMethod())

	// loginValue1 := reflect.ValueOf(&testReflect{A: "luyang"})
	// loginValue2 := reflect.ValueOf(testReflect{A: "luyang"})
	// fmt.Println(loginValue1, loginValue2, loginValue1.NumMethod(), loginValue2.NumMethod())

	// fmt.Println(loginValue1.Type(), loginValue2.Type(), loginValue1.Type().NumMethod(), loginValue2.Type().NumMethod())

	// fmt.Println(loginValue1.Elem(), loginValue1.Elem().NumMethod())
	// //fmt.Println(loginValue2.Elem(), loginValue2.Elem().NumMethod())
	// fmt.Println(reflect.New(loginType1).Elem(), reflect.New(loginType1).Elem().NumMethod(), reflect.New(loginType2).Elem())
	// fmt.Println(reflect.New(indirect1.Type()), reflect.New(indirect2.Type()), reflect.New(indirect1.Type()).Type(), reflect.New(indirect2.Type()).Type())

	// bugValue1 := reflect.New(loginType1) //type.elem()也可以和上面redirect的效果一样
	// bugValue2 := reflect.New(loginType2)
	// fmt.Println(bugValue1, bugValue1.NumMethod(), bugValue1.Elem(), bugValue1.Elem().NumMethod(), bugValue1.Type(), bugValue1.Elem().Type())
	// fmt.Println(bugValue2, bugValue2.NumMethod(), bugValue2.Elem(), bugValue2.Elem().NumMethod(), bugValue2.Type(), bugValue2.Elem().Type())

	// fmt.Println(loginValue1, loginValue1.Interface())
	//fmt.Println(bugValue2.Elem().Elem()) //出错，struct的elem()报错
	//fmt.Printf("%x", &bugValue2)
	//new之后在原来的基础上返回指针，eLem()就是取指针指向的值
	// newLogin1 := reflect.New(loginType)
	// fmt.Println(newLogin1, newLogin1.NumMethod())

	// newLogin2 := reflect.New(loginType).Elem()
	// fmt.Println(newLogin2, newLogin2.NumMethod())

	// method := newLogin2.MethodByName("TestA")
	// fmt.Println(method)

	// test := &testReflect{A: "luyang"}
	// method.Call([]reflect.Value{reflect.ValueOf(test), reflect.ValueOf("hello")})

	// testTwoMethods(&testReflect{})
}

// func testTwoMethods(a testInterface) {
// 	typea := reflect.TypeOf(a)
// 	value := reflect.New(typea).Elem()
// 	method := value.MethodByName("TestA")
// 	test := &testReflect{A: "luyang"}
// 	method.Call([]reflect.Value{reflect.ValueOf(test), reflect.ValueOf("hello")})
// }
