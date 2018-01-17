package main

import (
	"io"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("conf.yml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	viper.SetDefault("fcgi", false)
	viper.SetDefault("port", "9999")
	viper.SetDefault("perms", 0666)

	router := httprouter.New()

	router.POST("/upload", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		file, handler, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile(
			handler.Filename+"---"+time.Now().UTC().Format("2006-01-02_15:04:05"),
			os.O_WRONLY|os.O_CREATE,
			os.FileMode(viper.GetInt("perms")))
		defer f.Close()

		log.Printf(handler.Filename)
		io.Copy(f, file)
	})

	if viper.GetBool("fcgi") {
		fcgi.Serve(nil, router)
	} else {
		http.ListenAndServe(":"+viper.GetString("port"), router)
	}
}
