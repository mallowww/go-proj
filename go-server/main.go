package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm err - %v\n", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name: %s\n", name)
	fmt.Fprintf(w, "Address: %s\n", address)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found", 400)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method is not supported", 400)
		return
	}

	fmt.Fprintf(w, "register")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/register", registerHandler)

	fmt.Printf("enter port number(as int)\n")
	reader := bufio.NewReader(os.Stdin)
	addrStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("can't read string", err)
	}
	// trims the whitespace from the input string before converting it to an integer with strconv.Atoi()
	addrStr = strings.TrimSpace(addrStr)
	addr, err := strconv.Atoi(addrStr)
	if err != nil {
		log.Fatal("can't convert to int", err)
	}

	fmt.Printf("starting server from port :%v\n", addr)
	if err := http.ListenAndServe(":"+strconv.Itoa(addr), nil); err != nil {
		log.Fatal(err)
	}
}
