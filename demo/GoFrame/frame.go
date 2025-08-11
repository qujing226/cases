package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	fmt.Println("Hello GoFrame:", gf.VERSION)
	s := g.Server()

	s.BindHandler("/", func(r *ghttp.Request) {
		var req HelloReq
		if err := r.Parse(&req); err != nil {
			r.Response.Write(err.Error())
			return
		}
		if req.Name == "" {
			r.Response.Write("name should not be empty")
			return
		}
		if req.Age <= 0 {
			r.Response.Write("invalid age value")
			return
		}
		r.Response.Writef(
			"Hello %s !,your Age is %d",
			req.Name, req.Age)

		r.Response.Writef("Hello %s ,your age is %d", r.Get("name", "unknown").String(),
			r.Get("age").Int())
	})
	s.SetPort(9000)
	s.Run()
}

type HelloReq struct {
	g.Meta `path:"/" method:"get"`
	Name   string `v:"required" dc:"姓名"`
	Age    int    `v:"required" dc:"年龄"`
}
type Hello struct {
}

type HelloRes struct {
}

func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Writef(
		"Hello %s ! you Age is %d",
		req.Name, req.Age,
	)
	return
}
