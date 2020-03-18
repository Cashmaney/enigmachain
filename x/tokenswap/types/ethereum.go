package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"gopkg.in/yaml.v2"
)

const (
	EthereumAddressLength = 20
	EthereumTxHashLength  = 32
)

//
// Address is a common interface for different types of addresses used by the SDK
type EthIf interface {
	Equals(EthIf) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var _ EthIf = EthereumAddress{}

//var (
//	hashT    = reflect.TypeOf(EthereumTxHash{})
//	addressT = reflect.TypeOf(EthereumAddress{})
//)
//
//type EthereumAddress [AddressLength]byte
// var _ sdk.Address = EthereumAddress{}

var _ yaml.Marshaler = EthereumAddress{}

type EthereumAddress [EthereumAddressLength]byte

func (a EthereumAddress) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}

func (a EthereumAddress) Equals(other EthIf) bool {
	if a.Empty() && other.Empty() {
		return true
	}

	return bytes.Equal(a.Bytes(), other.Bytes())
}

func (a EthereumAddress) Empty() bool {
	if a.Bytes() == nil {
		return true
	}

	aa2 := EthereumAddress{}
	return bytes.Equal(a.Bytes(), aa2.Bytes())
}

func (a EthereumAddress) Marshal() ([]byte, error) {
	return a[:], nil
}

func (a EthereumAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (a *EthereumAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := HexToAddress(s)
	if err != nil {
		return err
	}

	*a = aa2
	return nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf compatibility.
func (a *EthereumAddress) Unmarshal(data []byte) error {
	if len(data) != EthereumAddressLength {
		return errors.New("invalid ethereum address length")
	}
	copy(a[:], data)
	return nil
}

func (a EthereumAddress) Bytes() []byte {
	return a[:]
}

func (a EthereumAddress) Format(s fmt.State, verb rune) {
	_, _ = fmt.Fprintf(s, "%"+string(verb), a[:])
}

type EthereumTxHash [EthereumTxHashLength]byte

func BytesToTxHash(b []byte) EthereumTxHash {
	var a EthereumTxHash
	a.SetBytes(b)
	return a
}

func BytesToAddress(b []byte) EthereumAddress {
	var a EthereumAddress
	a.SetBytes(b)
	return a
}

func (h *EthereumTxHash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-EthereumTxHashLength:]
	}
	copy(h[EthereumTxHashLength-len(b):], b)
}

func (a *EthereumAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-EthereumAddressLength:]
	}
	copy(a[EthereumAddressLength-len(b):], b)
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func HexToAddress(s string) (addr EthereumAddress, err error) {
	addrBytes := FromHex(s)
	if len(addrBytes) != EthereumAddressLength {
		return addr, errors.New("invalid length")
	}
	return BytesToAddress(addrBytes), nil
}

func HexToTxHash(s string) EthereumTxHash {
	return BytesToTxHash(FromHex(s))
}

func (a EthereumAddress) String() string {
	return a.Hex()
}
func (h EthereumTxHash) String() string {
	return h.Hex()
}

// Hex returns an EIP55-compliant hex string representation of the address.
func (a EthereumAddress) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

func (h EthereumTxHash) Hex() string { return EncodeHex(h[:]) }

// Encode encodes b as a hex string with 0x prefix.
func EncodeHex(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}
