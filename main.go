package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

const (
	version           string = "v0.0.3"
	npubPrefix        string = "npub1"
	bech32charset     string = "023456789acdefghjklmnpqrstuvwxyz"
	ERR_SEARCH_FAILED string = "did not find prefix"
)

var (
	misses    uint64
	target    string
	limit     uint64
	startTime time.Time
	unbounded bool
)

type keypair struct {
	npub     string
	nsec     string
	pub      string
	sec      string
	attempts uint64
	dur      time.Duration
}

func isBech32(s string) bool {
	for _, c := range s {
		if !strings.Contains(bech32charset, string(c)) {
			return false
		}
	}
	return true
}

func main() {
	// var tries string
	mult := time.Second

	// whats up
	fmt.Printf("Glasnostr (%s)\nMine a vanity prefix for your Nostr npub\nhttps://github.com/eyelight/glasnostr\n\n", version)

	// handle no target
	if len(os.Args) < 2 {
		fmt.Printf("Please supply Glasnostr with a target and optionally an upper limit for the number of guesses (default: 21 million guesses, or use '0' for infinity)\n	example: $ glasnostr foo 50000\n\nAlso, private keys will be sent to the screen, so you may want to redirect the output\n	example: $ glasnostr foo 50000 > glasnostr.txt\n")
		os.Exit(1)
	}

	// validate target
	target = os.Args[1]
	if !isBech32(target) {
		fmt.Printf("Target '%s' is invalid\nThe valid character set for an encoded npub is bech32. Try again with only the following characters: \n\n	%s\n\n", target, bech32charset)
		os.Exit(1)
	}

	// set limit
	if len(os.Args) < 3 {
		limit = 21000000
	} else {
		tries := os.Args[2]
		l, e := strconv.ParseUint(tries, 10, 64)
		if e != nil {
			fmt.Printf("Cannot parse limit (stay positive): \n%s\n\n", e.Error())
			os.Exit(1)
		}
		limit = l
		if limit == 0 {
			unbounded = true
		}
	}

	numWorkers := runtime.GOMAXPROCS(runtime.NumCPU())
	success := make(chan bool, 1)
	cancelAll := make(chan struct{}, numWorkers)
	var wg sync.WaitGroup
	startTime = time.Now()
	if !unbounded {
		fmt.Printf("(%d attempts): ", limit)
	}
	fmt.Printf("Starting %d threads looking for prefix '%s' at %s\n\n", numWorkers, target, startTime.Local().Format(time.RFC822))

	// set up workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go work(i, &wg, success, cancelAll)
	}

Work:
	for {
		select {
		case <-success:
			// stop
			for i := 0; i < numWorkers; i++ {
				cancelAll <- struct{}{}
			}
			break Work
		default:
			// report progress
			if !unbounded {
				if misses >= limit {
					fmt.Printf("¯\\_(ツ)_/¯ Tried %d attempts.\n", misses)
					for i := 0; i < numWorkers; i++ {
						cancelAll <- struct{}{}
					}
					break Work
				}
			}
			if time.Since(startTime) > mult && misses%uint64(100000) == 0 {
				fmt.Printf("%d tries after %s...\n", misses, time.Since(startTime).Round(time.Second))
				mult *= 2
			}
		}
	}
	wg.Wait()
	fmt.Printf("Done hogging your CPU. Thanks for using Glasnostr.\n")
}

// evaluate returns true if `target` immediately follows 'npub1' in the key
func evaluate(key, target string) bool {
	target = npubPrefix + target
	return strings.HasPrefix(key, target)
}

// mine returns a keypair if a generated npub matches the target, or an error
func mine() (keypair, error) {
	// generate keys
	sk := nostr.GeneratePrivateKey()
	pk, pkerr := nostr.GetPublicKey(sk)
	if pkerr != nil {
		fmt.Printf("error getting public key: %s\n", pkerr.Error())
		return keypair{}, pkerr
	}

	// make the npub from the pubkey
	npub, nerr := nip19.EncodePublicKey(pk)
	if nerr != nil {
		fmt.Printf("error encoding NIP-19 npub from public key: %s\n", nerr.Error())
		return keypair{}, nerr
	}

	// evaluate the target against the npub; finalize & return keypair if found
	result := evaluate(npub, target)
	if !result {
		return keypair{}, errors.New(ERR_SEARCH_FAILED)
	} else {
		nsec, err := nip19.EncodePrivateKey(sk)
		if err != nil {
			fmt.Printf("error encoding NIP-19 nsec from secret key: %s\n", err.Error())
			return keypair{}, err
		}
		kp := keypair{
			npub:     npub,
			nsec:     nsec,
			sec:      sk,
			pub:      pk,
			attempts: misses,
			dur:      time.Since(startTime),
		}
		return kp, nil
	}
}

// work mines for keys and reports upon success, until receiving a cancelAll channel
func work(id int, wg *sync.WaitGroup, success chan<- bool, cancelAll <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-cancelAll:
			return
		default:
			kp, err := mine()
			if err == nil {
				fmt.Printf("Glasnostr found '%s' after %d tries (%s):\n	(pub)	%s\n	(sec)	%s\n	(npub)	%s\n	(nsec)	%s\n\n", target, kp.attempts, kp.dur, kp.pub, kp.sec, kp.npub, kp.nsec)
				success <- true
			}
			misses++
		}
	}
}
