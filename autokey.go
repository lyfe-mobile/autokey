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
	ak.keylen = utf8.RuneCountInString(key)
	ak.alphalen = utf8.RuneCountInString(alphabet)
	ak.alphabet = []rune(alphabet)
	ak.betaalph = make(map[rune]int)
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

func (ak *AutoKey) MakeSeed() (seed []int, seedstring []rune, seedlen int) {
	if ak.seedlen == 0 {
		seed = []int{0}
		seedstring = []rune{}
		seedlen = 1
		return
	}
	seedlen = ak.seedlen
	seed = make([]int, seedlen)
	seedstring = make([]rune, seedlen)
	for i := 0; i < seedlen; i++ {
		seed[i] = rand.Intn(len(ak.alphabet))
		seedstring[i] = ak.alphabet[seed[i]]
	}
	return
}

func (ak *AutoKey) FindSeed(raw *string) (seed []int, seedlen int) {
	if ak.seedlen == 0 {
		seed = []int{0}
		seedlen = 1
		return
	}
	seedlen = ak.seedlen
	seed = make([]int, seedlen)
	for i := 0; i < seedlen; i++ {
		runeValue, width := utf8.DecodeRuneInString(*raw)
		*raw = (*raw)[width:]
		seed[i] = ak.betaalph[runeValue]
	}
	return
}

func (ak *AutoKey) Encode(raw string) string {
	var (
		res  []rune
		last []int
		n    int
	)
	cc := ak.keylen

	seed, seedstring, seedlen := ak.MakeSeed()
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
	var (
		res  []rune
		last []int
	)
	cc := ak.keylen

	seed, seedlen := ak.FindSeed(&raw)
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
