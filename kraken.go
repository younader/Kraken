package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"kraken/src/mgr"

	"github.com/urfave/cli"
)

func testHash(key string) {
	hasher := md5.New()
	hasher.Write([]byte(key))

	if "719d5e9b60811224f0c9366a9cd55023" == hex.EncodeToString(hasher.Sum(nil)) {
		fmt.Println("FOUND PASSWORD:")
		fmt.Println(string(key))

	}
}
func main() {
	var wg = sync.WaitGroup{}
	fmt.Println(time.Now().Format(time.RFC850))
	var mode string
	var len string
	var charset string
	app := cli.NewApp()
	app.Name = "kraker"
	app.Usage = "strongest password cracking tool on the planet"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "mode, m",
			Value:       "1",
			Usage:       "attack mode: 1 for brute force, 2 for dictionary, 3 for hybrid",
			Destination: &mode,
		},
		cli.StringFlag{
			Name:        "charset, c",
			Value:       "luns",
			Usage:       "character set: l for lowercase, u for uppercase, n for numbers and s for special chars",
			Destination: &charset,
		},
		cli.StringFlag{
			Name:        "len, l",
			Value:       "5",
			Usage:       "password length",
			Destination: &len,
		},
	}

	app.Action = func(c *cli.Context) error {
		//fmt.Println("Hello friend!", c.Args().Get(0))
		// fmt.Println("attack mode", c.String("mode"))
		// fmt.Println("character set", c.String("charset"))
		// fmt.Println("len", c.String("len"))

		if c.NArg() == 0 {
			log.Fatal("No hash file supplied")
		}
		//fmt.Println(mode, len)
		filename := c.Args().Get(0)
		dictname := c.Args().Get(1)
		//mt.Println(filename, dictname)
		stream, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		readHash := string(stream)
		readDict := make([]string, 0)
		fmt.Println("hash to crack: ", readHash)
		modeI, _ := strconv.Atoi(mode)
		lenI, _ := strconv.Atoi(len)
		//fmt.Println("mode :  length: ", modeI, lenI)

		if mode != "1" {
			stream, err := ioutil.ReadFile(dictname)
			fmt.Println(dictname)
			readDict = append(readDict, strings.Split(string(stream), "\n")...)
			if err != nil {
				log.Fatal(err)
			}
		}
		//fmt.Println(readDict)
		mymgr := mgr.Mgr{Mode: modeI, Length: lenI, CharacterSet: charset, WorkGroup: &wg}
		mymgr.Attack(readHash, readDict)
		return nil

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()

}
