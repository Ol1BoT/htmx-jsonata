package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/blues/jsonata-go"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		bt, err := ioutil.ReadFile("index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
		}

		w.Write(bt)
		return
	})

	r.Post("/validate-mapping", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			io.WriteString(w, basiDiv(err.Error()))
			return
		}

		mapping := r.Form.Get("mapping_writer")
		input := r.Form.Get("json_input")

		v, err := ValidateJSONata(input, mapping)
		if err != nil {
			io.WriteString(w, basiDiv(err.Error()))
			return
		}

		io.WriteString(w, basiDiv(v))

		return

	})

	http.ListenAndServe(":8080", r)

}

func ValidateJSONata(input, mapping string) (string, error) {
	var i any

	if err := json.Unmarshal([]byte(input), &i); err != nil {
		return "", err
	}

	expr, err := jsonata.Compile(mapping)
	if err != nil {
		return "", err
	}

	out, err := expr.Eval(i)
	if err != nil {
		return "", err
	}

	bt, err := json.Marshal(out)
	if err != nil {
		return "", err
	}

	return string(bt), nil

}

func basiDiv(msg string) string {
	return fmt.Sprintf("<div>%s</div>", msg)
}
