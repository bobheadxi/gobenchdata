'use strict';

// Generate charts per suite
export async function generateCharts({
  div,        // div to populate with charts 
  json,       // path to JSON database
  rootImport, // import path of package, e.g. 'github.com/bobheadxi/gobenchdata'
}) {
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
  runs.forEach(run => {
    len++;

    // add data from each suite
    run.Suites.forEach(suite => {
      if (charts[suite.Pkg+'-'+chartsTypes[0]]) {
        // if the chart div was already set up, append data to chart.
        // if the dataset is isn't in the datasets, then it no longer exists,
        // and we'll ignore it.
        suite.Benchmarks.forEach(bench => {
          chartsTypes.forEach(c => {
            const { data: { datasets } } = charts[suite.Pkg + '-' + c];
            const dataset = datasets.find(e => (e.label === bench.Name))
            if (dataset) dataset.data.push(newRunPoint(c, run, bench));
          })
        });
      } else {
        // group benchmarks for a package under a div
        const group = document.createElement('div');
        group.id = suite.Pkg;
        const title = document.createElement('h3');
        const pkgLink = document.createElement('a');
        if (rootImport) {
          pkgLink.setAttribute('href', `https://${rootImport}/tree/master/${suite.Pkg.replace(rootImport, '')}`);
        } else {
          const parts = suite.Pkg.split('/');
          rootImport = parts.slice(0, 3).join('/');
          pkgLink.setAttribute('href', `https://${rootImport}/tree/master/${parts.slice(3).join('/')}`);
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
          const chartName = suite.Pkg + '-' + c;

          // create elements
          const canvas = document.createElement('canvas');
          canvas.id = chartName;
          const ctx = canvas.getContext('2d');

          // create chart
          let i = seedColor;
          let max = 0;
          charts[chartName] = new Chart(ctx, {
            type: 'line',
            data: {
              labels,
              datasets: benchmarks.map(bench => {
                const p = newRunPoint(c, run, bench);
                max = Math.max(p.y+p.y*0.1, max);
                i += 3;
                return {
                  label: bench.Name,
                  data: [p],

                  fill: false,
                  backgroundColor: getColor(i),
                  borderColor: getColor(i),
                  pointRadius: 4,
                  pointHoverRadius: 5.5,
                  lineTension: 0,
                }
              }),
            },
            options: chartOptions(c, max),
          });
          // TODO: this only works if you click on a point exactly, which is
          // dumb. can't seem to make it work for clicking anywhere (getting
          // the chart.js x-axis is nontrivial). ugh
          canvas.onclick = (e) => {
            const p = charts[chartName].getElementAtEvent(e);
            if (p && p.length) {
              const { _index: i, _xScale: x } = p[0];
              const commit = x.ticks[i].split(' ')[0];
              window.open(`https://${rootImport}/commit/${commit}`, '_blank');
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
      fontSize: 10,
      boxWidth: 10,
    },
  },
})

const newPoint = (run, val) => ({
  t: new Date(run.Date*1000),
  y: val,
})

const newRunPoint = (c, run, bench) => {
  switch (c) {
    case chartsTypes[0]: return newPoint(run, bench.NsPerOp);
    case chartsTypes[1]: return newPoint(run, bench.Mem.BytesPerOp);
    case chartsTypes[2]: return newPoint(run, bench.Mem.AllocsPerOp);
    default: console.error('unexpected chart type', c);
  }
}

const chartsTypes = [
  'ns/op',
  'bytes/op',
  'allocs/op',
]

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
  return `${run.Version.substring(0, 7)} (${++month}/${d.getDate()}/${d.getFullYear()})`;
}
