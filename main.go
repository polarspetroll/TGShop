package main

import (
	"app"
	"net/http"
	"os"
	"telegram"
	"time"
)

func main() {
  token := os.Getenv("TOKEN")
	webhookurl := "https://" + os.Getenv("DOMAIN") + "/" + token

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         ":" + os.Getenv("PORT"),
	}

	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/index", 302) })

	///////Webhook Handler/////////////
	http.HandleFunc("/"+token, app.BotHandler)

	///////// Main App Handlers////////////
	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/logout", app.LogOut)
	http.HandleFunc("/message", app.GroupMessage)
	http.HandleFunc("/del", app.Del)
	http.HandleFunc("/list", app.List)
	http.HandleFunc("/product", app.GetProduct)
	http.HandleFunc("/new", app.NewProduct)
	///////////////////////////////////////

	telegram.SetWebhook(webhookurl)
	server.ListenAndServe()
}
