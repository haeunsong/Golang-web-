package ex

import "net/http"

// type Middleware func(next HandlerFunc) HandlerFunc
type Server struct {
	*router
	middlewares  []Middleware
	startHandler HandlerFunc // 체인형태로 연결된 미들웨어의 시작점
}

// Server 생성자 함수
func NewServer() *Server {
	r := &router{make(map[string]map[string]HandlerFunc)}
	s := &Server{router: r}
	s.middlewares = []Middleware{
		logHandler, recoverHandler, staticHandler, parseFormHandler, parseJsonBodyHandler, // 기본 미들웨어로 지정
	}
	return s
}

func (s *Server) Run(addr string) {
	// startHandler를 라우터 핸들러 함수로 지정
	s.startHandler = s.router.handler()

	// 등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	// 웹 서버 시작
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Context 생성
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}
	// URL의 쿼리 매개변수를 Context에 담은 후 startHandler로 제어권을 넘겨준다.
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	s.startHandler(c)
}

// 커스텀 미들웨어
func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}
