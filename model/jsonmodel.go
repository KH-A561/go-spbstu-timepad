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

type Timepad struct {
	Weekday int    `json:"weekday"`
	Date    string `json:"date"`
	Lessons []struct {
		Subject        string `json:"subject"`
		SubjectShort   string `json:"subject_short"`
		Type           int    `json:"type"`
		AdditionalInfo string `json:"additional_info"`
		TimeStart      string `json:"time_start"`
		TimeEnd        string `json:"time_end"`
		Parity         int    `json:"parity"`
		TypeObj        struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Abbr string `json:"abbr"`
		} `json:"typeObj"`
		Teachers []struct {
			ID         int    `json:"id"`
			Oid        int    `json:"oid"`
			FullName   string `json:"full_name"`
			FirstName  string `json:"first_name"`
			MiddleName string `json:"middle_name"`
			LastName   string `json:"last_name"`
			Grade      string `json:"grade"`
			Chair      string `json:"chair"`
		} `json:"teachers"`
		Auditories []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Building struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Abbr    string `json:"abbr"`
				Address string `json:"address"`
			} `json:"building"`
		} `json:"auditories"`
		WebinarURL string `json:"webinar_url,omitempty"`
		LmsURL     string `json:"lms_url,omitempty"`
	} `json:"lessons"`
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
