# encrypter
Experimenting with Golang GPG.

I found this code to be incompatible with `gnugpg2` in some (potentially corener) cases.  In particular, a PDF encrypted with this code would be mangled when decrypted with `gnugpg2`.
