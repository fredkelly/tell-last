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
  globalApp := fb.New(os.Getenv("FACEBOOK_APP_ID"), os.Getenv("FACEBOOK_SECRET"))
  globalApp.RedirectUri = "http://localhost:3000/auth/facebook/callback" // TODO

  // https://developers.facebook.com/docs/graph-api/securing-requests
  globalApp.EnableAppsecretProof = true

  // instantiate Martini
  m := martini.Classic()
  m.Use(martini.Logger())

  m.Use(func(res http.ResponseWriter, req *http.Request) {
    accessToken := req.Header.Get("Authorization") // TODO expect "Bearer: {TOKEN}" format?
    session := globalApp.Session(accessToken)

    // TODO make currentUser globally available
    // TODO add User struct and use res.Decode(&user)
    currentUser, err := session.Get("/me", nil)

    if err != nil {
      // err can be an facebook API error.
      // if so, the Error struct contains error details.
      if e, ok := err.(*fb.Error); ok {
        log.Printf("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]", e.Message, e.Type, e.Code, e.ErrorSubcode)
      }

      res.WriteHeader(http.StatusUnauthorized)
    }

    log.Printf("Logged in as: %s", currentUser["email"])
  })

  m.Get("/", func() string {
    // serve SPA ?
    return "Hello world"
  })

  m.Get("/tells", func() string {
    // get all tells for currentUser and render JSON array
    return ""
  })

  m.Post("/tells", func(params martini.Params) string {
    // create new tell for currentUser
    return ""
  })

  m.Run()
}
