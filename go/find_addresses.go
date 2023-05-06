package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Address struct {
	Score      int
	Address    string
	PrivateKey string
}

func areLastNCharsSame(str string, n int) bool {
	lastNChars := str[len(str)-n:]
	lastChar := strings.ToLower(string(lastNChars[0]))
	return strings.Count(lastNChars, lastChar) == n
}

func areFirstNCharsSame(str string, n int) bool {
	firstNChars := str[2 : n+2]
	firstChar := strings.ToLower(string(firstNChars[0]))
	return strings.Count(firstNChars, firstChar) == n
}

func hasMoreThanNRepeatedCharsInARow(str string, N int) bool {
	repeated := string(str[0])
	for i := 1; i < len(str); i++ {
		if str[i] == repeated[len(repeated)-1] {
			repeated += string(str[i])
			if len(repeated) == N {
				return true
			}
		} else {
			repeated = string(str[i])
		}
	}
	return false
}

func scoreAddress(address string, numberChars int) int {
	score := 0
	hasRepeatedChars := hasMoreThanNRepeatedCharsInARow(address, numberChars)
	if hasRepeatedChars {
		firstCharsSame := areFirstNCharsSame(address, numberChars)
		lastCharsSame := areLastNCharsSame(address, numberChars)
		if firstCharsSame && lastCharsSame {
			score += 100
		} else if firstCharsSame {
			score += 5
		} else if lastCharsSame {
			score += 10
		} else {
			// score += 1 // uncomment if need chars in the middle
		}
	}
	return score
}

// not works correct
// privateKey := make([]byte, 32)
// _, err := rand.Read(privateKey)
// if err != nil {
// 	panic(err)
// }
// address := computeAddress(privateKey)
// address = strings.ToLower(address)
// func computeAddress(privateKey []byte) string {
// 	hasher := sha256.New()
// 	hasher.Write(privateKey)
// 	hash := hasher.Sum(nil)

// 	// Compute the address from the public key by taking the last 20 bytes
// 	// of the keccak256 hash of the public key.
// 	hasher = sha256.New()
// 	hasher.Write(hash)
// 	publicKeyHash := hasher.Sum(nil)

// 	hasher = sha256.New()
// 	hasher.Write(publicKeyHash)
// 	publicKeyHash = hasher.Sum(nil)

// 	keccakHasher := sha256.New()
// 	keccakHasher.Write(publicKeyHash)
// 	hash = keccakHasher.Sum(nil)

// 	address := hex.EncodeToString(hash[len(hash)-20:])

// 	return address
// }

func generateAddresses(numberChars int, out chan<- Address) {
	counter := 0
	for {
		privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			fmt.Println("Error generating private key:", err)
			return
		}
		// Get the corresponding address
		addressBytes := crypto.PubkeyToAddress(privateKey.PublicKey)
		address := hexutil.Encode(addressBytes[:])
		score := scoreAddress(address, numberChars)
		if score != 0 {
			privateKeyHex := "0x" + hex.EncodeToString(privateKey.D.Bytes())
			out <- Address{
				Score:      score,
				Address:    address,
				PrivateKey: privateKeyHex,
			}
		}

		counter++
		if counter%100000 == 0 {
			fmt.Printf("Processed %d keys\n", counter)
		}
	}
}

func main() {
	numberChars := 4
	numWorkers := 4
	out := make(chan Address, 1000)
	for i := 0; i < numWorkers; i++ {
		go generateAddresses(numberChars, out)
	}

	outputFile, err := os.OpenFile("beauty_addresses.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	counter := 0
	for {
		select {
		case addr := <-out:
			counter++
			fmt.Printf("%+v\n", addr)
			jsonStr, err := json.Marshal(addr)
			if err != nil {
				fmt.Printf("Error marshalling address %+v: %v\n", addr, err)
				continue
			}
			outputFile.Write(jsonStr)
			outputFile.Write([]byte("\n"))
		}
	}
}
