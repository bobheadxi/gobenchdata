'use strict';

function initChartsJS() {
  Chart.defaults['line-with-guides'] = Chart.defaults.line;
  Chart.controllers['line-with-guides'] = Chart.controllers.line.extend({
    draw: function(ease) {
      Chart.controllers.line.prototype.draw.call(this, ease);

      if (this.chart.tooltip._active && this.chart.tooltip._active.length) {
        let activePoint = this.chart.tooltip._active[0],
            ctx = this.chart.ctx,
            x = activePoint.tooltipPosition().x,
            topY = this.chart.scales['y-axis-0'].top,
            bottomY = this.chart.scales['y-axis-0'].bottom;

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
    }
  });
}

// Generate charts per suite
export async function generateCharts({
  div,             // div to populate with charts 
  json,            // path to JSON database
  source,          // source repository for package, e.g. 'github.com/bobheadxi/gobenchdata'
  canonicalImport, // import path of package, e.g. 'go.bobheadxi.dev/gobenchdata'
  chartsTypes,     // additional types of charts to generate
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

  const labels = runs.sort((a, b) => a.Date - b.Date).map(r => label(r));
  const charts = {};
  let len = 0;
  // runs should start from the most recent run
  runs.forEach((run, runIndex) => {
    len++;

    // add data from each suite
    run.Suites.forEach(suite => {
      if (charts[makeChartName(suite.Pkg, chartsTypes[0])]) {
        // if the chart div was already set up, append data to chart.
        suite.Benchmarks.forEach(bench => {
          chartsTypes.forEach(c => {
            const p = newRunPoint(c, run, bench);
            const chart = charts[makeChartName(suite.Pkg, c)];

            // find appropriate dataset
            const { data: { datasets } } = chart;
            const dataset = datasets.find(e => (e ? e.label === bench.Name : false));
            if (dataset) {
              // append new point to existing dataset
              dataset.data.push(p);
            } else {
              // generate missing points and create new dataset if it is missing
              const dataPoints = [];
              for (let runI = 0; runI < runIndex; runI++) dataPoints.push(newPoint(runs[runI], NaN));
              dataPoints.push(p)
              // append new dataset
              datasets.push(newDataset(bench.Name, datasets.length + 3, dataPoints));
            }
          })
        });
      } else {
        // group benchmarks for a package under a div
        const group = document.createElement('div');
        group.id = suite.Pkg;
        const title = document.createElement('h3');
        const pkgLink = document.createElement('a');
        if (canonicalImport) {
          pkgLink.setAttribute('href', `https://${source}/tree/master/${suite.Pkg.replace(canonicalImport || source, '')}`);
        } else {
          const parts = suite.Pkg.split('/');
          source = parts.slice(0, 3).join('/');
          pkgLink.setAttribute('href', `https://${source}/tree/master/${parts.slice(3).join('/')}`);
        }
        pkgLink.setAttribute('target', '_blank');
        pkgLink.innerText = suite.Pkg;
        title.innerText = `package\n`;
        title.appendChild(pkgLink);
        group.appendChild(title);

        // chart for each benchmark type
        let seedColor = randomInt();
        const { Benchmarks: benchmarks } = suite;
        chartsTypes.forEach(c => {
          const chartName = makeChartName(suite.Pkg, c);

          // create elements
          const canvas = document.createElement('canvas');
          canvas.id = chartName;
          const ctx = canvas.getContext('2d');

          // create chart
          let i = seedColor;
          let max = 0;
          const datasets = benchmarks.map(bench => {
            const p = newRunPoint(c, run, bench);
            max = Math.max(p.y+p.y*0.1, max);
            i += 3;
            return newDataset(bench.Name, i, [p]);
          })
          charts[chartName] = new Chart(ctx, {
            type: 'line-with-guides',
            data: {
              labels,
              datasets,
            },
            options: chartOptions(c, max),
          });
          console.log('NEW CHART', datasets.length, charts[chartName].data.datasets, benchmarks.length)
          // TODO: this only works if you click on a point exactly, which is
          // dumb. can't seem to make it work for clicking anywhere (getting
          // the chart.js x-axis is nontrivial). ugh
          canvas.onclick = (e) => {
            const p = charts[chartName].getElementAtEvent(e);
            if (p && p.length) {
              const { _index: i, _xScale: x } = p[0];
              const label = x.ticks[i].split(' ');
              // a commit label should have 2 parts (see label())
              if (label.length === 1) return;
              window.open(`https://${source}/commit/${label[0]}`, '_blank');
            }
          }

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
      }
    });

    // fill missing data from datasets
    Object.values(charts).forEach(c => {
      const { data: { datasets } } = c;
      datasets.forEach(d => {
        const { data } = d;
        if (data.length < len) data.unshift(newPoint(run, NaN));
      });
    })
  })
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
    }]
  },
  legend: {
    position: 'bottom',
    labels: {
      font: "'Open Sans', sans-serif",
      fontSize: 10,
      boxWidth: 10,
    },
  },
})

const makeChartName = (pkg, chart) => `${pkg}-${chart}`;

const newPoint = (run, val) => ({
  t: new Date(run.Date*1000),
  y: val,
})

const newDataset = (label, colorIndex, data) => ({
  label,
  data,
  fill: false,
  backgroundColor: getColor(colorIndex),
  borderColor: getColor(colorIndex),
  pointRadius: 4,
  pointHoverRadius: 5.5,
  lineTension: 0,
})

const newRunPoint = (c, run, bench) => {
  switch (c) {
    case 'ns/op': return newPoint(run, bench.NsPerOp);
    case 'bytes/op': return newPoint(run, bench.Mem.BytesPerOp);
    case 'allocs/op': return newPoint(run, bench.Mem.AllocsPerOp);
    default:
      if (!bench.custom || bench.custom[c] === undefined) {
        console.error('no value for chart type', c)
      }
      return newPoint(run, bench.Cusom[c]);
  }
}

const chartColors = {
	red: 'rgb(255, 99, 132)',
	orange: 'rgb(255, 159, 64)',
	yellow: 'rgb(255, 205, 86)',
	green: 'rgb(75, 192, 192)',
	blue: 'rgb(54, 162, 235)',
	purple: 'rgb(153, 102, 255)',
	grey: 'rgb(201, 203, 207)'
};

const chartColorsList = Object.values(chartColors);

function randomInt() {
  return Math.floor(Math.random() * (chartColorsList.length + 1));
}

function getColor(i) {
  return chartColorsList[i % chartColorsList.length];
}

async function readJSON(path) {
  return (await fetch(path)).json();
}

function label(run) {
  const d = new Date(run.Date*1000);
  let month = d.getMonth();
  const ds = `${++month}/${d.getDate()}/${d.getFullYear()}`
  // if no version is available, just return the datestamp
  return run.Version ? `${run.Version.substring(0, 7)} (${ds})` : ds;
}
