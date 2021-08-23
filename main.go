package main

import "fmt"

//go:generate go build -tags 'main evm java' -x -o libshare_signer_main.so -ldflags '-s -w' -buildmode=c-shared  signer/mobile
//go:generate go build -tags 'local evm java' -x -o libshare_signer_test.so -ldflags '-s -w' -buildmode=c-shared  signer/mobile
//go:generate go.exe build -tags 'main evm java' -x -o share_signer_main.dll -ldflags '-s -w' -buildmode=c-shared signer/mobile
//go:generate go.exe build -tags 'local evm java' -x -o share_signer_test.dll -ldflags '-s -w' -buildmode=c-shared signer/mobile
func main() {
	fmt.Println("")
}
