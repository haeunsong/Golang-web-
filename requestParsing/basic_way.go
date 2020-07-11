// 1. No request parsing (가장 기본적인 방법)

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

func (d database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range d {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

// https://localhost:8000/
  http://localhost:8000/abc
  http://localhost:8000/xyz 세 개 모두 정상작동.
// db 라는 변수에 임시로 사용할 map타입의 database 정의
// ListenAndServe는 두 번째 인자로 ServeHTTP(Handler타입)를 받는다.
// cf. type Handler interface{
          ServeHTTP(ResponseWriter,*Request)
        }

