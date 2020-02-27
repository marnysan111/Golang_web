package main

import (
	"html/template"
	"log"
	"net/http"
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
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil { //ListenAndServeメソッドでポート8080を使いWebサーバを開始する
		log.Fatal("ListenAndServe:", err)
	}
}
