package main
 
import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateSecretKey() string {
	key := make([]byte, 32) // 32 bytes é um bom tamanho para a chave de assinatura HS256

	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(key)
}

func main() {
	fmt.Println(generateSecretKey())
}
