package main

import (
  "os"
  "fmt"
  "net/http"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/sessions"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/markbates/goth/providers/facebook"
)

func main() {

  goth.UseProviders(
    facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
  )

  gothic.GetProviderName = func(req *http.Request) (string, error) {
    return "facebook", nil
  }

  // instantiate Martini
  m := martini.Classic()

  m.Use(func(res http.ResponseWriter, req *http.Request) {
    if 1 == 2 {
      res.WriteHeader(http.StatusUnauthorized)
    }
  })

  m.Get("/", func(session sessions.Session) string {
    //user, err := gothic.FetchUser()
    //return fmt.Sprintf("Hello world!, access_token=%s", token.(string))
    return "Hello world!"
  })

  m.Get("/auth/:provider/callback", func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
    user, err := gothic.CompleteUserAuth(res, req)
    if err != nil {
      return "Something went wrong."
    }
    return fmt.Sprintf("Logged in as %s, access_token=%s", user.Name, user.AccessToken)
  })

  m.Get("/auth/:provider", gothic.BeginAuthHandler)

  m.Run()
}
