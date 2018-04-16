package main

import (
	"net/http"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
)

func Decode(raw []byte) []byte {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(raw)))
	base64.StdEncoding.Decode(decoded, raw)
	return decoded
}

func main() {
	cmd := parseCmd()
	client := http.Client{}
	resp, _ := client.Get(cmd.gfwListUrl)
	by, _ := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(Decode(by)), "\n")
	pattern := `^([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])$`
	compile := regexp.MustCompile(pattern)
	var l []string
	for line := range lines {
		if len(lines[line]) == 0 {
			continue
		} else if lines[line][0] == '!' {
			continue
		} else if lines[line][0:2] == "@@" {
			continue
		} else if strings.Contains(lines[line], "/") || strings.Contains(lines[line], "*") || strings.Contains(lines[line], "[") || strings.Contains(lines[line], "%") || !strings.Contains(lines[line], ".") {
			continue
		} else if compile.MatchString(lines[line]) {
			continue
		} else if lines[line][0:2] == "||" {
			l = append(l, lines[line][2:])
		} else if lines[line][0] == '.' {
			l = append(l, lines[line][1:])
		} else {
			l = append(l, lines[line])
		}
	}
	fp, _ := os.OpenFile(cmd.saveFile, os.O_CREATE|os.O_TRUNC, 0755)
	defer fp.Close()
	perLine := ""
	for e := range l {
		if perLine != l[e] {
			println(l[e])
			fp.WriteString("server=/" + l[e] + "/" + cmd.dns + "\n")
		}
		perLine = l[e]
	}
}
