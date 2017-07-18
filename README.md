# secret-share
a secret sharing tool based on Shamir's Secret Sharing algorithm implemented with pure Golang

## About
Shamir's Secret Sharing is an algorithm in cryptography created by Adi Shamir. 
It is a form of secret sharing, where a secret is divided into parts, giving each participant its own unique part, where some of the parts or all of them are needed in order to reconstruct the secret.

Counting on all participants to combine the secret might be impractical, and therefore sometimes the threshold scheme is used where any ```k``` of the parts are sufficient to reconstruct the original secret.

This is a small tool based on Shamir's Secret Sharing algorithm implemented by [codahale](https://github.com/codahale/sss), which can create up to 255 share parts from one file or directory,
and recover from some of the parts (not less then the given param ```k```)
