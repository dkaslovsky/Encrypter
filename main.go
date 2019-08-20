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

func encrypt(buf *bytes.Buffer, message []byte, password []byte) (err error) {

	// encryption writer
	encWriter, err := openpgp.SymmetricallyEncrypt(buf, password, nil, nil)
	if err != nil {
		return err
	}
	defer encWriter.Close()

	// encrypt
	_, err = encWriter.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func decrypt(encrypted *bytes.Buffer, password []byte) (decrypted []byte, err error) {

	// function to return password to openpgp.ReadMessage
	getPwd := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return password, nil
	}

	md, err := openpgp.ReadMessage(encrypted, nil, getPwd, nil)
	if err != nil {
		return decrypted, err
	}
	decrypted, err = ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return decrypted, err
	}

	return decrypted, nil
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
	encBuf := bytes.NewBuffer(nil)
	err = encrypt(encBuf, bytesIn, password)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		os.Exit(1)
	}

	// write encrypted file
	err = ioutil.WriteFile(encFile, encBuf.Bytes(), 0644)
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
	encBufIn := bytes.NewBuffer(encBytesIn)
	dec, err := decrypt(encBufIn, password)
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
