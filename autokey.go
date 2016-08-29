package autokey

type AutoKey struct {
	Key              map[int32]rune
	betaalph         map[rune]int32
	alphabet         []rune
	keylen, alphalen int32
}

func NewAutoKey(key, alphabet string) *AutoKey {
	ak := &AutoKey{}
	ak.Init(key, alphabet)
	return ak
}

func (ak *AutoKey) Init(key, alphabet string) {
	ak.Key = make(map[int32]rune)
	ak.betaalph = make(map[rune]int32)
	for i, c := range key {
		ak.Key[int32(i)] = c
	}
	for i, c := range alphabet {
		ak.betaalph[c] = int32(i)
	}
	ak.alphabet = []rune(alphabet)
	ak.keylen = int32(len(key))
	ak.alphalen = int32(len(alphabet))
}

func (ak *AutoKey) Encode(raw string) string {
	var res []rune
	cc := int32(ak.keylen)
	var last []rune

	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int32
			if cc > 0 {
				k = ak.Key[cc]
			} else {
				k = int32(ak.betaalph[last[0]])
				last = last[1:]
			}
			cc--
			ch := ak.alphabet[(ak.betaalph[c]+k)%ak.alphalen]
			res = append(res, ch)
			last = append(last, ch)
		}
	}
	return string(res)
}

func (ak *AutoKey) Decode(raw string) string {
	var res []rune
	cc := int32(ak.keylen)
	var last []rune

	for _, c := range raw {
		if _, ok := ak.betaalph[c]; !ok {
			res = append(res, c) // Not in alphabet
		} else {
			var k int32
			if cc > 0 {
				k = ak.Key[cc] % ak.alphalen
			} else {
				k = int32(ak.betaalph[last[0]])
				last = last[1:]
			}
			cc--
			ch := ak.alphabet[(ak.betaalph[c]-k+ak.alphalen)%ak.alphalen]
			res = append(res, ch)
			last = append(last, c)
		}
	}
	return string(res)
}

/*
Autokey.prototype.decode = function (raw) {
  var res = ''
  var cc = this.keylen
  var last = []
  var k
  for (var i = 0, ii = raw.length; i < ii; ++i) {
    if (this.betaalph[raw[i]] === undefined) {
      res += raw[i]
    } else {
      if (cc-- > 0) {
        k = this.key[cc] % this.alphalen
      } else {
        k = this.betaalph[last.shift()]
      }
      res += this.alphabet[(this.betaalph[raw[i]] - k + this.alphalen) % this.alphalen]
      last.push(raw[i])
    }
  }
  return res
}

module.exports = Autokey
*/
