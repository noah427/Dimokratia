package main


import (
	"github.com/imroc/req"

	"encoding/base64"
	"fmt"
)


func urlToDataScheme(url string) string {
	resp,_ := req.Get(url)
	scheme := fmt.Sprintf("data:image/png;base64,%s",base64.StdEncoding.EncodeToString(resp.Bytes()))
	return scheme
}


func Has(array []string, value string) bool {
	for _,v := range array {
		if v == value{
			return true
		}
	}

	return false
}