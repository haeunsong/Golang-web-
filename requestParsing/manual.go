// cf. https://dejavuqa.tistory.com/314
// 2. Manual request parsing 

package main

import (
	"fmt"
	"log"
	"net/http"
)

type pounds float32

func (p pounds) String() string {
	return fmt.Sprintf("£%.2f", p)
}

type database map[string]pounds

// 요청되는 경로 구분하여 처리
func (d database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/foo":
		fmt.Fprintf(w,"foo: %s\n", d["foo"])
	case "/bar":
		fmt.Fprintf(w,"bar: %s\n", d["foo"])
	default:
		w.WriteHeader(http.StatusNotFound) // 지정된 경로가 아닐 경우 404 메시지 보냄(WriteHeader()사용)
		fmt.Fprintf(w, "No page found for: %s\n", r.URL)
	}	
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

// localhost:8000/foo
// localhost:8000/bar 정상접속
