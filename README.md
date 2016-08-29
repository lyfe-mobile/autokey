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
    akey := autokey.NewAutoKey("MyKeyIsThis", Alphabet, 2)

    input := "Something very important!";
    cipherText := akey.Encode(input)
    println(cipherText)
    recoveredText := akey.Decode(cipherText)
    println(recoveredText)
}
```

`NewAutoKey` takes a key and an alphabet as the first two
arguments. The key should be chosen from the given alphabet.  

The third argument to `NewAutoKey` is a seed length. When plaintext
consists of many smaller strings with a very high number of
duplicates, the regular Autokey cipher will spit out many duplicates
on its own. This could allow a viewer to infer values based on the
frequency of the duplicated strings.

Using a seed will make this more difficult. A new seed is determined
randomly for each call to `Encode()` and is then tacked on the front
of the emitted string. This means that output will be larger than input.

A 0 may be given as the seed length to use the regular Autokey
algorithm with no random seed.

## Caveat

Autokey is a very old cipher with known vulnerabilities and is not as
secure as modern encryption systems.  You should consider this
algorithm as a slightly better alternative to
[ROT13](https://en.wikipedia.org/wiki/ROT13) or
[base64](https://en.wikipedia.org/wiki/Base64) for obfuscating text.
