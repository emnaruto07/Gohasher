package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const Version = `v1.0`

const banner = `_________     ______               ______              
__  ____/________  /_______ __________  /______________
_  / __ _  __ \_  __ \  __  /_  ___/_  __ \  _ \_  ___/
/ /_/ / / /_/ /  / / / /_/ /_(__  )_  / / /  __/  /    
\____/  \____//_/ /_/\__,_/ /____/ /_/ /_/\___//_/      
`

type Options struct {
	Hash        string
	List        string
	Concurrency int
	Version     bool
}

var md5, sha1, sha256, sha384, sha512 []func(param string, param2 string) string

var result []string

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

	res := strings.Split(r, "=")
	ress := res[1]

	return ress[:len(ress)-1]

}

func Theta(hashvalue string, hashtype string) string {

	path := fmt.Sprintf("https://md5decrypt.net/Api/api.php?hash=%s&hash_type=%s&email=deanna_abshire@proxymail.eu&code=1152464b80a61728", hashvalue, hashtype)

	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return string(body)

}

func hashCrack(hashvalue string) []string {

	if len(hashvalue) == 32 {
		println("[!] Hash Function : MD5\n")
		for _, api := range md5 {
			r := api(hashvalue, "md5")
			if r != "" {
				result = append(result, r)

			}
		}

	} else if len(hashvalue) == 40 {
		println("[!] Hash Function : SHA1\n")
		for _, api := range sha1 {
			r := api(hashvalue, "sha1")
			if r != "" {
				result = append(result, r)

			}
		}

	} else if len(hashvalue) == 64 {
		println("[!] Hash Function : SHA-256\n")
		for _, api := range sha1 {
			r := api(hashvalue, "sha256")
			if r != "" {
				result = append(result, r)

			}
		}

	} else if len(hashvalue) == 96 {
		println("[!] Hash Function : SHA-384\n")
		for _, api := range sha1 {
			r := api(hashvalue, "sha384")
			if r != "" {
				result = append(result, r)

			}
		}

	} else if len(hashvalue) == 128 {
		println("[!] Hash Function : SHA-512\n")
		for _, api := range sha1 {
			r := api(hashvalue, "sha512")
			if r != "" {
				result = append(result, r)

			}
		}

	} else {
		println("[!!] This hash ype is not supported")
		os.Exit(0)
	}

	return result
}

func hashOnly(hashvalue string) {
	res := hashCrack(hashvalue)
	println("Cracked hash of " + hashvalue + " value is: " + res[1])

}

func main() {

	md5 = append(md5, Beta, Theta)
	sha1 = append(sha1, Beta, Theta)
	sha256 = append(sha256, Beta, Theta)
	sha384 = append(sha384, Beta, Theta)
	sha512 = append(sha512, Beta, Theta)

	fmt.Println(banner)
	options := ParseOptions()

	if options.Hash == "" && options.List == "" {
		fmt.Println("hash string or hash file must be provided")
		flag.Usage()
		return
	}
	hashOnly(options.Hash)

}
