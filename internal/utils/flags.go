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
func (f *Form5500Flags) SetConnection(host string, port string, db string, sslmode string, user string, password string) {

  if db == "" {
    log.Fatal("-db flag is required")
  }

  var connectionPartial string
	
	if user != "" && password != "" {
		connectionPartial = fmt.Sprintf("user=%s password=%s ", user, password)
	} else {
		connectionPartial = ""
	}

	
	if !(sslmode == "require" || sslmode == "verify-full" || sslmode == "verify-ca" || sslmode == "disable") {
		log.Fatal("-ssl only supports disable (default), 'verify-full', 'verify-ca', 'require'")
	}
  
  f.Connection = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s %s", host, port, db, sslmode, connectionPartial)
}