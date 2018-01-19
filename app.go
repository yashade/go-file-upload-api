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
	viper.SetDefault("key", "mykey")

	router := httprouter.New()

	router.POST("/upload", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		key := r.Header.Get("Key")

		if key == viper.GetString("key") {
			file, handler, err := r.FormFile("file")
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			name := handler.Filename

			f, err := os.OpenFile(
				name,
				os.O_WRONLY|os.O_CREATE,
				os.FileMode(viper.GetInt("perms")))
			defer f.Close()

			log.Printf(handler.Filename)
			io.Copy(f, file)

			w.Write([]byte("\n" + name + "\n\n"))
		} else {
			w.WriteHeader(401)
			w.Write([]byte("\n" + "Unauthorized" + "\n\n"))
		}
	})

	if viper.GetBool("fcgi") {
		fcgi.Serve(nil, router)
	} else {
		http.ListenAndServe(":"+viper.GetString("port"), router)
	}
}
