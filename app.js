/**
 * Set up ChartsJS configuration
 */
function initChartsJS() {
  Chart.defaults['line-with-guides'] = Chart.defaults.line;
  Chart.controllers['line-with-guides'] = Chart.controllers.line.extend({
    draw(ease) {
      Chart.controllers.line.prototype.draw.call(this, ease);

      if (this.chart.tooltip._active && this.chart.tooltip._active.length) {
        const activePoint = this.chart.tooltip._active[0];
        const { ctx } = this.chart;
        const { x } = activePoint.tooltipPosition();
        const topY = this.chart.scales['y-axis-0'].top;
        const bottomY = this.chart.scales['y-axis-0'].bottom;

        // draw line
        ctx.save();
        ctx.beginPath();
        ctx.moveTo(x, topY);
        ctx.lineTo(x, bottomY);
        ctx.lineWidth = 2;
        ctx.strokeStyle = '#666';
        ctx.stroke();
        ctx.restore();
      }
    },
  });
}

/**
 * Color definitions
 */
const chartColors = {
  red: 'rgb(255, 99, 132)',
  orange: 'rgb(255, 159, 64)',
  yellow: 'rgb(255, 205, 86)',
  green: 'rgb(75, 192, 192)',
  blue: 'rgb(54, 162, 235)',
  purple: 'rgb(153, 102, 255)',
  grey: 'rgb(201, 203, 207)',
};
const chartColorsList = Object.values(chartColors);

/**
 * Generate a random color ID
 */
function randomColorID() {
  return Math.floor(Math.random() * (chartColorsList.length + 1));
}

/**
 * Retrieves color of given ID
 *
 * @param {number} i color ID
 */
function getColor(i) {
  return chartColorsList[i % chartColorsList.length];
}

/**
 * Generate an identifier for a chart
 *
 * @param {string} pkg package identifier
 * @param {string} chart chart identifier
 */
const makeChartName = (pkg, chart) => `${pkg}-${chart}`;

/**
 * Instantiates a chart point from a run
 *
 * @param {{ Date: number }} run point represents a run
 * @param {number} val benchmark value
 */
const newPoint = (run, val) => ({
  t: new Date(run.Date * 1000),
  y: val,
});

/**
 * Generates a new dataset for charting
 *
 * @param {string} label identifier for dataset
 * @param {number} colorIndex identifier for color
 * @param {{ t: Date, y: number }[]} data array of points
 */
const newDataset = (label, colorIndex, data) => ({
  label,
  data,
  fill: false,
  backgroundColor: getColor(colorIndex),
  borderColor: getColor(colorIndex),
  pointRadius: 4,
  pointHoverRadius: 5.5,
  lineTension: 0,
});

/**
 * Sets up run point for benchmark run
 *
 * @param {string} c benchmark type
 * @param {{ Date: number }} run benchmark run metadata
 * @param {{
 *  NsPerOp: number,
 *  BytesPerOp: number,
 *  AllocsPerOp: number,
 *  Custom: any,
 * }} bench benchmark results
 */
const newRunPoint = (c, run, bench) => {
  switch (c) {
    case 'ns/op': return newPoint(run, bench.NsPerOp);
    case 'bytes/op': return newPoint(run, bench.Mem.BytesPerOp);
    case 'allocs/op': return newPoint(run, bench.Mem.AllocsPerOp);
    default:
      if (!bench.Custom || bench.Custom[c] === undefined) {
        console.error('no value for chart type', c);
      }
      return newPoint(run, bench.Custom[c]);
  }
};

async function readJSON(path) {
  return (await fetch(path)).json();
}

function makeLabel(run) {
  const d = new Date(run.Date * 1000);
  const ds = `${d.getMonth() + 1}/${d.getDate()}/${d.getFullYear()}`;
  // if no version is available, just return the datestamp
  return run.Version ? `${run.Version.substring(0, 7)} (${ds})` : ds;
}


const chartOptions = (c, yMax) => ({
  responsive: true,
  aspectRatio: 1,
  title: {
    display: true,
    text: c,
  },
  layout: {
    padding: {
      right: 10,
    },
  },
  tooltips: {
    mode: 'index',
    intersect: false,
    position: 'nearest',
  },
  hover: {
    mode: 'index',
    intersect: false,
  },
  scales: {
    yAxes: [{
      display: true,
      ticks: { beginAtZero: true, suggestedMax: yMax },
      gridLines: {
        display: true,
        drawBorder: true,
        drawOnChartArea: false,
      },
    }],
    xAxes: [{
      display: false,
    }],
  },
  legend: {
    position: 'bottom',
    labels: {
      font: "'Open Sans', sans-serif",
      fontSize: 10,
      boxWidth: 10,
    },
  },
});

// Generate charts per suite
export async function generateCharts({
  div, // div to populate with charts
  json, // path to JSON database
  source, // source repository for package, e.g. 'github.com/bobheadxi/gobenchdata'
  canonicalImport, // import path of package, e.g. 'go.bobheadxi.dev/gobenchdata'
  chartsTypes, // additional types of charts to generate
  perBenchmark,
}) {
  chartsTypes.unshift('ns/op');
  initChartsJS();
  let runs = [];
  try {
    runs = await readJSON(json);
  } catch (e) {
    const err = document.createElement('div');
    div.appendChild(err);
    err.innerText = e;
  }

  const labels = runs.sort((a, b) => a.Date - b.Date).map((r) => makeLabel(r));
  const charts = {};
  let len = 0;
  // runs should start from the most recent run
  runs.forEach((run, runIndex) => {
    len += 1;

    // add data from each suite
    run.Suites.forEach((suite) => {
      for (let i = 0; i < suite.Benchmarks.length; i += 1) {
        const bench = suite.Benchmarks[i];

        let configCharName;
        if (perBenchmark) {
          configCharName = bench.Name;
        } else {
          configCharName = suite.Pkg;
        }

        if (charts[makeChartName(configCharName, chartsTypes[0])]) {
          // if the chart div was already set up, append data to chart.
          chartsTypes.forEach((c) => {
            const p = newRunPoint(c, run, bench);
            const chart = charts[makeChartName(configCharName, c)];

            // find appropriate dataset
            const { data: { datasets } } = chart;
            const dataset = datasets.find((e) => (e ? e.label === bench.Name : false));
            if (dataset) {
              // append new point to existing dataset
              dataset.data.push(p);
            } else {
              // generate missing points and create new dataset if it is missing
              const dataPoints = [];
              for (let runI = 0; runI < runIndex; runI += 1) {
                dataPoints.push(newPoint(runs[runI], NaN));
              }
              dataPoints.push(p);
              // append new dataset
              datasets.push(newDataset(bench.Name, datasets.length + 3, dataPoints));
            }
          });
        } else {
          // group benchmarks for a package under a div
          const group = document.createElement('div');
          group.id = configCharName;
          const title = document.createElement('h3');
          const pkgLink = document.createElement('a');

          if (perBenchmark) {
            title.innerText = configCharName;
            group.appendChild(title);
          } else {
            if (canonicalImport) {
              pkgLink.setAttribute('href', `https://${source}/tree/master/${suite.Pkg.replace(canonicalImport || source, '')}`);
            } else {
              const parts = suite.Pkg.split('/');
              source = parts.slice(0, 3).join('/');
              pkgLink.setAttribute('href', `https://${source}/tree/master/${parts.slice(3).join('/')}`);
            }
            pkgLink.setAttribute('target', '_blank');
            pkgLink.innerText = suite.Pkg;
            title.innerText = 'package\n';
            title.appendChild(pkgLink);
            group.appendChild(title);
          }

          // chart for each benchmark type
          const seedColor = randomColorID();

          const benchmarks = [bench];

          const src = source;
          chartsTypes.forEach((c) => {
            const chartName = makeChartName(configCharName, c);

            // create elements
            const canvas = document.createElement('canvas');
            canvas.id = configCharName;
            const ctx = canvas.getContext('2d');

            // create chart
            let colorI = seedColor;
            let max = 0;
            const datasets = benchmarks.map((bc) => {
              const p = newRunPoint(c, run, bc);
              max = Math.max(p.y + p.y * 0.1, max);
              colorI += 3;
              return newDataset(bc.Name, colorI, [p]);
            });

            charts[chartName] = new Chart(ctx, {
              type: 'line-with-guides',
              data: {
                labels,
                datasets,
              },
              options: chartOptions(c, max),
            });

            console.log('NEW CHART', datasets.length, charts[chartName].data.datasets, benchmarks.length);
            // TODO: this only works if you click on a point exactly, which is
            // dumb. can't seem to make it work for clicking anywhere (getting
            // the chart.js x-axis is nontrivial). ugh
            canvas.onclick = (e) => {
              const p = charts[chartName].getElementAtEvent(e);
              if (p && p.length) {
                const { _index: pointID, _xScale: x } = p[0];
                const label = x.ticks[pointID].split(' ');
                // a commit label should have 2 parts (see label())
                if (label.length === 1) {
                  return;
                }
                window.open(`https://${src}/commit/${label[0]}`, '_blank');
              }
            };

            // attach to dom
            const canvasDiv = document.createElement('div');
            canvasDiv.setAttribute('class', 'canvaswrapper');
            canvasDiv.appendChild(canvas);
            group.appendChild(canvasDiv);
          });

          // attach group to parent
          group.appendChild(document.createElement('hr'));
          group.appendChild(document.createElement('br'));
          div.appendChild(group);

          // In the event that we're adding graphs per package and not per benchmark we need to
          // break so we're not creating a new chart for each benchmark.
          if (!perBenchmark) {
            break;
          }
        }
      }
    });
    // fill missing data from datasets
    Object.values(charts).forEach((c) => {
      const { data: { datasets } } = c;
      datasets.forEach((d) => {
        const { data } = d;
        if (data.length < len) {
          data.unshift(newPoint(run, NaN));
        }
      });
    });
  });
}
