// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 195.

// Http4 is an e-commerce server that registers the /list and /price
// endpoint by calling http.HandleFunc.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//!+main

func main() {
	db := database{
		data:    make(map[string]dollars),
		RWMutex: sync.RWMutex{},
	}
	db.data["shoes"] = 50
	db.data["socks"] = 5

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!-main

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// type database map[string]dollars

type database struct {
	data map[string]dollars
	sync.RWMutex
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	db.RLock()
	defer db.RUnlock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.RLock()
	defer db.RUnlock()
	if price, ok := db.data[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, err := strconv.Atoi(req.URL.Query().Get("price"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "new price is invalid: %d", price)
		return
	}
	db.Lock()
	db.data[item] = dollars(price)
	db.Unlock()

	db.RLock()
	defer db.RUnlock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db.data[item]; !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "item need to be deleted is not existed: %s\n", item)
		return
	}
	db.Lock()
	delete(db.data, item)
	db.Unlock()

	db.RLock()
	defer db.RUnlock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, err := strconv.Atoi(req.URL.Query().Get("price"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "new price is invalid: %d", price)
		return
	}

	db.Lock()
	defer db.Unlock()
	if _, ok := db.data[item]; ok {
		db.data[item] = dollars(price)
		for item, price := range db.data {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
