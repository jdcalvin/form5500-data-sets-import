package utils

import (
  "fmt"
  "log"
)
 
 // Form5500Flags store cli flags to a struct
type Form5500Flags struct {
  Connection string
	Years []string
	Section string
}

// SetConnection Set attributes from cli inputs
func (f *Form5500Flags) SetConnection(host string, port string, db string, sslmode bool, user string, password string) {

  if db == "" {
    log.Fatal("-db flag is required")
  }

  var connectionPartial string
	
	if user != "" && password != "" {
		connectionPartial = fmt.Sprintf("user=%s password=%s ", user, password)
	} else {
		connectionPartial = ""
	}

	var sslMode string

	if sslmode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
  
  f.Connection = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s %s", host, port, db, sslMode, connectionPartial)
}