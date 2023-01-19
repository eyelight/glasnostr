# Glasnostr
Mine a vanity prefix for your Nostr npub

## Usage
Glasnostr is simple. Just invoke it on the command line with your target prefix. 
```
$ glasnostr foo
```
Or, if the executable is not in your `$PATH`, try this
```
$ ./glasnostr foo
```
You can also specify a limit to the number of attempts. The default is 21 million.
```
$ glasnostr foo 50000
```

## Output
```
$ glasnostr cherrymerkle69
Glasnostr (v0.0.1)
Mine a vanity prefix for your Nostr npub
https://github.com/eyelight/glasnostr

Starting 21000000 attempts for prefix 'cherrymerkle69'

Glasnostr found 'cherrymerkle69' after 420 tries:
    (pub)   adb1deba1bcca55h9775f65e2cb2f6a23b0d9bbe69b9d9fac36efe9ff871eae3
    (sec)   7443d822990o9f4d3ba65183f0251d75c92e31c3e8fedc109ad35e0f0699317e
    (npub)  npub14kcherrymerkle69qe0zevhk5gasmxa7dxuan7krdmlfl7r3at3su5qpl0  
    (nsec)  nsec1w3pasg5fp78563ax2xolqfgawhyjuvwrarldcyy66d0q7p5ex9lquw9d0c

Done hogging your CPU. Thanks for using Glasnostr.
```

## Privacy
As shown above (keys are fake), successful output will display private keys on the screen in plaintext, so you might want to redirect the output into a file using `>` (note that this will overwrite `newfile.txt`).
```
$ glasnostr foo 50000 > newfile.txt
```
You can append to an existing file by redirecting with `>>` (append) instead of `>` (overwrite). 
```
$ glasnostr foo >> exsistingfile.txt
```

## Miscellaneous
The character set for a valid npub is `qpzry9x8gf2tvdw0s3jn54khce6mua7l` so only mine prefixes with these characters. (Yes, 'foo' has been illegal this whole time). 

If you do find 'cherrymerkle69' please let me know on Twitter ([@eyelightyou](https://twitter.com/eyelightyou)).

BTC tips much appreciated. Address is on my github page https://github.com/eyelight