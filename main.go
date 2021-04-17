package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	s := re.FindString(string(body))

	decodedValue, err := url.QueryUnescape(s)

	fmt.Println(decodedValue)

}

func main() {

	fmt.Println(banner)
	options := ParseOptions()

	//hash := flag.String("hash", "", "Contains the hash values")
	//hashfile := flag.String("l", "", "This contains file containing hashes")
	//threads := flag.Int("c", 20, " Contains the thread value")
	//version := flag.Bool("v", false, "Show current program version")
	//flag.Parse()

	//if *version {
	//	fmt.Printf("The current version of program is :%v", Version)
	//	os.Exit(0)
	//}

	if options.Hash == "" && options.File == "" {
		fmt.Println("hash string or hash file must be provided")
		flag.Usage()
		return
	}

	beta(options.Hash, "")

}
