package trace

import (
	"fmt"
	"io"
)

//Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェース

type Tracer interface {
	Trace(...interface{}) //任意の型の引数を受け取れる
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w} //実装を隠しているテクニック
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

//OffはTraceメソッドの呼び出しを無視するTracerを返す

func Off() Tracer {
	return &nilTracer{}
}
