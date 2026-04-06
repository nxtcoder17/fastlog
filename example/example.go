package main

import (
	"flag"
	"fmt"

	"github.com/nxtcoder17/fastlog"
)

type user struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "--debug")
	flag.Parse()

	attrs := []any{
		"string", "benchmark",
		"number", 17,
		"floating.number", 17.17,
		"bool", true,
		"map",
		map[string]any{"hello": "world"},
		"array",
		[]any{1, 2, 3, 4},
		"err", fmt.Errorf("this is an error"),
		"user",
		user{FirstName: "sample", LastName: "kumar", Age: 17},
		"large-map",
		map[string]any{
			"k1":  "v1",
			"k2":  "v2",
			"k3":  "v3",
			"k4":  "v4",
			"k5":  "v5",
			"k6":  "v6",
			"k7":  "v7",
			"k8":  "v8",
			"k9":  "v9",
			"k10": "v10",
			"k11": "v11",
			"k12": "v12",
			"k13": "v13",
		},
	}

	logfmt := fastlog.New().DebugMode(debug).Logfmt()
	fmt.Printf("# LOGFMT:\n\n")
	logfmt.Debug("hello", attrs...)
	logfmt.Info("hello", attrs...)
	logfmt.Warn("hello", attrs...)
	logfmt.Error("hello", attrs...)

	logfmt.Clone().Logfmt().Info("hello [from clone]", attrs...)

	fmt.Printf("\n# LOGFMT (slog):\n\n")
	logfmtSlog := logfmt.Slog()
	logfmtSlog.Debug("hello", attrs...)
	logfmtSlog.Info("hello", attrs...)
	logfmtSlog.Warn("hello", attrs...)
	logfmtSlog.Error("hello", attrs...)

	consolefmt := fastlog.New().DebugMode(debug).Console()

	fmt.Printf("\n# CONSOLE:\n\n")
	consolefmt.Debug("hello", attrs...)
	consolefmt.Info("hello", attrs...)
	consolefmt.Warn("hello", attrs...)
	consolefmt.Error("hello", attrs...)

	consolefmt.Clone().Console().Info("hello [from clone]", attrs...)

	fmt.Printf("\n# CONSOLE (slog):\n\n")
	consoleFmtSlog := consolefmt.Slog()
	consoleFmtSlog.Debug("hello", attrs...)
	consoleFmtSlog.Info("hello", attrs...)
	consoleFmtSlog.Warn("hello", attrs...)
	consoleFmtSlog.Error("hello", attrs...)

	jsonfmt := fastlog.New().DebugMode(debug).JSON()

	fmt.Printf("\n# JSON:\n\n")
	jsonfmt.Debug("hello", attrs...)
	jsonfmt.Info("hello", attrs...)
	jsonfmt.Warn("hello", attrs...)
	jsonfmt.Error("hello", attrs...)
	jsonfmt.Clone().JSON().Info("hello [from clone]", attrs...)

	fmt.Printf("\n# JSON (slog):\n\n")
	jsonFmtSlog := jsonfmt.Slog()
	jsonFmtSlog.Debug("hello", attrs...)
	jsonFmtSlog.Info("hello", attrs...)
	jsonFmtSlog.Warn("hello", attrs...)
	jsonFmtSlog.Error("hello", attrs...)
}
