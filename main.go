package main

import (
	"net/http"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
	"log"
)

func Decode(raw []byte) []byte {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(raw)))
	base64.StdEncoding.Decode(decoded, raw)
	return decoded
}

func main() {
	cmd := parseCmd()
	client := http.Client{}
	resp, err := client.Get(cmd.gfwListUrl)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(Decode(body)), "\n")
	pattern := `^([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])$`
	compile := regexp.MustCompile(pattern)
	var l []string
	perLine := ""
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
			if perLine != lines[line][2:] {
				l = append(l, "server=/"+lines[line][2:]+"/"+cmd.dns)
			}
			perLine = lines[line][2:]
		} else if lines[line][0] == '.' {
			if perLine != lines[line][1:] {
				l = append(l, "server=/"+lines[line][1:]+"/"+cmd.dns)
			}
			perLine = lines[line][1:]
		} else {
			if perLine != lines[line] {
				l = append(l, "server=/"+lines[line]+"/"+cmd.dns)
			}
			perLine = lines[line]
		}
	}

	fp, err := os.OpenFile(cmd.saveFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	defer fp.Close()

	dnsmasqServer := strings.Join(l, "\n")

	if cmd.v {
		println(dnsmasqServer)
	}

	fp.WriteString(dnsmasqServer)
}
