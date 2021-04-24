package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var Teal = Color("\033[1;36m%s\033[0m")

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

const Version = `v1.0`

const banner = `_________     ______               ______              
__  ____/________  /_______ __________  /______________
_  / __ _  __ \_  __ \  __  /_  ___/_  __ \  _ \_  ___/
/ /_/ / / /_/ /  / / / /_/ /_(__  )_  / / /  __/  /    
\____/  \____//_/ /_/\__,_/ /____/ /_/ /_/\___//_/      

             created by @emnaruto07 & @s0u1z 
`

type Options struct {
	Hash string
	List string
	//Concurrency int
	Version bool
}

var crackersByLength = map[int]HashCracker{
	32:  NewGeneralCracker("md5"),
	40:  NewGeneralCracker("sha1"),
	64:  NewGeneralCracker("sha256"),
	96:  NewGeneralCracker("sha384"),
	128: NewGeneralCracker("sha512"),
}

func ReadstdFile() []string {
	var rows []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
		i := 0
		for i < len(rows) {
			rows[i] = strings.TrimSpace(rows[i])
			if rows[i] == "" {
				rows = append(rows[:i], rows[i+1:]...)
			}
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return rows
}

func main() {

	start := time.Now()

	fmt.Println(Teal(banner))
	options := ParseOptions()
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	/*
		    File.Mode() return a FileMode flag.

		You can see what is each letter here FileMode.
			The flag we are looking for is os.ModeNamedPipe. When this flag is on it means that we have a pipe.
			This way we can know when our command is receiving stdout from another process.

	*/
	if fi.Mode()&os.ModeNamedPipe == 0 {
		if options.Hash == "" && options.List == "" {
			fmt.Print("[!!!]Hash string Or Hash file must be provided.\n")
			flag.Usage()
			os.Exit(1)
			return
		}

		if options.Hash != "" {

			cracker, found := crackersByLength[len(options.Hash)]
			if !found {
				fmt.Println("[!!]Unsupported hash")
				return
			}
			res, err := cracker.Crack(options.Hash)
			if err != nil {
				fmt.Println(err)
			}
			if res != "" {
				fmt.Printf("Hash(%v): %v value: %v", cracker.String(), options.Hash, res)
			} else {
				fmt.Printf("Hash(%v): %v value: Not found\n", cracker.String(), options.Hash)
			}
		} else if options.List != "" {

			file, err := ParseFile(options.List)

			if err != nil {
				fmt.Println(err)
			}
			for _, f := range file {
				cracker, found := crackersByLength[len(f)]
				if !found {
					fmt.Println("[!!]Usupported hash")
					return
				}

				res, err := cracker.Crack(f)
				if err != nil {
					fmt.Println(err)
				}
				if res != "" {
					//fmt.Println("Hash(%v):" + %v + " value:" + %v) // Hash(Md5):hash value:value
					fmt.Printf("Hash(%v): %v value: %v", cracker.String(), f, res)
				} else {
					fmt.Printf("Hash(%v): %v value: Not found\n", cracker.String(), f)
				}

			}

		}

	} else {
		for _, r := range ReadstdFile() {

			var res string
			var err error
			cracker, found := crackersByLength[len(r)]
			if !found {
				fmt.Println("[!!]Unsupported hash")
				return
			}

			res, err = cracker.Crack(r)

			if err != nil {
				fmt.Println(err)
			}
			if res != "" {
				fmt.Printf("Hash(%v): %v value: %v\n", cracker.String(), r, res)
			} else {
				fmt.Printf("Hash(%v): %v value: Not found\n", cracker.String(), r)
			}

		}

	}

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

}

func ParseFile(filename string) ([]string, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print("Error:", err)
	}
	rows := strings.Split(string(d), "\n")
	i := 0
	for i < len(rows) {
		rows[i] = strings.TrimSpace(rows[i])
		if rows[i] == "" {
			rows = append(rows[:i], rows[i+1:]...)
		}
		i++
	}
	return rows, nil

}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Hash, "hash", "", "Provide the hash string with this flag.")
	flag.StringVar(&options.List, "l", "", "Provide the file, Which contain the hashes.")
	//flag.IntVar(&options.Concurrency, "c", 10, "Number of concurrent goroutines for resolving")
	flag.BoolVar(&options.Version, "version", false, "Show the version of GoHasher.")
	flag.Parse()

	if options.Version {
		fmt.Println("[+]Current Version:", Version)
		os.Exit(0)
	}

	return options
}

type GeneralCracker struct {
	hashType string
}

func (c *GeneralCracker) hashToolkit(hashvalue string) string {

	resp, err := http.Get("https://hashtoolkit.com/decrypt-hash/?hash=" + hashvalue)
	if err != nil {
		fmt.Print("Error:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Error:", err)
	}
	re := regexp.MustCompile(`text=.*?"`)

	s := re.Find(body)
	r := string(s)

	if r != "" {

		res := strings.Split(r, "=")
		ress := res[1]
		decodedValue, err := url.QueryUnescape(ress[:len(ress)-1])
		if err != nil {
			fmt.Println(err)
		}

		return decodedValue

	} else {
		return r
	}

}

func (c *GeneralCracker) md5Decrypt(hashvalue string) string {
	// make md5crypt.net call, use c.hashType to get hash type
	path := fmt.Sprintf("https://md5decrypt.net/Api/api.php?hash=%s&hash_type=%s&email=deanna_abshire@proxymail.eu&code=1152464b80a61728", hashvalue, c.hashType)

	resp, err := http.Get(path)
	if err != nil {
		fmt.Print("Error:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print("Error:", err)
	}

	return string(body)

}

func (c *GeneralCracker) nitrxGen(hashvalue string) string {
	// make md5crypt.net call, use c.hashType to get hash type

	resp, err := http.Get("http://www.nitrxgen.net/md5db/" + hashvalue)
	if err != nil {
		fmt.Print("Error:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print("Error:", err)
	}

	return string(body)

}

type HashCracker interface {
	String() string
	Crack(hash string) (string, error)
}

func (c *GeneralCracker) String() string {
	return c.hashType
}

func (c *GeneralCracker) Crack(hash string) (string, error) {
	var result = make(chan string)

	go func() {

		if c.hashType == "md5" {
			result <- c.nitrxGen(hash)
		} else if c.hashType == "sha1" {
			result <- c.hashToolkit(hash)
		} else {
			result <- c.md5Decrypt(hash)
		}

	}()

	return <-result, nil

}

func NewGeneralCracker(hashType string) *GeneralCracker {
	return &GeneralCracker{hashType}
}
