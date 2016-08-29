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

var _ = Describe("Autokey", func() {
	const fullAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	Context("symmetric operations", func() {
		const capsAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		const key = "SODFBOIBOEIBOISBDOBLBUVLIUBIEBEHFWF"
		cipher := autokey.NewAutoKey(key, capsAlphabet)

		It("actually encodes", func() {
			Ω(cipher.Encode("WELLDONEISBETTERTHANWELLSAID")).ShouldNot(Equal("WELLDONEISBETTERTHANWELLSAID"))
		})
		It("encodes properly", func() {
			Ω(cipher.Encode("WELLDONEISBETTERTHANWELLSAID")).Should(Equal("WJHQKSOIQTVMEOYSEIOQXWTZTIMR"))
		})
		It("decodes properly", func() {
			Ω(cipher.Decode("WJHQKSOIQTVMEOYSEIOQXWTZTIMR")).Should(Equal("WELLDONEISBETTERTHANWELLSAID"))
		})
	})
	It("is generally reversible", func() {
		cipher := autokey.NewAutoKey("somekey here", fullAlphabet)
		plain := shuffle(fullAlphabet)
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
	It("doesn't touch characters not in the plaintext alphabet", func() {
		cipher := autokey.NewAutoKey("some other key her", fullAlphabet)
		plain := "This is my text!"
		Ω(cipher.Encode(plain)).Should(MatchRegexp(`[A-Za-z0-9]{4} [A-Za-z0-9]{2} [A-Za-z0-9]{2} [A-Za-z0-9]{4}\!`))
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})

	It("works with lots of non-alphabet in it", func() {
		lowerAlphabet := "abcdefghijklmnopqrstuvwxyz"
		cipher := autokey.NewAutoKey("moy", lowerAlphabet)
		plain := "This has a lot of non-alphabet!"
		Ω(cipher.Decode(cipher.Encode(plain))).Should(Equal(plain))
	})
	It("can deal with non-ASCII", func() {
		gAlphabet := "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
		cipher := autokey.NewAutoKey("ქრისტე", gAlphabet)
		plain := "თითქოს კლძეზე ხელი ჰკრეს"
		encoded := cipher.Encode(plain)
		Ω(encoded).Should(Equal("თნჩვქა რღუკძე ოცჭს ცოჭხო"))
		Ω(cipher.Decode(encoded)).Should(Equal(plain))
	})
})

// Ginkgo boilerplate, this runs all tests in this package
func TestCalcSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autokey Tests")
}
