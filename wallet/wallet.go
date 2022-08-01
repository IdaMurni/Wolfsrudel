package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"../utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	wallet := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wallet.privateKey = privateKey
	wallet.publicKey = &wallet.privateKey.PublicKey
	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(wallet.publicKey.X.Bytes())
	h2.Write(wallet.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(dc8)
	wallet.blockchainAddress = address
	return wallet
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", wallet.privateKey.D.Bytes())
}

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func (wallet *Wallet) BlockchainAddress() string {
	return wallet.blockchainAddress
}

func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        wallet.PrivateKeyStr(),
		PublicKey:         wallet.PublicKeyStr(),
		BlockchainAddress: wallet.BlockchainAddress(),
	})
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (transaction *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.senderBlockchainAddress,
		Recipient: transaction.recipientBlockchainAddress,
		Value:     transaction.value,
	})
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func (transactionRequest *TransactionRequest) Validate() bool {
	if transactionRequest.SenderPrivateKey == nil ||
		transactionRequest.SenderBlockchainAddress == nil ||
		transactionRequest.RecipientBlockchainAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil {
		return false
	}
	return true
}
