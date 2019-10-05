# gpg-go
Experimenting with Golang GPG.

I found `golang.org/x/crypto/openpgp` to be incompatible with `gnugpg2` in some (potentially corner) cases.  In particular, a PDF encrypted with this code would be mangled when decrypted with `gnugpg2`.  It is possible that the error lies in my implementation, but this concern was enough to abandon developing or using this code for encrypting important files.
