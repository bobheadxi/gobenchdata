'use strict';

// Generate one chart per suite
export async function generateCharts({
  div,  // div to populate with charts 
  json, // path to JSON database
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
      if (charts[suite.Pkg]) {
        // if the chart div was already set up, append data to chart.
        // if the dataset is isn't in the datasets, then it no longer exists,
        // and we'll ignore it.
        suite.Benchmarks.forEach(bench => {
          const { data: { datasets } } = charts[suite.Pkg];
          const dataset = datasets.find(e => (e.label === bench.Name))
          if (dataset) dataset.data.push(newPoint(run, bench.NsPerOp));
        });
      } else {
        // create elements
        const canvas = document.createElement('canvas');
        canvas.id = suite.Pkg;
        const ctx = canvas.getContext('2d');

        // create chart
        let i = randomInt();
        const { Benchmarks: benchmarks } = suite;
        charts[suite.Pkg] = new Chart(ctx, {
          type: 'line',
          data: {
            labels,
            datasets: benchmarks.map(bench => {
              i += 3;
              return {
                label: bench.Name,
                data: [newPoint(run, bench.NsPerOp)],

                fill: false,
                backgroundColor: getColor(i),
                borderColor: getColor(i),
                pointRadius: 4,
              }
            }),
          },
          options: chartOptions(suite),
        });

        // attach to dom
        const canvasDiv = document.createElement('div');
        canvasDiv.appendChild(canvas);
        div.appendChild(canvasDiv);
        div.appendChild(document.createElement('br'));
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

const chartOptions = (suite) => ({
  responsive: true,
  title: {
    display: true,
    text: suite.Pkg,
  },
  tooltips: {
    mode: 'index',
    intersect: false,
  },
  hover: {
    mode: 'nearest',
    intersect: true
  },
})

const newPoint = (run, val) => ({
  t: new Date(run.Date*1000),
  y: val,
})


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
, val) => ({
  t: new Date(run.Date*1000),
  y: val,
})
