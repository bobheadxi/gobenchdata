package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"

	"go.bobheadxi.dev/gobenchdata/bench"
	"go.bobheadxi.dev/gobenchdata/checks"
)

func statusEmoji(s checks.Status) string {
	switch s {
	case checks.StatusPass:
		return "✅"
	case checks.StatusFail:
		return "❌"
	case checks.StatusNotFound:
		return "❓"
	default:
		return "🤷‍♂️"
	}
}

func newTable(out io.Writer) *tablewriter.Table {
	t := tablewriter.NewWriter(out)
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")
	return t
}

func output(results []bench.Run) {
	if !*noSort {
		sort.Sort(bench.RunHistory(results))
	}
	if *prune > 0 && len(results) > *prune {
		results = results[:*prune]
	}

	var b []byte
	var err error
	if *flat {
		b = make([]byte, 0)
		b = append(b, '[')
		for i, run := range results {
			runBytes, err := json.Marshal(run)
			if err != nil {
				panic(err)
			}
			b = append(b, '\n', ' ', ' ')
			b = append(b, runBytes...)
			if i != (len(results) - 1) {
				b = append(b, ',')
			}
		}
		b = append(b, '\n', ']', '\n')
	} else {
		b, err = json.MarshalIndent(results, "", "  ")
		if err != nil {
			panic(err)
		}
	}
	if *jsonOut != "" {
		if err := ioutil.WriteFile(*jsonOut, b, os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Printf("successfully output results as json to '%s'\n", *jsonOut)
	} else {
		println(string(b))
	}
}

func outputChecksReport(r *checks.Report) {
	// overview conters
	var (
		checksPassed int
		checksFailed int
		checksTotal  int
	)

	results := newTable(os.Stdout)
	results.SetHeader([]string{"", "Check", "Package", "Benchmark", "Diff", "Comment"})
	for checkName, check := range r.Checks {
		checksTotal++
		switch check.Status {
		case checks.StatusPass:
			checksPassed++
		case checks.StatusFail:
			checksFailed++
		}

		t := check.Thresholds
		for _, bench := range check.Diffs {
			var reason string
			if t.Min != nil && bench.Value < *t.Min {
				reason = fmt.Sprintf("exceeded minimum %f (-%.2f)", *t.Min, *t.Min-bench.Value)
			} else if t.Max != nil && bench.Value > *t.Max {
				reason = fmt.Sprintf("exceeded maximum %f (+%.2f)", *t.Max, bench.Value-*t.Max)
			}

			results.Append([]string{
				statusEmoji(bench.Status), checkName, bench.Package, bench.Benchmark,
				fmt.Sprintf("%.2f", bench.Value), reason,
			})
		}
		if len(check.Diffs) == 0 {
			results.Append([]string{
				statusEmoji(checks.StatusNotFound), checkName, "", "",
				"", "no benchmarks found",
			})
		}
	}

	meta := newTable(os.Stdout)
	meta.SetHeader([]string{"", "Base", "Current", "Passed", "Failed", "Total"})
	meta.Append([]string{
		statusEmoji(r.Status), r.Base, r.Current, fmt.Sprintf("%d", checksPassed),
		fmt.Sprintf("%d", checksFailed), fmt.Sprintf("%d", checksTotal),
	})

	// output results to stdout
	meta.Render()
	println()
	results.Render()
	println()
}

func load(files ...string) []bench.RunHistory {
	hist := []bench.RunHistory{}
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		var runs bench.RunHistory
		if err := json.Unmarshal(b, &runs); err != nil {
			panic(err)
		}
		hist = append(hist, runs)
	}
	return hist
}

func showHelp() {
	println(helpText)
	println("FLAGS:\n")
	pflag.PrintDefaults()
	println("\nsee https://go.bobheadxi.dev/gobenchdata for more documentation.")
}
