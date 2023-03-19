package digester

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get the DigesterInfo for a directory
func Directory(dirPath string) ([]DigestInfo, error) {
	var files []DigestInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {

			digestInfo, err := File(path, info)
			if err != nil {
				return err
			}
			files = append(files, digestInfo)
		} else {
			// fmt.Println("Directory: ", path)
			if path != dirPath {
				fmt.Println("skip Directory: ", path)
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// jsonBytes, err := json.Marshal(files)
	// if err != nil {
	// 	return nil, err
	// }

	// return jsonBytes, nil
	return files, nil
}
