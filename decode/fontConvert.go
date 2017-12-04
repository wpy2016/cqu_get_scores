package decode

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
)

func Decode(str string) (string, error) {
	reader := bytes.NewReader([]byte(str))
	uftReader := transform.NewReader(reader, simplifiedchinese.GBK.NewDecoder())
	result, e := ioutil.ReadAll(uftReader)
	if e != nil {
		log.Printf("Decode : error=[%v]", e)
		return "", e
	}
	return string(result), nil
}
