package autokey_test

import (
	"testing"

	"github.com/lyfe-mobile/autokey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Autokey", func() {
	const Possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	Context("symmetric operations", func() {
		const Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		const key = "SODFBOIBOEIBOISBDOBLBUVLIUBIEBEHFWF"
		cipher := autokey.NewAutoKey(key, Alphabet)

		It("encodes properly", func() {
			Ω(cipher.Encode("WELLDONEISBETTERTHANWELLSAID")).Should(Equal("WWUDXFBVDGIZRBLFRVBDKJGMGVZE"))
		})
		It("decodes properly", func() {
			Ω(cipher.Decode("WWUDXFBVDGIZRBLFRVBDKJGMGVZE")).Should(Equal("WELLDONEISBETTERTHANWELLSAID"))
		})
	})
	It("is generally reversible", func() {
		cipher := autokey.NewAutoKey("somekey here", Possible)
		//var str = shuffle(possible)
		str := Possible
		Ω(cipher.Decode(cipher.Encode(str))).Should(Equal(str))
	})
	It("doesn't touch characters not in the plaintext alphabet", func() {
		cipher := autokey.NewAutoKey("some other key her", Possible)
		str := "This is my text!"
		Ω(cipher.Encode(str)).Should(MatchRegexp(`[A-Za-z0-9]{4} [A-Za-z0-9]{2} [A-Za-z0-9]{2} [A-Za-z0-9]{4}\!`))
		Ω(cipher.Decode(cipher.Encode(str))).Should(Equal(str))
	})

	It("works with lots of non-alphabet in it", func() {
		lessPossible := "abcdefghijklmnopqrstuvwxyz"
		cipher := autokey.NewAutoKey("moy", lessPossible)
		str := "This has a lot of non-alphabet!"
		Ω(cipher.Decode(cipher.Encode(str))).Should(Equal(str))
	})

})

// Ginkgo boilerplate, this runs all tests in this package
func TestCalcSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autokey Tests")
}
