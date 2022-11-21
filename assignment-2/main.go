package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/antonlindstrom/pgstore"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var schema = `CREATE TABLE users (
	id serial primary key,
	username varchar(50) DEFAULT NULL,
	first_name varchar(200) NOT NULL,
	last_name varchar(200) NOT NULL,
	password varchar(120) DEFAULT NULL);`

type Users struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Password   string `json:"password"`
}

type Renderer struct {
	template *template.Template
	debug    bool
	location string
}

func NewRenderer(location string, debug bool) *Renderer {
	tpl := new(Renderer)
	tpl.location = location
	tpl.debug = debug

	tpl.ReloadTemplates()

	return tpl
}

func (t *Renderer) ReloadTemplates() {
	t.template = template.Must(template.ParseGlob(t.location))
}

func (t *Renderer) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	if t.debug {
		t.ReloadTemplates()
	}

	return t.template.ExecuteTemplate(w, name, data)
}

func newPostgresStore() *pgstore.PGStore {
	url := "postgres://postgresuser:postgrespassword@localhost:5435/postgres?sslmode=disable"
	authKey := []byte("my-auth-key-very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store, err := pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	return store
}

var store = newPostgresStore()

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgresuser:postgrespassword@localhost:5435/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// db.MustExec(schema)

	r := echo.New()
	r.Renderer = NewRenderer("./*.html", true)

	r.GET("/", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "register.html", nil)
	})

	r.POST("/register", func(ctx echo.Context) error {
		_, err := db.NamedQuery(`INSERT INTO public.users
		(username, first_name, last_name, "password")
		VALUES(:username, :firstName, :lastName, :password)`, map[string]interface{}{
			"username":  ctx.FormValue("username"),
			"firstName": ctx.FormValue("firstName"),
			"lastName":  ctx.FormValue("lastName"),
			"password":  ctx.FormValue("psw"),
		})

		if err != nil {
			fmt.Println(err)
			return err
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	r.GET("/login-html", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", func(ctx echo.Context) error {
		var user Users
		err := db.Get(&user, `SELECT * FROM users WHERE username = $1 AND password =$2`, ctx.FormValue("username"), ctx.FormValue("password"))
		if err != nil {
			fmt.Println("ERRRRROR", err)
			return err
		}

		session, err := store.Get(ctx.Request(), fmt.Sprintf("SESSION-ID-%v", user.ID))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		session.Values["username"] = user.Username
		if err = session.Save(ctx.Request(), ctx.Response()); err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, "/home")
	})

	r.GET("/home", func(ctx echo.Context) error {
		ctx.Render(http.StatusOK, "home.html", nil)

		return nil
	})

	r.POST("/home", func(ctx echo.Context) error {
		ctx.Render(http.StatusOK, "home.html", nil)

		return nil
	})

	r.GET("/logout", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	r.POST("/logout", func(ctx echo.Context) error {
		var user Users
		err := db.Get(&user, `SELECT * FROM users WHERE username = $1 AND password =$2`, ctx.FormValue("username"), ctx.FormValue("password"))
		if err != nil {
			fmt.Println("ERRRRROR", err)
			return err
		}

		session, err := store.Get(ctx.Request(), fmt.Sprintf("SESSION-ID-%v", user.ID))
		session.Options.MaxAge = -1
		session.Save(ctx.Request(), ctx.Response())
		return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	r.Start(":9000")
}
