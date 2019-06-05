package kron

import (
  "os"
  "fmt"
  "github.com/spf13/viper"
)

// For storing the location of a job for select and favorite.
type location struct {
	Contexts []string  `json:"contexts"`
	Namespaces []string  `json:"namespaces"`
}

func toLocation(obj interface{}) location {
  m, ok := obj.(map[string]interface{})
  if !ok {
    fmt.Println("Failed")
    return location{} // maybe panic??
  }
  c := m["contexts"].([]string)
  n := m["namespaces"].([]string)
  return location{c, n}
}

func createConfig() {
  os.Mkdir(os.Getenv("HOME") + "/.kron/", 0777)
  err := viper.WriteConfigAs(os.Getenv("HOME") + "/.kron/config.yaml")
  if err != nil {
    panic(err.Error())
  }
}
