package autokey

import (
	"math/rand"
	"unicode/utf8"
)

type AutoKey struct {
	Key                       map[int]int
	betaalph                  map[rune]int
	alphabet                  []rune
	keylen, alphalen, seedlen int
}

func NewAutoKey(key, alphabet string, seedlen int) *AutoKey {
	ak := &AutoKey{}
	ak.Init(key, alphabet, seedlen)
	return ak
}

func (ak *AutoKey) Init(key, alphabet string, seedlen int) {
	ak.Key = make(map[int]int)
	ak.betaalph = make(map[rune]int)

	ak.alphalen = utf8.RuneCountInString(alphabet)
	ak.alphabet = []rune(alphabet)
	ak.keylen = utf8.RuneCountInString(key)
	ak.seedlen = seedlen
	i := 0
	for _, c := range alphabet {
		ak.betaalph[c] = i
		i++
	}
	i = 0
	for _, c := range key {
		ak.Key[i] = ak.betaalph[c]
		i++
	}
}

func (ak *AutoKey) RandSeed() ([]int, []rune) {
	seed := make([]int, ak.seedlen)
	seedstring := make([]rune, ak.seedlen)
	for i := 0; i < ak.seedlen; i++ {
		seed[i] = rand.Intn(len(ak.alphabet))
		seedstring[i] = ak.alphabet[seed[i]]
	}
	return seed, seedstring
}

func (ak *AutoKey) Encode(raw string) string {
	var (
		res        []rune
		last       []int
		n          int
		seed       []int
		seedstring []rune
		seedlen    int
	)
	cc := ak.keylen

	if ak.seedlen == 0 {
		seed = []int{0}
		seedstring = []rune{}
		seedlen = 1
	} else {
		seed, seedstring = ak.RandSeed()
		seedlen = ak.seedlen
	}
	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int
			if cc > 0 {
				k = ak.Key[cc]
				cc--
			} else {
				k = last[0]
				last = last[1:]
			}
			ch := ak.alphabet[(ak.betaalph[c]+k+seed[n%seedlen])%ak.alphalen]
			res = append(res, ch)
			last = append(last, ak.betaalph[ch])
			n++
		}
	}
	return string(seedstring) + string(res)
}

func (ak *AutoKey) Decode(raw string) string {
	var res []rune
	cc := ak.keylen
	var last []int
	var seed []int
	var seedstring []rune
	var seedlen int

	if ak.seedlen == 0 {
		seed = []int{0}
		seedstring = []rune{}
		seedlen = 1
	} else {
		seedlen = ak.seedlen
		seed = make([]int, seedlen)
		seedstring = make([]rune, seedlen)
		for i := 0; i < seedlen; i++ {
			runeValue, width := utf8.DecodeRuneInString(raw)
			raw = raw[width:]
			seed[i] = ak.betaalph[runeValue]
			seedstring[i] = runeValue
		}
	}

	n := 0
	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int
			if cc > 0 {
				k = ak.Key[cc]
				cc--
			} else {
				k = last[0]
				last = last[1:]
			}
			ch := ak.alphabet[(ak.betaalph[c]-(k+seed[n%seedlen])+ak.alphalen*2)%ak.alphalen]
			res = append(res, ch)
			last = append(last, ak.betaalph[c])
			n++
		}
	}
	return string(res)
}
