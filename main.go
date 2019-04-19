package main

import "fmt"


import "crypto/sha1" 
import "encoding/hex"
import "strings"

func Perm(a []rune, f func([]rune)) {
    perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []rune, f func([]rune), i int) {
    if i > len(a) {
        f(a)
        return
    }
    perm(a, f, i+1)
    for j := i + 1; j < len(a); j++ {
        a[i], a[j] = a[j], a[i]
        perm(a, f, i+1)
        a[i], a[j] = a[j], a[i]
    }
}
func createHash(key string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
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
func main() {
    myhash:=createHash("fatma")
    // Perm([]rune("abcd"), func(a []rune) {
    //     fmt.Println(createHash(string(a)))
    // })
    // fmt.Println("Hello fatouma")
    // fmt.Println(createHash("hello fatouma"))
    // fmt.Println(createHash("hello fatouma"))
    possible_combinations:=combrep(5,
        strings.Split("a b c d e f g h i j k l m n o p q r s t u v w x y z 0 1 2 3 4 5 6 7 8 9 $ @ &", " "))
      
        for index,element:=range possible_combinations{
        //fmt.Println(strings.Join(element,""))
       found:=0 
       
        Perm([]rune(strings.Join(element,"")), func(a []rune) {
            if myhash==createHash(string(a)){
                fmt.Println("FOUND PASSWORD:")
                 fmt.Println(string(a))
                 
                 
            }
        
            // fmt.Println(string(a))
        })

        if found==1{fmt.Println(index)}
    }    
}
