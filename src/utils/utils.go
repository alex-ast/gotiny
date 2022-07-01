package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

// Removed symbols '0', 'l', '1', 'I', 'O', 'o' to avoid reading mistakes
const letters = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

var prng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenRandomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = letters[prng.Int()%len(letters)]
	}
	return string(bytes)
}

func Marshal(obj any) *bytes.Buffer {
	data, _ := json.Marshal(obj)
	return bytes.NewBuffer(data)
}

func Unmarshal(data io.ReadCloser, obj any) error {
	bytes, err := ioutil.ReadAll(data)
	if err == nil {
		err = json.Unmarshal(bytes, obj)
	}
	return err
}

// Marshals obj to the bytes array
func MarshalToBytes(obj any) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}

// Unmarshals object from the byte array
func UnmarshalFromBytes(bytes []byte, obj any) error {
	return json.Unmarshal(bytes, obj)
}

// panic() if err is not nil
func Assert(err error) {
	if err != nil {
		panic(err)
	}
}
