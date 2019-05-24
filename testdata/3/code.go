package code

import "regexp"

const c = "0-9"

var re = regexp.MustCompile(`[`+c+``)
