# GO-MODELS

## USAGE
Install the library with `go get github.com/frenki123/go-models`

In `main.go` create new struct User:
```go
type User struct {
	Id       int
	Username string
	Email    string
}
```
Next we need to define schema for this struct:
```go
func (u User) Schema() db.Table {
	return db.MustDefSchema(&u,
		db.CharField(&u.Username, 25).Unique(),
		db.EmailField(&u.Email).Unique(),
	)
}
```
If *db* is not imported from github.com/frenki123/go-models/db, import it manually
```go
import "github.com/frenki123/go-models/db"

```
To start to use this table we need to register them in main function. Also we need to import the driver which we need. Here we will use go-sqlite3. Go-models is using .env file to read DB config.
```go
import (
	"github.com/frenki123/go-models/db"
	_ "github.com/mattn/go-sqlite3"
)
....
func main() {
	db.MustRegister(User{})
}
```
Create `.env` file with following data:
```.env
DATABASE = "sqlite3"
DBCONNSTRING = "example.db"
```
To create database we can run `go run .`.
Then we can save some users in db. If you created db file, delete it before.
```go
user1 := User{Username: "user1", Email: "email"}
user2 := User{Username: "username", Email: "email"}
if err := db.Save(&user1); err != nil {
	fmt.Println("Error saving in database", err)
}
if err := db.Save(&user2); err != nil {
	fmt.Println("Error saving in database", err)
}
```
And get them:
```go
user, err := db.Get(&User{}, 1)
fmt.Println("Error getting user is:", err)
fmt.Println("User is", user)
```
Full code
```go
package main

import (
	"fmt"
	"github.com/frenki123/go-models/db"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       int
	Username string
	Email    string
}

func (u User) Schema() db.Table {
	return db.MustDefSchema(&u,
		db.CharField(&u.Username, 25).Unique(),
		db.EmailField(&u.Email).Unique(),
	)
}

func main() {
	db.MustRegister(User{})
    user1 := User{Username: "user1", Email: "email"}
    user2 := User{Username: "username", Email: "email"}
    if err := db.Save(&user1); err != nil {
    	fmt.Println("Error saving in database", err)
    }
    if err := db.Save(&user2); err != nil {
    	fmt.Println("Error saving in database", err)
    }
    user, err := db.Get(&User{}, 1)
    fmt.Println("Error getting user is:", err)
    fmt.Println("User is", user)
}
```
