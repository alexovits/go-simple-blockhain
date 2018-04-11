package main

import (
	"./blockchain"
	"./util"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Blockchain []blockchain.Block

type Message struct {
	Data int
}

var mutex = &sync.Mutex{}

func run() {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Creates the necessary routes for the GET and POST requests
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		util.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// Ensure the block creation is atomic
	mutex.Lock()
	newBlock := blockchain.GenerateBlock(Blockchain[len(Blockchain)-1], m.Data)
	mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	util.RespondWithJSON(w, r, http.StatusCreated, newBlock)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := blockchain.Block{}
		genesisBlock = blockchain.Block{0, t.String(), 0, blockchain.CalculateHash(genesisBlock), "", 1, ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	run()
}
