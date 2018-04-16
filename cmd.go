package main

import "flag"

type Cmd struct {
	gfwListUrl string
	dns        string
	saveFile   string
	v          bool
}

func parseCmd() Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.gfwListUrl, "url", "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt", "gfw list url")
	flag.StringVar(&cmd.dns, "dns", "208.67.220.220#5353", "dns")
	flag.StringVar(&cmd.saveFile, "file", "./dnsmasq.servers.conf", "save file")
	flag.BoolVar(&cmd.v, "v", false, "v")
	flag.Parse()
	return cmd
}
