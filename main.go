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

	fs := http.FileServer(http.Dir("public"))

	r.Handle("/public/*", http.StripPrefix("/public/", fs))

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
			io.WriteString(w, basicDiv(err.Error()))
			return
		}

		mapping := r.Form.Get("mapping_writer")
		input := r.Form.Get("json_input")

		v, err := ValidateJSONata(input, mapping)
		if err != nil {
			io.WriteString(w, basicDiv(err.Error()))
			return
		}

		io.WriteString(w, basicTextArea(v))

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

func basicDiv(msg string) string {
	return fmt.Sprintf("<div>%s</div>", msg)
}

func basicTextArea(msg string) string {
	return fmt.Sprintf(`
      <textarea
        class="block p-2.5 w-full h-full text-sm text-gray-900 bg-gray-200 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
        hx-target="#mapping-results" hx-post="/validate-mapping"
        hx-trigger="keyup changed delay:500ms" hx-include="#json-input">
    	%s</textarea>
	`, msg)
}
