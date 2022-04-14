package main

import (
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	router := router()
	middleware := negroni.Classic()
	middleware.UseHandler(router)

	if os.Getenv("ENV") == "local" {
		router.Headers("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Methods")
		router.Headers("Access-Control-Allow-Origin", "*")
		router.Headers("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	}

	http.ListenAndServe(":"+os.Getenv("PORT_NUM"), corseHandler(middleware))
}

func router() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("GET").
		Path("/health-check").
		HandlerFunc(controllers.HealthcheckControllerHandle)
	router.
		Methods("POST").
		Path("/upload-image").
		HandlerFunc(controllers.UploadImageControllerHandle)
	router.
		Methods("POST").
		Path("/delete-image").
		HandlerFunc(controllers.DeleteImageControllerHandle)
	router.
		PathPrefix("/cdn/").
		Handler(http.StripPrefix("/cdn/", http.FileServer(http.Dir("/media/"))))
	router.
		Methods("GET").
		PathPrefix("/").
		HandlerFunc(controllers.MediaControllerHandle)

	return router
}

func corseHandler(handler http.Handler) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, webRequest *http.Request) {
		setHeaders(responseWriter)
		if webRequest.Method != "OPTIONS" {
			handler.ServeHTTP(responseWriter, webRequest)
		}
	}
}

func setHeaders(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Methods")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
}
