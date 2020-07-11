// cf. https://dejavuqa.tistory.com/314
// 3. Multiplexer
// ServeHTTP 대신 ServeMux 사용

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

func main(){
	db := database{
		"foo":1,
		"bar":2,
		"baz":3,
	}

	mux := http.NewServeMux()

	mux.Handle("/foo",http.HandlerFunc(db.foo))
	mux.Handle("/bar",http.HandlerFunc(db.bar))

	// Convenience method for longer form mux.Handle
	mux.HandleFunc("/baz", db.baz)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

// ServeMux는 HandlerFunc를 제공하며 개발 함수를 바로 사용 가능.
// ServeMux 인스턴스는 Multiplexer이며 ListenAndServe의 두 번째 인자로 들어갈 수 있다.
// ServeMux를 사용하려면 http.NewServeMux를 선언해야함.
