package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/fatih/color"
)

var (
	green          = color.New(color.FgGreen, color.Bold).PrintfFunc()
	faint          = color.New(color.Faint).PrintfFunc()
	red            = color.New(color.FgRed).SprintfFunc()
	requestCounter = 0
)

type Server struct {
	listenAddr string
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
	listenAddr := flag.String("listenaddr", ":8000", "the server address")
	flag.Parse()
	s := NewServer(*listenAddr)
	green("Server running on port: http://localhost%v\n", *listenAddr)
	log.Fatal(red("%v", s.Start()))
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", home)
	http.HandleFunc("/quote", s.getQuoteOfTheDay)
	return http.ListenAndServe(s.listenAddr, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Quote of The Day API, use /quote to get a random quote."))
}

func (s *Server) getQuoteOfTheDay(w http.ResponseWriter, r *http.Request) {
	requestCounter++

	quotesData, err := readQuotesJson()
	if err != nil {
		log.Fatal(red("%v", err))
	}

	quote := getRandomQuote(quotesData)

	json.NewEncoder(w).Encode(quote)
	faint("Request #%v: %v\n", requestCounter, *r)
}

func readQuotesJson() (*QuotesData, error) {
	jsonFile, err := os.Open("quotes.json")
	if err != nil {
		log.Fatal(red("%v", err))
		return &QuotesData{}, err
	}
	defer jsonFile.Close()

	var quotesData QuotesData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&quotesData); err != nil {
		log.Fatal(red("%v", err))
		return &QuotesData{}, err
	}
	return &quotesData, err
}

func getRandomQuote(quotes *QuotesData) Quote {
	randomInt := rand.IntN(len(quotes.Quotes))
	quote := quotes.Quotes[randomInt]
	return quote
}
