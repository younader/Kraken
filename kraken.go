package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Equanox/gotron"
	"github.com/urfave/cli"
	"kraken/src/mgr"
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
		var temp string
		mymgr := mgr.Mgr{Mode: modeI, Length: lenI, CharacterSet: charset, WorkGroup: &wg}
		mymgr.Attack(readHash, readDict, &temp)
		return nil

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	window, err := gotron.New("./webapp")
	if err != nil {
		panic(err)
	}

	// Alter default window size and window title.
	window.WindowOptions.Width = 1200
	window.WindowOptions.Height = 980
	window.WindowOptions.Title = "Kraken"

	// Start the browser window.
	// This will establish a golang <=> nodejs bridge using websockets,
	// to control ElectronBrowserWindow with our window object.
	done, err := window.Start()

	if err != nil {
		panic(err)
	}
	opt := option.WithCredentialsFile("./serviceAccount.json")
	myapp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}
	client, err := myapp.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for true {
		res, err := client.Collection("req_queue").Doc("req").Get(context.Background())

		for err != nil {
			time.Sleep(15000)
			res, err = client.Collection("req_queue").Doc("req").Get(context.Background())
		}
		data := res.Data()
		_, err = client.Collection("req_queue").Doc("req").Delete(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("document data: %#v\n", data)
		fmt.Println(data["mode"])
		fmt.Println(data["mode"])
		dic := make([]string, 0)
		s := reflect.ValueOf(data["dictionary"])
		for i := 0; i < s.Len(); i++ {
			dic = append(dic, s.Index(i).Elem().String())
		}

		Mode, _ := strconv.Atoi(data["mode"].(string))
		Length, _ := strconv.Atoi(data["length"].(string))

		mymgr := mgr.Mgr{Mode: Mode, Length: Length, CharacterSet: data["charset"].(string), WorkGroup: &wg}
		var pass string
		mymgr.Attack(data["hash"].(string), dic, &pass)
		wg.Wait()

		fmt.Println("the result is", pass)

		_, err = client.Collection("results").Doc(data["hash"].(string)).Set(context.Background(), map[string]interface{}{
			"password": pass,
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}

	}
	defer client.Close()
	// Open dev tools must be used after window.Start
	// window.OpenDevTools()

	// Wait for the application to close
	<-done
}
