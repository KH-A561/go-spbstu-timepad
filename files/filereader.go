package files

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"universityTimepad/model"
)

const (
	filePathFormat = "%s/%s"
	facultiesDir   = "files/faculties"
)

var (
	facultiesWithGroups *map[*model.Faculty]*[]model.Group
	once                sync.Once
)

func ReadFacultiesWithGroups() (*map[*model.Faculty]*[]model.Group, error) {
	once.Do(func() {
		files, err := os.ReadDir(facultiesDir)
		if err != nil {
			fmt.Println(err)
		}

		result := make(map[*model.Faculty]*[]model.Group)
		for _, entry := range files {
			file, err := os.ReadFile(fmt.Sprintf(filePathFormat, facultiesDir, entry.Name()))
			if err != nil {
				fmt.Println(err)
			}
			var faculties model.FacultyJson
			err = json.Unmarshal(file, &faculties)
			if err != nil {
				fmt.Println(err)
			}

			resK := &faculties.Faculty
			resV := faculties.Data.Data[strconv.Itoa(resK.Id)]
			result[resK] = &(resV)
		}
		facultiesWithGroups = &result
	})

	return facultiesWithGroups, nil
}
