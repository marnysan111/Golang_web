package main

import (
	"flag"
	"go/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//templは一つのテンプレートを表す
type templateHandler struct {
	once     sync.Once //一回だけコンパイル？
	filename string
	templ    *template.Template
}

//ServeHTTPはHTTPリクエストを処理する
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename))) //htmlのあるフォルダ指定
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/github"),
		google.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Webサーバを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil { //ListenAndServeメソッドでポート8080を使いWebサーバを開始する
		log.Fatal("ListenAndServe:", err)
	}
}
