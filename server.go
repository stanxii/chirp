package main

import (
	"flag"
	"fmt"
	"net/http"

	"chirp.com/app"
	"chirp.com/config"
	"chirp.com/controllers"
	"chirp.com/email"
	"chirp.com/middleware"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()
	cfg := config.LoadConfig(*boolPtr)
	services := app.Setup(cfg)
	defer services.Close()
	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Lenslocked.com Support", "support@mg.lenslocked.com"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	)

	router := app.NewRouter()

	tweetsAPI := controllers.NewTweets(services.Tweet, services.Like, services.Tag, services.Tagging)
	tagsAPI := controllers.NewTags(services.Tag, services.Tagging)
	usersAPI := controllers.NewUsers(services.User, services.Like, services.Follow, services.Tweet, emailer)

	//init middleware
	userMw := middleware.NewUserMw(services.User)
	requireUserMw := middleware.NewRequireUserMw(userMw)

	//test route
	router.HandleFunc("/ping", ping).Methods("GET")
	//api routes
	subRouter := router.PathPrefix("/api").Subrouter()
	controllers.ServeUserResource(subRouter, usersAPI, &requireUserMw)
	controllers.ServeTweetResource(subRouter, tweetsAPI, &requireUserMw)
	controllers.ServeTagResource(subRouter, tagsAPI, &requireUserMw)

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port),
		userMw.Apply(router))
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pinging the server...Success!\n")
}
