package utils

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
)

func PrettyPrint(data interface{}, pre string, post string) {
	dataJson, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf("Data cannot be reformatted, the original data is:\n %#v\n", data)
	} else {
		if pre != "" {
			fmt.Printf("%s\n", pre)
		}
		fmt.Printf("%s\n", dataJson)
		if post != "" {
			fmt.Printf("%s\n", post)
		}
	}
}

func GetLoginOptions() (string, string, string, error) {
	var server, username, password string
	flag.StringVar(&server, "server", "", "Server address")
	flag.StringVar(&username, "username", "admin", "Server login username, admin as default")
	flag.StringVar(&password, "password", "", "Server login password")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		flag.PrintDefaults()
		return "", "", "", errors.New("server, username, and password must all be specified.")
	} else {
		return server, username, password, nil
	}
}
