package main

import "fmt"


import "crypto/sha1" 
import "encoding/hex"


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
func main() {
    Perm([]rune("abcd"), func(a []rune) {
        fmt.Println(createHash(string(a)))
    })
    fmt.Println("Hello fatouma")
    fmt.Println(createHash("hello fatouma"))
    fmt.Println(createHash("hello fatouma"))
}
