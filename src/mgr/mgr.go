package mgr

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

func combrep(n int, lst []string) [][]string {
	if n == 0 {
		return [][]string{nil}
	}
	if len(lst) == 0 {
		return nil
	}
	r := combrep(n, lst[1:])
	for _, x := range combrep(n-1, lst) {
		r = append(r, append(x, lst[0]))
	}
	return r
}

//Mgr is the password cracking manager responsible for the logic of the program
type Mgr struct {
	Mode         int
	Length       int
	CharacterSet string
	WorkGroup    *sync.WaitGroup
}

func concat(idx []int, str []string) string {
	finalStr := ""
	for _, elem := range idx {
		finalStr += str[elem]
	}
	return finalStr
}
func (m *Mgr) dictionaryAttack(dictionary []string, hash string, pass *string) {
	possibleCombinations := combrep(m.Length, dictionary)
	fmt.Println(possibleCombinations)
	hasher := md5.New()

	for _, element := range possibleCombinations {
		a := make([]int, 0)
		for i := range element {
			a = append(a, i)
			N := len(a)
			p := make([]int, N)
			i := 0
			sols := 1
			hasher.Write([]byte(concat(a, element)))

			if hash == hex.EncodeToString(hasher.Sum(nil)) {
				fmt.Println("FOUND PASSWORD:")
				fmt.Println(concat(a, element))
				*pass = concat(a, element)
			}
			hasher.Reset()
			for i < N {

				if p[i] < i {
					j := 0
					if (i % 2) == 1 {
						j = p[i]
					}
					a[j], a[i] = a[i], a[j]
					p[i]++
					i = 1
					sols++
					hasher.Write([]byte(concat(a, element)))

					if hash == hex.EncodeToString(hasher.Sum(nil)) {
						fmt.Println("FOUND PASSWORD:")
						*pass = concat(a, element)
						fmt.Println(concat(a, element))

					}
					hasher.Reset()
				} else {
					p[i] = 0
					i++
				}
			}
		}
	}

}
func (m *Mgr) bruteForce(lexicon string, hash string, pass *string) {
	runtime.GOMAXPROCS(8)
	fmt.Println("the parameters are:", m.Mode, m.Length, m.CharacterSet, hash)
	possibleCombinations := combrep(m.Length,
		strings.Split(lexicon, " "))
	fmt.Println(len(possibleCombinations))

	numOfRoutines := 8
	start := 0
	increment := len(possibleCombinations) / numOfRoutines
	for i := 0; i < numOfRoutines; i++ {
		m.WorkGroup.Add(1)
		if i == (numOfRoutines - 1) {
			go crackpass(possibleCombinations[start:], hash, m.WorkGroup, pass)

		} else {
			go crackpass(possibleCombinations[start:start+increment], hash, m.WorkGroup, pass)

		}
		start += increment
	}

}

//Attack is the function responsible of launching the cracking attack
func (m *Mgr) Attack(hash string, dictionary []string, pass *string) {
	fmt.Println("dictionary is :", dictionary)
	lowercase := "a b c d e f g h i j k l m n o p q r s t u v w x y z"
	uppercase := "A B C D E F G H I J K L M N O P Q R S T U V W X Y Z"
	numbers := "0 1 2 3 4 5 6 7 8 9"
	specialChars := "$ @ &  "
	var lexicon string

	if strings.Contains(m.CharacterSet, "l") {
		lexicon += lowercase + " "
	}
	if strings.Contains(m.CharacterSet, "u") {
		lexicon += uppercase + " "
	}
	if strings.Contains(m.CharacterSet, "n") {
		lexicon += numbers + " "
	}
	if strings.Contains(m.CharacterSet, "s") {
		lexicon += specialChars + " "
	}
	if m.Mode == 1 {
		//fmt.Println(lexicon)
		fmt.Println("Initiating Brute Force Attack...")

		m.bruteForce(lexicon, hash, pass)
	} else if m.Mode == 2 {
		m.dictionaryAttack(dictionary, hash, pass)
	} else if m.Mode == 3 {
		fmt.Println("not yet implemented")
	}
}

func crackpass(possibleCombinations [][]string, myhash string, wg *sync.WaitGroup, pass *string) {
	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()
	//fmt.Println("routine is up and running")
	//*pass = "hello"
	hasher := md5.New()
	for _, element := range possibleCombinations {
		// fmt.Println(strings.Join(element,""))

		// Perm([]rune(strings.Join(element, "")), func(a []rune) {

		// 	hasher.Write([]byte(string(a)))

		// 	if myhash == hex.EncodeToString(hasher.Sum(nil)) {
		// 		fmt.Println("FOUND PASSWORD:")
		// 		fmt.Println(string(a))

		// 	}
		// 	hasher.Reset()
		// 	// fmt.Println(string(a))
		// })
		a := []rune(strings.Join(element, ""))
		N := len(a)
		p := make([]int, N)
		i := 0
		sols := 1
		hasher.Write([]byte(string(a)))

		if myhash == hex.EncodeToString(hasher.Sum(nil)) {
			fmt.Println("FOUND PASSWORD:")
			fmt.Println(string(a))
			*pass = string(a)
		}
		hasher.Reset()
		for i < N {

			if p[i] < i {
				j := 0
				if (i % 2) == 1 {
					j = p[i]
				}
				a[j], a[i] = a[i], a[j]
				p[i]++
				i = 1
				sols++
				hasher.Write([]byte(string(a)))

				if myhash == hex.EncodeToString(hasher.Sum(nil)) {
					fmt.Println("FOUND PASSWORD:")
					*pass = string(a)
					fmt.Println(string(a))

				}
				hasher.Reset()
			} else {
				p[i] = 0
				i++
			}
		}

	}
	//fmt.Println("go routine finished at", time.Now().Format(time.RFC850))
	wg.Done()
}
