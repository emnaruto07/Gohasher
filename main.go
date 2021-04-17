package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const Version = `v1.0`

const banner = `_________     ______               ______              
__  ____/________  /_______ __________  /______________
_  / __ _  __ \_  __ \  __  /_  ___/_  __ \  _ \_  ___/
/ /_/ / / /_/ /  / / / /_/ /_(__  )_  / / /  __/  /    
\____/  \____//_/ /_/\__,_/ /____/ /_/ /_/\___//_/v1.0 `

type Options struct {
	Hash        string
	List        string
	Concurrency int
	Version     bool
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Hash, "hash", "", "Provide the hash string with this flag.")
	flag.StringVar(&options.List, "l", "", "Provide the file, Which contain the hashes.")
	flag.IntVar(&options.Concurrency, "c", 10, "Number of concurrent goroutines for resolving")
	flag.BoolVar(&options.Version, "version", false, "Show the version of GoHasher.")
	flag.Parse()

	if options.Version {
		fmt.Println("[+]Current Version:", Version)
		os.Exit(0)
	}

	return options
}

func alpha(hashvalue string, hashtype string) {
	return
}

func beta(hashvalue string, hashtype string) {

	resp, err := http.Get("https://hashtoolkit.com/decrypt-hash/?hash=" + hashvalue)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`text=.*?"`)

	//decodedValue, err := url.QueryUnescape(s)
	s := re.Find(body)

	fmt.Println(string(s))

}

func main() {

	fmt.Println(banner)
	options := ParseOptions()

	if options.Hash == "" && options.List == "" {
		fmt.Println("hash string or hash file must be provided")
		flag.Usage()
		return
	}

	beta(options.Hash, "")

}
