package main

import(
	"encoding/json"
	"log"
	"fmt"
)

type Account struct {
    Email string
    Password string
    Money float64
}

func main() {
    account := Account{
        Email: "rsj217@gmail.com",
        Password: "123456",
        Money: 100.5,
    }

    rs, err := json.Marshal(account)
    if err != nil{
        log.Fatalln(err)
    }

    fmt.Println(rs)
    fmt.Println(string(rs))
}