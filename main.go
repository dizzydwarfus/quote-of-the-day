package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/fatih/color"
)

var (
	Green          = color.New(color.FgGreen, color.Bold).PrintfFunc()
	Faint          = color.New(color.Faint).PrintfFunc()
	Red            = color.New(color.FgRed).SprintfFunc()
	RequestCounter = 0
)

type Server struct {
	listenAddr string
	tmpl       *template.Template
}

type Quote struct {
	ID       int    `json:"id"`
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

type QuotesData struct {
	Quotes []Quote `json:"quotes"`
}

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "the server address")
	flag.Parse()

	tmpl, err := template.ParseFiles("templates/quote.html")
	if err != nil {
		log.Fatalf("Template error: %v", err)
	}

	s := NewServer(*listenAddr, tmpl)

	log.Printf("Server running on port: http://localhost%v\n", *listenAddr)
	log.Fatal(s.Start())
}

func NewServer(listenAddr string, tmpl *template.Template) *Server {
	return &Server{
		listenAddr: listenAddr,
		tmpl:       tmpl,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	fsHandler := http.FileServer(http.Dir("./web"))

	mux.Handle("/", fsHandler)
	mux.HandleFunc("/quote", s.getQuoteOfTheDay)

	return http.ListenAndServe(s.listenAddr, mux)
}

func (s *Server) getQuoteOfTheDay(w http.ResponseWriter, r *http.Request) {
	RequestCounter++

	quotesData, err := readQuotesJson()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	quote := getRandomQuote(quotesData)

	// Execute template with Quote struct
	w.Header().Set("Content-Type", "text/html")
	if err := s.tmpl.Execute(w, quote); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func readQuotesJson() (*QuotesData, error) {
	jsonFile, err := os.Open("quotes.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var quotesData QuotesData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&quotesData); err != nil {
		return nil, err
	}
	return &quotesData, nil
}

func getRandomQuote(quotes *QuotesData) Quote {
	randomInt := rand.IntN(len(quotes.Quotes))
	quote := quotes.Quotes[randomInt]
	return quote
}
