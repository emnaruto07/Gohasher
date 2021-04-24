package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
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
	Hash        string
	List        string
	Concurrency int
	Version     bool
}

var crackersByLength = map[int]HashCracker{
	32:  NewGeneralCracker("md5"),
	40:  NewGeneralCracker("sha1"),
	64:  NewGeneralCracker("sha256"),
	96:  NewGeneralCracker("sha384"),
	128: NewGeneralCracker("sha512"),
}

func main() {

	start := time.Now()
	fmt.Println(Teal(banner))
	options := ParseOptions()
	threads := options.Concurrency / 2
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

			hp, err := MakePipe(options.List, threads)
			if err != nil {
				fmt.Println(err)
				return
			}
			wg := new(sync.WaitGroup)
			wg.Add(threads)
			for i := 0; i < threads; i++ {
				go worker(hp, wg)
			}
			wg.Wait()

		}

	} else {

		hp, err := MakePipeFile(threads)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg := new(sync.WaitGroup)
		wg.Add(threads)
		for i := 0; i < threads; i++ {
			go worker(hp, wg)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	fmt.Printf("Gohasher took %s", elapsed)

}

func MakePipe(fname string, con int) (hashPipe <-chan string, err error) {
	in, err := os.Open(fname)
	if err != nil {
		return
	}
	thr := con
	hp := make(chan string, thr) // Todo: re-add the -c flag
	hashPipe = hp
	go func(input io.ReadCloser, hp chan<- string) {
		defer input.Close()
		s := bufio.NewScanner(input)
		for s.Scan() {
			hp <- s.Text()
		}
		close(hp)
	}(in, hp)

	return
}

func worker(in <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for it := range in {
		cracker, found := crackersByLength[len(it)]
		if !found {
			fmt.Println("[!!]Unsupported hash", it)
			continue
		}

		res, err := cracker.Crack(it)

		if err != nil {
			fmt.Println(err)
		}
		if res != "" {
			fmt.Printf("Hash(%v): %v value: %v\n", cracker.String(), it, res)
		} else {
			fmt.Printf("Hash(%v): %v value: Not found\n", cracker.String(), it)
		}
	}
}
func MakePipeFile(con int) (hashPipe <-chan string, err error) {
	thr := con
	hp := make(chan string, thr) // Todo: re-add the -c flag
	hashPipe = hp
	go func(hp chan<- string) {

		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			hp <- s.Text()
		}
		close(hp)
	}(hp)

	return
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
	var result1 = make(chan string)
	var result2 = make(chan string)
	var result3 = make(chan string)

	go func() {
		if c.hashType == "md5" {
			result1 <- c.nitrxGen(hash)
		} else if c.hashType == "sha1" {
			result2 <- c.hashToolkit(hash)
		} else {
			result3 <- c.md5Decrypt(hash)
		}

	}()

	select {
	case v1 := <-result1:
		return v1, nil
	case v2 := <-result2:
		return v2, nil
	case v3 := <-result3:
		return v3, nil
	}
}

func NewGeneralCracker(hashType string) *GeneralCracker {
	return &GeneralCracker{hashType}
}
