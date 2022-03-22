package main

import (
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func main() {
	rd = render.New()
	m := MakeWebHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
