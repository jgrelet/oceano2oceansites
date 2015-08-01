// writeAsciiBtl
package main

import (
	"fmt"
)

func (nc *Btl) WriteAscii(map_format map[string]string, hdr []string) {
	fmt.Println("Hello from WriteAscii for bottle !")
}

func (nc *Btl) WriteHeader(map_format map[string]string, hdr []string) {
	fmt.Println("Hello from WriteHeader for bootle !")
}
