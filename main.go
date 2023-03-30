package main

import (
    "bufio"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "log"
    "os"
)

func main() {
    // Open the input file
    file, err := os.Open("transactions.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Read the transactions into a slice
    scanner := bufio.NewScanner(file)
    var transactions []string
    for scanner.Scan() {
        transactions = append(transactions, scanner.Text())
    }

    // Compute the Merkle Tree Root
    root := computeMerkleRoot(transactions)
    fmt.Println(root)
}

func computeMerkleRoot(transactions []string) string {
    // Convert the transactions to byte arrays
    var txBytes [][]byte
    for _, tx := range transactions {
        txBytes = append(txBytes, hexToBytes(tx))
    }

    // Compute the Merkle Tree
    tree := buildMerkleTree(txBytes)

    // Return the root hash
    return bytesToHex(tree[0])
}

func buildMerkleTree(data [][]byte) [][]byte {
    if len(data) == 0 {
        return [][]byte{}
    }

    // Special case for leaf nodes
    if len(data) == 1 {
        hash := sha256.Sum256(data[0])
        return [][]byte{hash[:]}
    }

    // Compute the left and right subtrees
    var left, right [][]byte
    if len(data)%2 == 1 {
        left = append(data[:len(data)-1], data[len(data)-1])
        right = left
    } else {
        left = data[:len(data)/2]
        right = data[len(data)/2:]
    }

    leftTree := buildMerkleTree(left)
    rightTree := buildMerkleTree(right)

    // Concatenate the left and right subtrees
    var tree [][]byte
    tree = append(tree, leftTree...)
    tree = append(tree, rightTree...)

    // Compute the parent nodes
    for i := 0; i < len(tree)-1; i += 2 {
        parent := sha256.Sum256(append(tree[i], tree[i+1]...))
        tree = append(tree, parent[:])
    }

    return tree
}

func hexToBytes(hexStr string) []byte {
    bytes, err := hex.DecodeString(hexStr)
    if err != nil {
        log.Fatal(err)
    }
    return bytes
}

func bytesToHex(bytes []byte) string {
    return hex.EncodeToString(bytes)
}
