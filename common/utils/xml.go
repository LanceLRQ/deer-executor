package utils

import (
	"encoding/xml"
	"strings"
)

// XMLStringObject 解析XML字符串
func XMLStringObject(xmlStr string, obj interface{}) bool {
	// for testlib
	xmlStr = strings.Replace(xmlStr, "<?xml version=\"1.0\" encoding=\"windows-1251\"?>", "<?xml version=\"1.0\" encoding=\"utf-8\"?>", -1)
	return XMLBytesObject([]byte(xmlStr), obj)
}

// XMLBytesObject 解析XML字节数组
func XMLBytesObject(xmlStr []byte, obj interface{}) bool {
	err := xml.Unmarshal(xmlStr, &obj)
	return err == nil
}
