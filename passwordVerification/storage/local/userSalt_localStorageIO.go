package local 

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"os"
)

// Stores user salts on a local file
func PutUserSalt(userhash string, salt []byte) {
    data, err := os.ReadFile(LOCAL_FILE_PATH_USER_SALT)
    if err != nil {
        userSaltMap := make(map[string][]byte)
        userSaltMap[userhash] = salt

        writeMapToFile(userSaltMap)

    } else {
        userSaltMap := deserializeFileData(data)

        userSaltMap[userhash] = salt

        writeMapToFile(userSaltMap)
    }
}

// Reads salts on a local file 
func GetUserSalt(userhash string) ([]byte, error) {
    data, err := os.ReadFile(LOCAL_FILE_PATH_USER_SALT)
    if err != nil {
        return nil, err 
    }

    userSaltMap := deserializeFileData(data)

    salt, ok := userSaltMap[userhash]
    if !ok {
        log.Printf("User %s does not exist in the DB!\n", userhash)
        return nil, errors.New("User does not exist in the DB!")
    }

    return salt, nil
}

// Removes a user salt from the storage
func RemoveUserSalt(userhash string) {
    data, err := os.ReadFile(LOCAL_FILE_PATH_USER_SALT)
    if err != nil {
        log.Printf("Failed to read DB file: %s\n", err)
    }
 
    userSaltMap := deserializeFileData(data)

    delete(userSaltMap, userhash)

    writeMapToFile(userSaltMap)
}

func writeMapToFile(userSaltMap map[string][]byte) {
    os.Mkdir(LOCAL_DIR, os.ModePerm)
    file, err := os.Create(LOCAL_FILE_PATH_USER_SALT)
    if err != nil {
        log.Fatalf("failed creating file: %s", err)
    }
    defer file.Close()
        
    b := new(bytes.Buffer)
    e := gob.NewEncoder(b)

    err = e.Encode(userSaltMap)
    if err != nil {
         panic(err)
    }

    file.Write(b.Bytes())
}

func deserializeFileData(data []byte) map[string][]byte {
    b := bytes.NewBuffer(data)
    d := gob.NewDecoder(b)

    var decodedUserSaltMap map[string][]byte
    err := d.Decode(&decodedUserSaltMap)
    if err != nil {
        panic(err)
    }

    return decodedUserSaltMap
}

