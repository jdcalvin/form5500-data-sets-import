package utils

import (
  "fmt"
)

func ExampleMappingStruct() {
  m := Mapping{LongForm: "TESTLONG", ShortForm: "TESTSHORT", Alias: "TESTALIAS", DataType: "TESTDATATYPE"}
  fmt.Println(m.LongForm)
  fmt.Println(m.ShortForm)
  fmt.Println(m.Alias)
  fmt.Println(m.DataType)
  fmt.Println(m.IndexName())
  // Output:
  // TESTLONG
  // TESTSHORT
  // TESTALIAS
  // TESTDATATYPE
  // idx_TESTALIAS
}