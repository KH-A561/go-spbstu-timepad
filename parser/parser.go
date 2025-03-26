package parser

import (
	json2 "encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"universityTimepad/model"
)

func ReadTimepad(data []byte, groupId int) ([]model.Timepad, error) {
	strs := strings.Split(string(data), "<script>")
	var json string
	for _, str := range strs {
		str = strings.Trim(str, "\r\n ")
		str, isCut := strings.CutPrefix(str, "window.__INITIAL_STATE__ =")
		if isCut {
			json = strings.Split(str, ";")[0]
		}
	}
	var result []model.Timepad
	get := gjson.Get(json, fmt.Sprintf("lessons.data.%s", strconv.Itoa(groupId)))
	err := json2.Unmarshal([]byte(get.Raw), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
