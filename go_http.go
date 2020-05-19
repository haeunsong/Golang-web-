
// ex> 주소창에 localhost:9090/?test=1234&test2=5678 입력
package main

import(
	"fmt"
	"log"
	"net/http"
	"strings"
)


func defaultHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	// Get 파라미터 및 정보 출력
	fmt.Println("default : ",r.Form)
	fmt.Println("path : ",r.URL.Path)
	fmt.Println("param : ",r.Form["test_param"])
	// Parameter 전체 출력
	for k,v := range r.Form{
		fmt.Println("key: ",k)
		fmt.Println("val: ",strings.Join(v,""))
	}
	// 기본 출력
	fmt.Fprintf(w, "Haeun, Golang Webserver Working!")
	
}

func main(){
	// 기본 url 핸들러 메소드 지정
	// http를 통해 /에 접근하면 이 핸들러를 작동시키겠다는 의미!
	// defaultHandler를 인수로 받음
	http.HandleFunc("/",defaultHandler)
	// 서버 시작
	err:= http.ListenAndServe(":9090",nil)
	// 예외 처리
	if err!= nil{
		log.Fatal("ListenAndServe: ",err)
	}else{
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
