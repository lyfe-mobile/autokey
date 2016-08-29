# Autokey

[![Build Status](https://travis-ci.org/lyfe-mobile/autokey.svg?branch=master)](https://travis-ci.org/lyfe-mobile/autokey)
[![Coverage Status](https://coveralls.io/repos/github/lyfe-mobile/autokey/badge.svg?branch=master)](https://coveralls.io/github/lyfe-mobile/autokey?branch=master)

A Go implementation of the simple
[Autokey cipher](https://en.wikipedia.org/wiki/Autokey_cipher).
Autokey has some useful characteristics: the output is the same length
as the input and looks like the input was merely jumbled as it uses
the same alphabet as the input:

* Input: `all our words from loose using have lost their edge.`
* Output: `kbu cjd hprjy qmud ilujd xpxkq xxsw hvng zgksj nqwa.`

This is similar to and based on the stream cipher implementation at
[https://github.com/hex7c0/autokey](https://github.com/hex7c0/autokey)
but restricts the input and output to a given alphabet.  This means
you can let spaces, punctuation, and other characters go through the
algorithm with no change.

## How it works

To install:

    go get github.com/lyfe-mobile/autokey

## Example usage

```go
package main

import "github.com/lyfe-mobile/autokey"

const Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func main() {
    // The key does not have to be chosen from the alphabet.
    akey := autokey.NewAutoKey("My Key Is This@)(#*$", Alphabet)

    input := "Something very important!";
    cipherText := akey.Encode(input)
    println(cipherText)
    recoveredText := akey.Decode(cipherText)
    println(recoveredText)
}
```

## Caveat

Autokey is a very old cipher with known vulnerabilities and is not as
secure as modern encryption systems.  You should consider this
algorithm as a slightly better alternative to
[ROT13](https://en.wikipedia.org/wiki/ROT13) or
[base64](https://en.wikipedia.org/wiki/Base64) for obfuscating text.
