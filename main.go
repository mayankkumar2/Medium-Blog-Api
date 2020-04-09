package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"mediumFeedAPI/pkg"
	"os"
)
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsGet() && string(ctx.Path()) == "/api/v1/articles"{
		//log.Println(ctx.QueryArgs().Has("orgid"))
		if ctx.QueryArgs().Has("userid"){
			id := "@"+string(ctx.QueryArgs().Peek("userid"))
			m := pkg.GetArticles(id)
			ctx.Response.SetBody(m)
		} else if ctx.QueryArgs().Has("orgid"){

			id := string(ctx.QueryArgs().Peek("orgid"))
			m := pkg.GetArticles(id)
			ctx.Response.SetBody(m)
		} else {

			m := `{"status":"fail","articles":null}`
			ctx.Response.SetBody([]byte(m))
		}
		ctx.Response.Header.Set("content-type", "application/json")

	} else {
		ctx.Response.Header.Set("content-type","application/json")
		ctx.Response.SetBody([]byte(`{"status":"fail"}`))
	}
}


func main() {
	log.Println("Starting the server.")
	var port string
	if os.Getenv("PORT") == "" {
		port = ":80"
	} else {
		port = ":" + os.Getenv("PORT")
	}
	log.Println("Env PORT : " + os.Getenv("PORT"))
	log.Println("Run at port "+port[1:])
	fasthttp.ListenAndServe(port, fastHTTPHandler)
}
