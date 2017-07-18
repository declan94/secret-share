# secret-share
a secret sharing tool based on Shamir's Secret Sharing algorithm implemented with pure Golang

## About
Shamir's Secret Sharing is an algorithm in cryptography created by Adi Shamir. 
It is a form of secret sharing, where a secret is divided into parts, giving each participant its own unique part, where some of the parts or all of them are needed in order to reconstruct the secret.

Counting on all participants to combine the secret might be impractical, and therefore sometimes the threshold scheme is used where any ```k``` of the parts are sufficient to reconstruct the original secret.

This is a small tool based on Shamir's Secret Sharing algorithm implemented by [codahale](https://github.com/codahale/sss), which can create up to 255 share parts from one file or directory,
and recover from some of the parts (not less then the given param ```k```)

[More about Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir's_Secret_Sharing)

## Install

#### Install from source code
You have to install go tools first, see [here](https://golang.org/doc/install#install).

After that, execute commands below:
```
$ go get -u github.com/declan94/secret-share
$ sudo cp `go env GOPATH`/bin/secret-share /usr/local/bin
```

#### Install from tarball
For macos or linux, you can download [Release Tarball](https://github.com/declan94/secret-share/releases), unzip it then directly use the pre-built binary program.

## Usage

#### Create sharing parts
```
$ secret-share -k 3 path/to/secretfile path/to/part1 path/to/part2 path/to/part3 path/to/part4
```
This will create 4 sharing parts (part1, part2, part3, part4), and at least 3 of them are needed to recover the original secretfile.

When ```-k``` is not specified, all created sharing parts are needed to recover the original secretfile.

When secretfile is a directory, all sharing parts will be directory with same hierarchical structure.

#### Recover secret file (directory)
```
$ secret-share -r path/to/recover path/to/part1 path/to/part2 path/to/part3
```
This will recover original file or directory from the given sharing parts.

