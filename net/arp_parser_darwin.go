package net

import "regexp"

var ArpTableParser = regexp.MustCompile("^[^\\d\\.]+([\\d\\.]+).+\\s+([a-f0-9:]{11,17})\\s+on\\s+([^\\s]+)\\s+.+$")
var ArpTableTokens = 4
