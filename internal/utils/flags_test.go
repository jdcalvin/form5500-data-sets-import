package utils

import (
  "fmt"
)

var port, host, dbName, user, password, sslmode string
var flagEx Form5500Flags

func setVariables() {
  host = "testhost"
  port = "5555"
  sslmode = "disable"
  user = ""
  password = ""
  dbName = "testdb"
}

// SetConnection(host string, port string, db string, sslmode bool, user string, password string)

func ExampleSetConnectionWithoutUserPassword() {
  setVariables()
  
  flagEx.SetConnection(host,port,dbName,sslmode,user,password)
  fmt.Println(flagEx.Connection) // no user or password
  // Output:
  // host=testhost port=5555 dbname=testdb sslmode=disable
}

func ExampleSetConnectionWithoutPassword() {
  setVariables()
  
  flagEx.SetConnection(host,port,dbName,sslmode,"testuser",password)
  fmt.Println(flagEx.Connection) // no user or password
  // Output:
  // host=testhost port=5555 dbname=testdb sslmode=disable
}

func ExampleSetConnectionWithUserPassword() {
  setVariables()
  
  flagEx.SetConnection(host,port,dbName,sslmode,"testuser", "abc123")
  fmt.Println(flagEx.Connection) // no user or password
  // Output:
  // host=testhost port=5555 dbname=testdb sslmode=disable user=testuser password=abc123
}
