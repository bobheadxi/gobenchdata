import { Run, RunSuiteBenchmark, RunSuite } from "@/generated";

// copied from 'apexcharts.ApexAxisChartSeries'
type ApexAxisChartSingleSeries = {
  name?: string;
  type?: string;
  data: { x: number; y: number | null }[];
};

type ApexAxisChartSeries = ApexAxisChartSingleSeries[];

enum MetricBuiltins {
  NSPEROP = "NsPerOp",
  MEM_BYPTESPEROP = "Mem.BytesPerOp",
  MEM_ALLOCSPEROP = "Mem.AllocsPerOp",
  MEM_MBPERS = "Mem.MBPerSec",
}

type ChartSet = { [metric: string]: ApexAxisChartSeries };

const defaultMetrics = {
  [MetricBuiltins.NSPEROP]: true,
  [MetricBuiltins.MEM_BYPTESPEROP]: true,
  [MetricBuiltins.MEM_ALLOCSPEROP]: true,
};

export function iterateSuites(
  runs: Run[],
  func: (s: RunSuite, r: Run) => void
) {
  for (let rID = 0; rID < runs.length; rID += 1) {
    const run = runs[rID];
    // iterate suites in each run
    for (let sID = 0; sID < run.Suites.length; sID += 1) {
      func(run.Suites[sID], run);
    }
  }
}

export function iterateBenchmarks(
  runs: Run[],
  func: (b: RunSuiteBenchmark, s: RunSuite, r: Run) => void,
  pkg?: RegExp
) {
  iterateSuites(runs, (suite: RunSuite, run: Run) => {
    if (pkg && !pkg.test(suite.Pkg)) return;

    // iterate benchmarks
    for (let bID = 0; bID < suite.Benchmarks.length; bID += 1) {
      func(suite.Benchmarks[bID], suite, run);
    }
  });
}

function benchName(b: string): string {
  return b.replace(/^(Benchmark)/, "");
}

/**
 * Generates one set ApexAxisChartSeries per metric for the provided group of pkg, benches
 *
 * @param runs
 * @param pkg
 * @param benches
 * @param metrics
 */
export function generateSeries(
  runs: Run[],
  pkg: RegExp,
  benches: RegExp[],
  metrics: { [metric: string]: boolean } = defaultMetrics
): {
  charts: ChartSet;
  xaxis: { [metric: string]: number[] };
} {
  // by default, include all builtins
  if (Object.keys(metrics).length === 0) {
    metrics = defaultMetrics;
  }
  const metricKeys = Object.keys(metrics).filter((m: string) => metrics[m]);

  // set up index:
  // * each metric creates a chart
  // * each chart gets a set of series, corresponding to matching benchmarks
  const index = metricKeys.reduce((acc: ChartSet, cur) => {
    acc[cur] = [];
    return acc;
  }, {});

  // we need to pad data, so track x axis of each metric
  const xaxis = metricKeys.reduce(
    (acc: { [metric: string]: { [x: number]: boolean } }, cur) => {
      acc[cur] = {};
      return acc;
    },
    {}
  );

  // check each run for suites
  iterateBenchmarks(
    runs,
    (bench, _, run) => {
      for (let i = 0; i < benches.length; i += 1) {
        const benchMatch = benches[i];
        if (!benchMatch.test(bench.Name)) continue;

        // add benchmark data
        for (let mID = 0; mID < metricKeys.length; mID += 1) {
          const metric = metricKeys[mID];
          const existingSeries = index[metric]
            .filter((v) => v.name === benchName(bench.Name))
            .pop();
          if (!existingSeries) {
            index[metric].push({
              name: benchName(bench.Name),
              data: [],
            });
          }

          // get existing series or get the one we just made
          const series =
            existingSeries || index[metric][index[metric].length - 1];

          // add appropriate value
          const push = (y: number) => {
            series.data.push({ x: run.Date, y });
            xaxis[metric][run.Date] = true;
          };
          switch (metric) {
            case MetricBuiltins.NSPEROP:
              push(bench.NsPerOp);
              break;
            case MetricBuiltins.MEM_ALLOCSPEROP:
              push(bench.Mem.AllocsPerOp);
              break;
            case MetricBuiltins.MEM_BYPTESPEROP:
              push(bench.Mem.BytesPerOp);
              break;
            case MetricBuiltins.MEM_MBPERS:
              push(bench.Mem.MBPerSec);
              break;
            default:
              // assume custom if metric is not a builtin
              if (bench.Custom && metric in bench.Custom) {
                push(bench.Custom[metric]);
              }
          }
        }
      }
    },
    pkg
  );

  // sort out x axis values for each metric
  const xaxisArrays: { [metric: string]: number[] } = {};
  metricKeys.forEach((metric) => {
    xaxisArrays[metric] = Object.keys(xaxis[metric]).map((k) => parseInt(k));
    xaxisArrays[metric].sort();
  });

  // fill missing data for each metric
  metricKeys.forEach((metric) => {
    // index the points each series has data for
    const seriesX = index[metric].reduce(
      (acc: { [s: string]: { [x: number]: boolean } }, series) => {
        acc[series.name || "?"] = series.data.reduce(
          (acc2: { [x: number]: boolean }, point) => {
            acc2[point.x] = true;
            return acc2;
          },
          {}
        );
        return acc;
      },
      {}
    );
    // generate missing data
    const seriesNames = Object.keys(seriesX);
    xaxisArrays[metric].forEach((x) => {
      seriesNames.forEach((s) => {
        if (!seriesX[s][x]) {
          index[metric]
            .filter((series) => series.name === s)
            .pop()
            ?.data.push({ x, y: null });
        }
      });
    });
    // sort
    index[metric].forEach((s) => {
      s.data.sort((p1, p2) => (p1.x < p2.x ? -1 : 1));
    });
  });

  return {
    charts: index,
    xaxis: xaxisArrays,
  };
}
