package autokey

import "unicode/utf8"

type AutoKey struct {
	Key              map[int]int
	betaalph         map[rune]int
	alphabet         []rune
	keylen, alphalen int
}

func NewAutoKey(key, alphabet string) *AutoKey {
	ak := &AutoKey{}
	ak.Init(key, alphabet)
	return ak
}

func (ak *AutoKey) Init(key, alphabet string) {
	ak.Key = make(map[int]int)
	ak.betaalph = make(map[rune]int)

	ak.alphalen = utf8.RuneCountInString(alphabet)
	ak.alphabet = []rune(alphabet)
	ak.keylen = utf8.RuneCountInString(key)

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

func (ak *AutoKey) Encode(raw string) string {
	var res []rune
	cc := ak.keylen
	var last []int

	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int
			if cc > 0 {
				k = ak.Key[cc]
			} else {
				k = last[0]
				last = last[1:]
			}
			cc--
			ch := ak.alphabet[(ak.betaalph[c]+k)%ak.alphalen]
			res = append(res, ch)
			last = append(last, ak.betaalph[ch])
		}
	}
	return string(res)
}

func (ak *AutoKey) Decode(raw string) string {
	var res []rune
	cc := ak.keylen
	var last []int

	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int
			if cc > 0 {
				k = ak.Key[cc]
			} else {
				k = last[0]
				last = last[1:]
			}
			cc--
			ch := ak.alphabet[(ak.betaalph[c]-k+ak.alphalen)%ak.alphalen]
			res = append(res, ch)
			last = append(last, ak.betaalph[c])
		}
	}
	return string(res)
}
