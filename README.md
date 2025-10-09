## Install Golang

Download golang : `https://go.dev/dl/`

cek instalasi : `go version`

buat file dengan nama 'main.go' dg isi code seperti ini:

```
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}

```

Kemudian jalankan: 

```
go run main.go

```

## inisialisi Modul GO
Setiap project Go butuh modul (semacam package.json kalau di Node.js).

jalankan ini:

```
go mod init task-manager

```

setelah itu nanti akan muncul file: `go.mod`