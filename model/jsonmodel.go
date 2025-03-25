package model

import (
	"encoding/json"
	"strconv"
)

type FacultyJson struct {
	Faculty Faculty       `json:"faculty"`
	Data    GroupDataJson `json:"data"`
}

type GroupDataJson struct {
	Data map[string][]Group
}

func (m *GroupDataJson) UnmarshalJSON(b []byte) error {
	resultMap := make(map[string][]Group)
	err := json.Unmarshal(b, &resultMap)
	if err != nil {
		return err
	}
	for k, vList := range resultMap {
		for index := range vList {
			vList[index].FacultyId, _ = strconv.Atoi(k)
		}
	}
	m.Data = resultMap
	return nil
}
