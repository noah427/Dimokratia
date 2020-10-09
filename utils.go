package main

import (
	"github.com/imroc/req"

	"encoding/base64"
	"fmt"
	"regexp"
)

var (
	findFileTypeRegex = regexp.MustCompile(`.+\.(.+)`)
)

func urlToDataScheme(url string) string {
	resp, _ := req.Get(url)
	scheme := fmt.Sprintf("data:image/%s;base64,%s", findFileType(url), base64.StdEncoding.EncodeToString(resp.Bytes()))
	return scheme
}

func findFileType(url string) string {
	res := findFileTypeRegex.FindAllStringSubmatch(url, -1)
	return res[0][1]
}

func Has(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func findActionType(actionType string) ActionType {
	for _, action := range actionTypes {
		if action.name == actionType {
			return action
		}
	}
	return ActionType{}
}