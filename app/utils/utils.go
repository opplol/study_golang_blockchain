package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	err := encoder.Encode(i)
	HandleErr(err)

	return aBuffer.Bytes()
}

func FromByte(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))

}
