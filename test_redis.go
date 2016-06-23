package main

import (
    redis "gopkg.in/redis.v3"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // client := redis.NewClient(&redis.Options{
    //     Addr: "192.168.1.100:7001",
    // })

    db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/dbname")
    if err != nil  {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT username, password FROM uc_user")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    for rows.Next() {
        var name, password string
        if err := rows.Scan(&name, &password); err != nil {
            panic(err)
        }

        fmt.Printf("%-40s %-100s\n", name, password)
    }

    if err = rows.Err(); err != nil {
        panic(err)
    }

    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{"127.0.0.1:7001"},
        Password: "password",
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    err = client.Set("test", "Haha", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := client.Get("test").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println(val)

    err = client.LPush("message_queue", "Hello, world").Err()
    if err != nil {
        panic(err)
    }

    val, _ = client.LPop("message_queue").Result()
    fmt.Println(val)

}
