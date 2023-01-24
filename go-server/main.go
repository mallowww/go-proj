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

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	// http.HandleFunc("/form", formHandler)
	// http.HandleFunc("/register", registerHandler)

	fmt.Printf("enter port number(as int)\n")
	reader := bufio.NewReader(os.Stdin)
	addrStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("can't read string", err)
	}
	addrStr = strings.TrimSpace(addrStr)
	addr, err := strconv.Atoi(addrStr)
	if err != nil {
		log.Fatal("can't convert to int", err)
	}

	fmt.Printf("starting server from port :%v\n", addr)
	http.ListenAndServe(":"+strconv.Itoa(addr), nil)
}
