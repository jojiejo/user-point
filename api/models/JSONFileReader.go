package models

import (
    "io/ioutil"
    "errors"

    "github.com/jlaffaye/ftp"
)

func ReadJsonFile(filePath string, ftpConn *ftp.ServerConn) ([]byte, error) {
	jsonFile, err := ftpConn.Retr(filePath)
    
    if err != nil {
        err = errors.New("Report not found")
        return nil, err
    }

    defer jsonFile.Close()

    // READ FILE
    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        return nil, err
    }

	return byteValue, nil
}
