- http.ResponseWriter 파라미터: HTTP Response에 무언가를 쓸 수 있게 한다.
- http.Request 파라미터: 입력된 Request요청을 검토할 수 있게 한다.

- ListenAndServe()
: 첫번째 파라미터로 포트 몇(숫자)에서 Request를 Listen할 것을 지정하고,
 두번째는 파라미터로 어떤 ServeMux를 사용할 지를 지정. nil인 경우 DefaultServeMux사용
 DefaultServeMux를 사용하는 경우, Handle() 또는 HandleFunc()을 사용하여 라우팅 패턴을 추가하게 된다.

- http.Handle()
: HTTP Handler를 정의하는 또 다른 방식으로 http.Handle()메서드를 사용할 수 있다. 
첫번째 파라미터로 URL(혹은 URL패턴)을 받아들이고,
두번째 파라미터로 http.Handler인터페이스를 갖는 객체를 받아들인다. 

- ServeHTTP()
: HTTP Response에 데이타를 쓰기 위한 Writer와 HTTP Request입력데이터를 파라미터로 갖는다.
