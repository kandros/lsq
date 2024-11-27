package config

import "strings"

type Config struct {
	CfgVers      int    `edn:"meta/version"`
	PreferredFmt string `edn:"preferred-format"`
	FileNameFmt  string `edn:"journal/file-name-format"`
}

func ConvertDateFormat(cfgFileFormat string) string {
	lsqFmts := [][]string{
		{"yyyy", "2006"},
		{"yy", "06"},
		{"MM", "01"},
		{"M", "1"},
		{"dd", "02"},
		{"d", "2"},
	}

	goFormat := cfgFileFormat
	for _, val := range lsqFmts {
		goFormat = strings.ReplaceAll(goFormat, val[0], val[1])
	}

	return goFormat
}
