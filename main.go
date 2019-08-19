package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/crypto/openpgp"
)

type Opts struct {
	FileName string `long:"filename" short:"f" required:"true" description:"Name of file to encrypt"`
}

var opts Opts
var parser = flags.NewParser(&opts, flags.Default)

func encrypt(message []byte, password []byte) (enc []byte, err error) {

	buf := bytes.NewBuffer(nil)

	// encryption writer
	encWriter, err := openpgp.SymmetricallyEncrypt(buf, password, nil, nil)
	if err != nil {
		return enc, err
	}

	// encrypt
	_, err = encWriter.Write(message)
	if err != nil {
		return enc, err
	}
	err = encWriter.Close()
	if err != nil {
		return enc, err
	}

	return buf.Bytes(), nil
}

func decrypt(enc []byte, password []byte) (dec []byte, err error) {

	// function to return password to openpgp.ReadMessage
	getPwd := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return password, nil
	}

	buf := bytes.NewBuffer(enc)
	md, err := openpgp.ReadMessage(buf, nil, getPwd, nil)
	if err != nil {
		return dec, err
	}
	dec, err = ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return dec, err
	}

	return dec, nil
}

func main() {
	_, err := parser.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fileIn := opts.FileName
	encFile := fileIn + ".gpg"
	decFile := fileIn + "_dec"

	password := []byte("gopher")

	// read whole file at once
	bytesIn, err := ioutil.ReadFile(fileIn)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	// encrypt
	enc, err := encrypt(bytesIn, password)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		os.Exit(1)
	}

	// write encrypted file
	err = ioutil.WriteFile(encFile, enc, 0644)
	if err != nil {
		fmt.Println("Error writing encrypted file:", err)
		os.Exit(1)
	}

	// read encrypted file
	encBytesIn, err := ioutil.ReadFile(encFile)
	if err != nil {
		fmt.Println("Error reading encrypted file:", err)
		os.Exit(1)
	}

	// decrypt
	dec, err := decrypt(encBytesIn, password)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		os.Exit(1)
	}

	// write decrypted file
	err = ioutil.WriteFile(decFile, dec, 0644)
	if err != nil {
		fmt.Println("Error writing decrypted file:", err)
		os.Exit(1)
	}

}
