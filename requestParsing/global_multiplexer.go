// cf. https://dejavuqa.tistory.com/314
// 4. Global multiplexer

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

func (d database) foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "foo: %s\n", d["foo"])
}

func (d database) bar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bar: %s\n", d["bar"])
}

func (d database) baz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "baz: %s\n", d["baz"])
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	}
	http.HandleFunc("/foo", db.foo)
	http.HandleFunc("/bar", db.bar)
	http.HandleFunc("/baz", db.baz)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

// ServeMux를 각 패키지 별로 만드는 대신 Go의 Global DefaultServeMux 사용
// ListenAndServe의 두 번째 인자에 nil이 들어감!!
