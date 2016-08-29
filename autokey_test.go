package autokey_test

import (
	"math/rand"
	"testing"

	"github.com/lyfe-mobile/autokey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func shuffle(s string) string {
	r := []rune(s)
	for i := range r {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

const (
	fullAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	SeedLen      = 2
)

var _ = Describe("Autokey", func() {
	BeforeEach(func() {
		rand.Seed(100)
	})
	Context("symmetric operations", func() {
		const capsAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		const key = "SODFBOIBOEIBOISBDOBLBUVLIUBIEBEHFWF"
		cipher := autokey.NewAutoKey(key, capsAlphabet, SeedLen)

		It("actually encodes", func() {
			Ω(cipher.Encode("WELLDONEISBETTERTHANWELLSAID")).ShouldNot(Equal("WELLDONEISBETTERTHANWELLSAID"))
		})
		It("encodes properly", func() {
			Ω(cipher.Encode("WELLDONEISBETTERTHANWELLSAID")).Should(Equal("PQLZWGZIDYFJKCTENITYDGMMIPIYBH"))
		})
		It("decodes properly", func() {
			Ω(cipher.Decode("PQLZWGZIDYFJKCTENITYDGMMIPIYBH")).Should(Equal("WELLDONEISBETTERTHANWELLSAID"))
		})
	})
	It("is generally reversible", func() {
		cipher := autokey.NewAutoKey("somekey here", fullAlphabet, SeedLen)
		plain := shuffle(fullAlphabet)
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
	It("doesn't touch characters not in the plaintext alphabet", func() {
		cipher := autokey.NewAutoKey("some other key her", fullAlphabet, SeedLen)
		plain := "This is my text!"
		Ω(cipher.Encode(plain)).Should(MatchRegexp(`[A-Za-z0-9]{4} [A-Za-z0-9]{2} [A-Za-z0-9]{2} [A-Za-z0-9]{4}\!`))
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
	It("works with lots of non-alphabet in it", func() {
		lowerAlphabet := "abcdefghijklmnopqrstuvwxyz"
		cipher := autokey.NewAutoKey("moy", lowerAlphabet, SeedLen)
		plain := "This has a lot of non-alphabet!"
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
	It("can deal with non-ASCII", func() {
		gAlphabet := "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
		cipher := autokey.NewAutoKey("ქრისტე", gAlphabet, SeedLen)
		plain := "თითქოს კლძეზე ხელი ჰკრეს"
		encoded := cipher.Encode(plain)
		Ω(cipher.Decode(encoded)).Should(Equal(plain))
	})
	It("can deal with zero seedlen", func() {
		cipher := autokey.NewAutoKey("somekey here", fullAlphabet, 0)
		plain := shuffle(fullAlphabet)
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
})

// Ginkgo boilerplate, this runs all tests in this package
func TestCalcSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autokey Tests")
}
