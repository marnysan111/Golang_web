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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Webサーバを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil { //ListenAndServeメソッドでポート8080を使いWebサーバを開始する
		log.Fatal("ListenAndServe:", err)
	}
}
