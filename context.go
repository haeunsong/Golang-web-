package ex

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"
	"text/template"
)

// Context 타입 Parmas 필드에 라우터에서 해석한 URL 매개변수를 담고, 핸들러 내부에는 Context 값이 전달되게 한다.
type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// 핸들러 함수 타입으로 http.HandlerFunc 대신에 HandlerFunc 를 사용한다.
type HandlerFunc func(*Context)

// 컨텍스트에 JSON 포맷으로 데이터를 렌더링
func (c *Context) RenderJson(v interface{}) {
	// HTTP Status 를 StatusOK 로 지정
	c.ResponseWriter.WriteHeader(http.StatusOK)
	// Content-Type 을 application/json 으로 지정
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	// v값을 json으로 출력
	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		// 에러 발생 시 RenderErr 메서드 호출
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

// 컨텍스트에 XML 포맷으로 데이터를 렌더링
func (c *Context) RenderXml(v interface{}) {
	// HTTP Status 를 StatusOK로 지정
	c.ResponseWriter.WriteHeader(http.StatusOK)
	// Content-Type 을 application/json 으로 지정
	c.ResponseWriter.Header().Set("Content-Type", "application/xml;charset=utf-8")

	// v값을 xml로 출력
	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		// 에러 발생 시 RenderErr 메서드 호출
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

// 에러 상태를 적절한 HTTP Status 로 렌더링
func (c *Context) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			// 정상적인 code를 전달하면 HTTP Status를 해당 code로 지정
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			// 정상적인 code가 아니면 HTTP Status를 StatusInternalServerError로 지정
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

// templates: 템플릿 객체를 보관하기 위한 map
var templates = map[string]*template.Template{}

func (c *Context) RenderTemplate(path string, v interface{}) {
	// path에 해당하는 템플릿이 있는지 확인
	t, ok := templates[path]
	if !ok {
		// path에 해당하는 템플릿이 없으면 템플릿 객체 생성
		t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		templates[path] = t
	}
	// v값을 템플릿 내부로 전달하여 만들어진 최종 결과를 c.ResponseWriter에 출력
	t.Execute(c.ResponseWriter, v)
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}
