package config

type Config struct {
	CfgVers      int    `edn:"meta/version"`
	PreferredFmt string `edn:"preferred-format"`
	FileNameFmt  string `edn:"journal/file-name-format"`
}
