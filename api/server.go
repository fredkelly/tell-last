package main

import (
  "os"
  "log"
  //"fmt"
  "net/http"

  "github.com/go-martini/martini"
  fb "github.com/huandu/facebook"
)

func main() {
  // create a global App var to hold your app id and secret.
  var globalApp = fb.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"))
  globalApp.RedirectUri = "http://localhost:3000/auth/facebook/callback"

  // instantiate Martini
  m := martini.Classic()
  m.Use(martini.Logger())

  m.Use(func(res http.ResponseWriter, req *http.Request) {
    _, err := fb.Get("/me", fb.Params{
      "access_token": req.Header.Get("Authorization"),
    })

    if err != nil {
      // err can be an facebook API error.
      // if so, the Error struct contains error details.
      if e, ok := err.(*fb.Error); ok {
        log.Printf("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]", e.Message, e.Type, e.Code, e.ErrorSubcode)
      }

      res.WriteHeader(http.StatusUnauthorized)
    }
  })

  m.Get("/", func() string {
    return "Hello world!"
  })

  m.Run()
}
