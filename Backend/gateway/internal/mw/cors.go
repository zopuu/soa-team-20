package mw

import (
	"github.com/go-chi/cors"
)

func CORS(origins []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET","POST","PUT","PATCH","DELETE","OPTIONS"},
		AllowedHeaders:   []string{"Accept","Authorization","Content-Type","X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
