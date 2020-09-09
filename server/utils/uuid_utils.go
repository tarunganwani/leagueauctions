package utils

import (
	"github.com/google/uuid"
)

//GetUUIDFromString - get uuid from string
func GetUUIDFromString(uuidstr string) (uuidOut uuid.UUID, err error){
	uuidOut, err = uuid.Parse(uuidstr)
	return
}

//GetStringFromUUID - get string from uuid
func GetStringFromUUID(uuidIn uuid.UUID) string{
	return uuidIn.String()
}