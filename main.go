package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"mediumFeedAPI/pkg"
	"os"
	"github.com/lab259/cors"
)
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsGet() && string(ctx.Path()) == "/api/v1/articles"{
		if ctx.QueryArgs().Has("userid"){
			id := "@"+string(ctx.QueryArgs().Peek("userid"))
			m,statusCode := utils.GetArticles(id)
			ctx.Response.SetBody(m)
			ctx.SetStatusCode(statusCode)
		} else if ctx.QueryArgs().Has("orgid"){

			id := string(ctx.QueryArgs().Peek("orgid"))
			m,statusCode := utils.GetArticles(id)
			ctx.Response.SetBody(m)
			ctx.SetStatusCode(statusCode)
		} else {

			m := `{"status":"fail","articles":null}`
			ctx.Response.SetBody([]byte(m))
			ctx.SetStatusCode(400)
		}
		ctx.Response.Header.Set("content-type", "application/json")

	} else {
		ctx.Response.Header.Set("content-type","application/json")
		ctx.SetStatusCode(400)
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
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(fastHTTPHandler)
	fasthttp.ListenAndServe(port, handler)
}
