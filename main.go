package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

const (
	version       string = "v0.0.2"
	npubPrefix    string = "npub1"
	bech32charset string = "023456789acdefghjklmnpqrstuvwxyz"
)

func isBech32(s string) bool {
	for _, c := range s {
		if !strings.Contains(bech32charset, string(c)) {
			return false
		}
	}
	return true
}

func main() {
	var target string
	var tries string
	var limit int
	var hits uint

	// whats up
	fmt.Printf("Glasnostr (%s)\nMine a vanity prefix for your Nostr npub\nhttps://github.com/eyelight/glasnostr\n\n", version)

	// handle no target
	if len(os.Args) < 2 {
		fmt.Printf("Please supply Glasnostr with a target and optionally an upper limit for the number of guesses (default: 21 million guesses)\n	example: $ glasnostr foo 50000\n\nAlso, private keys will be sent to the screen, so you may want to redirect the output\n	example: $ glasnostr foo 50000 > glasnostr.txt\n")
		os.Exit(1)
	}

	// validate target
	target = os.Args[1]
	if !isBech32(target) {
		fmt.Printf("Error: target '%s' is invalid\nThe valid character set for an encoded npub is bech32. Try again with only the following characters: \n\n	%s\n\n", target, bech32charset)
		os.Exit(1)
	}

	// set limit
	if len(os.Args) < 3 {
		limit = 21000000
	} else {
		tries = os.Args[2]
		l, e := strconv.ParseUint(tries, 10, 64)
		if e != nil {
			fmt.Printf("Cannot parse limit (stay positive): \n%s\n\n", e.Error())
			os.Exit(1)
		}
		limit = int(l)
	}

	fmt.Printf("Starting %d attempts for prefix '%s'\n\n", limit, target)

	// find target
	for i := 0; i < limit; i++ {
		// generate keys
		sk := nostr.GeneratePrivateKey()
		pk, pkerr := nostr.GetPublicKey(sk)
		if pkerr != nil {
			fmt.Printf("Error getting public key: %s\n", pkerr.Error())
			continue
		}

		// make the npub from the pubkey
		npub, nerr := nip19.EncodePublicKey(pk)
		if nerr != nil {
			fmt.Printf("Error encoding NIP-19 npub from public key: %s\n", nerr.Error())
			continue
		}

		// evaluate the target against the npub & report if found
		result := evaluate(npub, target)
		if !result {
			continue
		} else {
			hits++
			nsec, err := nip19.EncodePrivateKey(sk)
			if err != nil {
				fmt.Printf("Error encoding NIP-19 nsec from secret key: %s\n", err.Error())
			}
			fmt.Printf("Glasnostr found '%s' after %d tries:\n	(pub)	%s\n	(sec)	%s\n	(npub)	%s\n	(nsec)	%s\n\n", target, i, pk, sk, npub, nsec)
		}
	}
	if hits == 0 {
		fmt.Printf("¯\\_(ツ)_/¯\n")
	} else {
		fmt.Printf("Done hogging your CPU. Thanks for using Glasnostr.\n")
	}
}

func evaluate(key, target string) bool {
	target = npubPrefix + target
	return strings.HasPrefix(key, target)
}
