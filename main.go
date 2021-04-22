package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

//go:embed index.html
var index []byte

type response struct {
	Header     []attribute
	Query      []attribute
	FormParams []attribute
}

type attribute struct {
	Name   string
	Values []string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	if p, err := strconv.Atoi(os.Getenv("PORT")); err != nil || p == 0 {
		os.Setenv("PORT", "8000")
	}

	fmt.Printf("Listening on: http://0.0.0.0:%s\n", os.Getenv("PORT"))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), mux); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp := response{getSortedAttrs(r.Header), getSortedAttrs(r.URL.Query()), getSortedAttrs(r.PostForm)}

	if strings.Contains(r.Header.Get("Accept"), "application/json") || r.URL.Query().Get("format") == "json" {
		jsonOutput(w, resp)

		return
	}

	htmlOutput(w, resp)
}

func getSortedAttrs(values map[string][]string) []attribute {
	names := []string{}
	for name := range values {
		names = append(names, name)
	}

	sort.Strings(names)

	attrs := []attribute{}
	for _, h := range names {
		attrs = append(attrs, attribute{h, values[h]})
	}

	return attrs
}

func jsonOutput(w http.ResponseWriter, resp response) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func htmlOutput(w http.ResponseWriter, resp response) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(index)

	if len(resp.Header) > 0 {
		printValues(w, "Header", resp.Header)
	}

	if len(resp.Query) > 0 {
		printValues(w, "Query", resp.Query)
	}

	if len(resp.FormParams) > 0 {
		printValues(w, "Form Params", resp.FormParams)
	}
}

func printValues(w io.Writer, title string, attrs []attribute) {
	fmt.Fprintf(w, "<h2>%s</h2>", title)

	for _, val := range attrs {
		fmt.Fprint(w, "<details><summary>")
		fmt.Fprint(w, val.Name+"</summary>")

		for _, v := range val.Values {
			fmt.Fprint(w, "<p>"+v+"</p>")
		}

		fmt.Fprintln(w, "</details>")
	}
}
