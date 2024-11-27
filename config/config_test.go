package config

import "testing"

func TestConvertDateFormat(t *testing.T) {
	tests := map[string]struct {
		cfgFileFormat string
		want          string
	}{
		"yyyy_MM_dd (Logseq Default)": {
			cfgFileFormat: "yyyy_MM_dd",
			want:          "2006_01_02",
		},
		"yyyy_MM_d": {
			cfgFileFormat: "yyyy_MM_d",
			want:          "2006_01_2",
		},
		"yyyy_M_d": {
			cfgFileFormat: "yyyy_M_d",
			want:          "2006_1_2",
		},
		"yyyy_M_dd": {
			cfgFileFormat: "yyyy_M_dd",
			want:          "2006_1_02",
		},
		"yy_MM_dd": {
			cfgFileFormat: "yy_MM_dd",
			want:          "06_01_02",
		},
		"yy_M_dd": {
			cfgFileFormat: "yy_M_dd",
			want:          "06_1_02",
		},
		"yy_M_d": {
			cfgFileFormat: "yy_M_d",
			want:          "06_1_2",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := ConvertDateFormat(tt.cfgFileFormat); got != tt.want {
				t.Errorf("ConvertDateFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
