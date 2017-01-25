package utils

import (
  "log"
)

// ReadAsset will read a filepath from form5500/ and return the string, or error if file is not found
func ReadAsset(filepath string) (string) {
  data, err := Asset("assets/"+filepath)
  if err != nil {
    log.Fatal(err)
  }

  return string(data)
}

