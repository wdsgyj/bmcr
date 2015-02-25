package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

var fileFlag = flag.String("crash", "", "Crash File")

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	regx, err = regexp.Compile(`\{([^\s\{\}=/w]{3,}?)=`)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	file, err := os.Open(*fileFlag)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	crashs := make([]*Crash, 0, 512)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			crash, err := New(text)
			if err != nil {
				log.Fatalln(err)
			}
			crashs = append(crashs, crash)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	//	for _, c := range crashs {
	//		log.Println(c.Feature)
	//	}
}
