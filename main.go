package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"ex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

const VerifyMessage = "verified"

// 인증된 웹 요청만 허용
// 인증이 확인되면 다음 핸들러로 넘어가고, 인증이 확인되지 않으면 로그인 페이지로 이동.
// 만약 요청 경로가 ignore 에 등록된 URL이면 인증을 확인하지 않고 바로 다음 핸들로로 넘어간다.
func AuthHandler(next ex.HandlerFunc) ex.HandlerFunc {
	ignore := []string{"/login", "/index.html"}
	return func(c *ex.Context) {
		// URL prefix가 "/login","main/index.html" 이면 auth를 체크하지 않는다.
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}
		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// "X_AUTH" 쿠키 값이 없으면 "/login" 으로 이동
			c.Redirect("/login")
			return
		} else if err != nil {
			// 에러 처리
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			// 쿠키 값으로 인증이 확인되면 다음 핸들러로 넘어감.
			next(c)
			return
		}
		// "/login"으로 이동
		c.Redirect("/login")
	}

}

// 인증 토큰 확인
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}
func CheckLogin(username, password string) bool {
	// 로그인 처리
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)
	return username == USERNAME && password == PASSWORD
}

//
// 인증 토큰 생성
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

// 모든 웹 요청을 라우터가 받아 처리하도록.
func main() {
	// 서버 생성
	s := ex.NewServer()

	s.HandleFunc("GET", "/", func(c *ex.Context) {
		c.RenderTemplate("/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *ex.Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", "/users/:id", func(c *ex.Context) {
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *ex.Context) {
		u := User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	s.HandleFunc("POST", "/users", func(c *ex.Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *ex.Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})

	s.HandleFunc("GET", "/login", func(c *ex.Context) {
		// "login.html" 렌더링
		c.RenderTemplate("/login.html", map[string]interface{}{"message": "로그인 해주세요."})
	})

	s.HandleFunc("POST", "/login", func(c *ex.Context) {
		// 로그인 정보를 확인하여 쿠키에 인증 토큰 값 기록
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		// id와 password가 맞지 않으면 다시 "/login" 페이지 렌더링
		c.RenderTemplate("/login.html",
			map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다"})
	})
	s.Use(AuthHandler)

	s.Run(":8080")

}
