package main

import (
  "os"
  "log"
  //"fmt"
  "net/http"
  "strconv"
  "time"

  "database/sql"
  "github.com/coopernurse/gorp"
  _ "github.com/go-sql-driver/mysql"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/encoder"
  "github.com/martini-contrib/binding"
  fb "github.com/huandu/facebook"
)

var dbmap *gorp.DbMap

type User struct {
  Id          int64   `db:"id" json:"id"`
  Uid         string  `db:"uid" json:"uid"`// Facebook identifier
  CreatedAt   int64   `db:"created_at" json:"created_at"`
  Email       string  `db:"email" json:"email"`
  FirstName   string  `db:"first_name" json:"first_name"`
  LastName    string  `db:"last_name" json:"last_name"`
}

type Tell struct {
  Id          int64   `db:"id" json:"id"`
  ToId        int64   `db:"to_id" json:"to_id" form:"to_id"`
  FromId      int64   `db:"from_id" json:"from_id" form:"from_id"`
  ReporterId  int64   `db:"reporter_id" json:"reporter_id,omitempty"`
  CreatedAt   int64   `db:"created_at" json:"created_at"`
  Body        string  `db:"body" json:"body" form:"body"`
}

// Initialise DB and setup tables
func initDb() *gorp.DbMap {
  db, err := sql.Open("mysql", "root@/tell_last")
  if err != nil {
    log.Printf("couldn't connect to database: %s", err)
  }

  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

  dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
  dbmap.AddTableWithName(Tell{}, "tells").SetKeys(true, "Id")

  err = dbmap.CreateTablesIfNotExists()
  if err != nil {
    log.Printf("couldn't create tables (%s)", err)
  }

  return dbmap
}

func findOrCreateUser(attrs fb.Result) *User {
  user := &User{}

  user.Uid       = attrs["id"].(string)
  user.FirstName = attrs["first_name"].(string)
  user.LastName  = attrs["last_name"].(string)

  // Store an email if we get one..
  // TODO complain if we don't?
  if email, ok := attrs["email"]; ok {
    user.Email = email.(string)
  }

  err := dbmap.SelectOne(&user, "SELECT * FROM users WHERE uid = ?", user.Uid)

  if err != nil {
    // create new user
    user.CreatedAt = time.Now().Unix()
    dbmap.Insert(user)
  } else {
    dbmap.Update(user)
  }

  return user
}

// all tells addressed to you
func (user User) getTells() []Tell {
  var tells []Tell
  _, err := dbmap.Select(&tells, "SELECT * FROM tells WHERE to_id = ?", user.Id)

  if err != nil {
    log.Printf("error loading tells for user id=%s (%s)", user.Id, err)
  }

  return tells
}

// create new tell
func (user User) Tell(tell Tell) Tell {
  // set reported by current user
  tell.ReporterId = user.Id

  err := dbmap.Insert(&tell)

  if err != nil {
    log.Printf("error inserting tell for user id=%s (%s)", user.Id, err)
  }

  return tell
}

func (tell Tell) Filter() interface{} {
  tell.ReporterId = -1
  return tell
}

func main() {
  // TODO move to own handler
  // create a global App var to hold your app id and secret.
  globalApp := fb.New(os.Getenv("FACEBOOK_APP_ID"), os.Getenv("FACEBOOK_SECRET"))
  globalApp.RedirectUri = "http://localhost:3000/auth/facebook/callback" // TODO

  // https://developers.facebook.com/docs/graph-api/securing-requests
  //globalApp.EnableAppsecretProof = true

  // instantiate Martini
  m := martini.Classic()
  m.Use(martini.Logger())

  // setup database
  dbmap = initDb()  
  defer dbmap.Db.Close()

  // Authentication
  m.Use(func(c martini.Context, res http.ResponseWriter, req *http.Request) {
    accessToken := req.Header.Get("Authorization") // TODO expect "Bearer: {TOKEN}" format?
    session := globalApp.Session(accessToken)

    // TODO make currentUser globally available
    // TODO add User struct and use res.Decode(&user)
    var user *User
    attrs, err := session.Get("/me", nil)

    if err != nil {
      // err can be an facebook API error.
      // if so, the Error struct contains error details.
      if e, ok := err.(*fb.Error); ok {
        log.Printf("Facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]", e.Message, e.Type, e.Code, e.ErrorSubcode)
      }

      res.WriteHeader(http.StatusUnauthorized)
    } else {
      user = findOrCreateUser(attrs)
      c.Map(user)
    }

    if user != nil {
      log.Printf("Logged in as: %s %s", user.FirstName, user.LastName)
    }
  })

  // Encoding
  m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
    // Use indentations. &pretty=1
    pretty, _ := strconv.ParseBool(r.FormValue("pretty"))
    // Use null instead of empty object for json &null=1
    null, _ := strconv.ParseBool(r.FormValue("null"))
    // JSON no matter what
    c.MapTo(encoder.JsonEncoder{PrettyPrint: pretty, PrintNull: null}, (*encoder.Encoder)(nil))
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
  })

  m.Get("/", func(user *User) string {
    // serve SPA ?
    return "Hello world" + user.FirstName
  })

  m.Get("/tells", func(enc encoder.Encoder, user *User) (int, []byte) {
    // get all tells for user and render JSON array
    tells := user.getTells()
    return http.StatusOK, encoder.Must(enc.Encode(tells))
  })

  m.Post("/tells", binding.Bind(Tell{}), func(enc encoder.Encoder, user *User, tell Tell) (int, []byte) {
    // create new tell for user
    tell = user.Tell(tell)
    return http.StatusOK, encoder.Must(enc.Encode(tell))
  })

  m.Run()
}
