package ex

import (
	"net/http"
	"strings"
)

// 라우터에 핸들러를 등록하기 위한 메서드인 HandleFunc 만들기

type router struct {
	// 키 : http 메서드
	// 값 : URL 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]http.HandlerFunc // 이차원 맵
}

// 브라우저가 서버로 요청을 만들면 서버는 해당요청을 처리하고 응답한다.
// 이러한 응답/패턴을 추상화한 것이 바로 Handler 인터페이스다.
type Handler interface {
	// ServeHttp 메소드가 응답헤더와 데이터를 ResponseWriter로 보내는 역할을 한다.
	ServeHTTP(http.ResponseWriter, *http.Request)
}

//매개변수로 전달된 http메서드, URL패턴, 핸들러 함수를 2차원 맵인 handlers필드에 등록한다.
func (r *router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새로운 맵 생성
		m = make(map[string]http.HandlerFunc)
		r.handlers[method] = m
	}
	// http 메서드로 등록된 맵에 URL패턴과 핸들러 함수 등록
	m[pattern] = h
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m, ok := r.handlers[req.Method]; ok {
		if h, ok := m[req.URL.Path]; ok {
			//요청 URL에 해당하는 핸들러 수행
			h(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

// 라우터에 등록된 동적 URL패턴과 실제 URL경로가 일치하는지 확인하는 match 함수
func match(pattern, path string) (bool, map[string]string) {
	// 패턴과 패스가 정확히 일치하면 바로 true 반환
	if pattern == path {
		return true, nil
	}

	// 패턴과 패스를 "/" 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 패턴과 패스를 "/"로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 URL 매개변수를 담기 위한 params 맵 생성
	params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// 패턴이 ':' 문자로 시작하면 params에 URL params 를 담은 후 다음 루프 수행
		default:
			// 일치하는 경우가 없으면 false 반환
			return false, nil
		}
	}

	// true와 params를 반환
	return true, params

}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// http메서드에 맞는 모든 handlers를 반복하여 요청 URL에 해당하는 handler를 찾음
	for pattern, handler := range r.handlers[req.Method] {
		if ok, _ := match(pattern, req.URL.Path); ok {
			// 요청 URL에 해당하는 handler 수행
			handler(w, req)
			return
		}
	}
	// 요청 URL에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
	http.NotFound(w, req)
	return
}
