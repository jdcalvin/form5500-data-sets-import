package utils

import (
  "os"
  "io/ioutil"
  "log"
)

// ReadFile will read a filepath from form5500/ and return the string, or error if file is not found
func ReadFile(filepath string) (string) {
  pwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }

  txt, err := ioutil.ReadFile(pwd+"/"+filepath)
  if err != nil {
    log.Fatal(err)
  }

  return string(txt)
}

