package main

import (
  "github.com/valyala/fasthttp"
  "github.com/go-redis/redis/v7"
  "os"
  "strconv"
  "log"
)

func main() {
  port := os.Getenv("PORT")
  if _, err := strconv.ParseInt(port, 10, 16); err != nil {
    port = "8080"
  }

  client := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })

  handler := func(ctx *fasthttp.RequestCtx) {
    path := string(ctx.Path())
    switch path {
    case "/":
      ctx.Redirect("https://teamwaterloop.ca", 302)
    case "/admin": // check if logged in?
      ctx.SendFile("admin.html")
      break
    default:
      val, err := client.Get("url_" + path[1:]).Result()
      if err == nil {
        ctx.Redirect(val, 302)
      } else {
        ctx.Error("404: redirect not found", fasthttp.StatusNotFound)
      }
      break
    }
  }

  log.Println("Listening on 127.0.0.1:" + port)
  log.Fatal(fasthttp.ListenAndServe("127.0.0.1:"+port, handler))
}
