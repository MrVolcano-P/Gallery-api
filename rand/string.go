package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	return base64.URLEncoding.EncodeToString(b), nil

}

