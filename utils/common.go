package utils

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"
)

// PrettyPrint Print json to console elegantly, which is helpful for degbugging
func PrettyPrint(data interface{}, pre string, post string) {
	dataJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf("Data cannot be reformatted, the original data is:\n %#v\n", data)
	} else {
		if pre != "" {
			fmt.Printf("%s\n", pre)
		}
		fmt.Printf("%s\n", dataJSON)
		if post != "" {
			fmt.Printf("%s\n", post)
		}
	}
}

// GetLoginOptions Get login information
func GetLoginOptions() (string, string, string, error) {
	var server, username, password string
	flag.StringVar(&server, "server", "", "Server address")
	flag.StringVar(&username, "username", "admin", "Server login username, admin as default")
	flag.StringVar(&password, "password", "", "Server login password")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		flag.PrintDefaults()
		return "", "", "", errors.New("server, username, and password must all be specified")
	}
	return server, username, password, nil
}

//  URL Compose a URL based on protocol, server address, URI, etc.
func URL(proto string, fqdn string, port string, uri ...string) string {
	var url string

	if proto == "" {
		url = fqdn
	} else {
		url = proto + "://" + fqdn
	}

	if port != "" {
		url += ":" + port
	}

	for _, uri := range uri {
		url += uri
	}

	return url
}

func EmptyStrExists(vars ...string) bool {
	for _, v := range vars {
		if c := strings.TrimSpace(v); c == "" {
			return true
		}
	}
	return false
}
