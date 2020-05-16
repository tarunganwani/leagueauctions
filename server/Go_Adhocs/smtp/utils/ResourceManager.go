package utils

import (
	"github.com/magiconair/properties"
)

var PropertyFiles = []string{"${GOPATH}/src/smtp/resources/mail.properties"}

var Props, _ = properties.LoadFiles(PropertyFiles, properties.UTF8, true)

type ResourceManager struct {
}

func (res ResourceManager) GetProperty(propertyName string) string {
	message := ""
	var ok bool
	message, ok = Props.Get(propertyName)
	if !ok {
		return Props.MustGet("not.found")
	} else {
		return message
	}
}
