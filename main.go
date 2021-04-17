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
\____/  \____//_/ /_/\__,_/ /____/ /_/ /_/\___//_/`

type Options struct {
	Hash        string
	List        string
	Concurrency int
	Version     bool
}

var md5 []func(param string, param2 string) string

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

func alpha(hashvalue string) {
	return
}

func Beta(hashvalue string, hashtype string) string {

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
	r := string(s)

	return r

}

func Theta(hashvalue string, hashtype string) string {

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
	r := string(s)

	return r

}

func hashCrack(hashvalue string) {

	if len(hashvalue) == 32 {
		println("[!] Hash Function : MD5")
		for _, api := range md5 {
			r := api(hashvalue, "md5")
			println(r)
		}

	}

}

func hashOnly(hashvalue string) {

}

func main() {

	md5 = append(md5, Beta, Theta)

	fmt.Println(banner)
	options := ParseOptions()

	if options.Hash == "" && options.List == "" {
		fmt.Println("hash string or hash file must be provided")
		flag.Usage()
		return
	}
	hashCrack(options.Hash)

}
