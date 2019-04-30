'use strict';

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
  const d = new Date(run.Date*1000)
  return `${run.Version.substring(0, 7)} (${d.getMonth()}/${d.getDay()}/${d.getFullYear()})`
}

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

  const charts = {};
  let len = 0;
  // runs should start from the most recent run
  // TODO account for missing data towards end of data and in the middle
  runs.forEach(run => {
    run.Suites.forEach(suite => {
      if (charts[suite.Pkg]) {
        // if the chart div was already set up, append data to chart.
        // if the dataset is isn't in the datasets, then it no longer exists,
        // and we'll ignore it.
        suite.Benchmarks.forEach(bench => {
          const chart = charts[suite.Pkg];
          const dataset = chart.data.datasets.find(e => (e.label === bench.Name))
          if (dataset) {
            dataset.data.push(newPoint(run, bench));
          }
        });
      } else {
        // create elements
        const canvasDiv = document.createElement('div');
        const canvas = document.createElement('canvas');
        canvasDiv.appendChild(canvas);
        div.appendChild(canvasDiv);
        canvas.id = suite.Pkg;
        const ctx = canvas.getContext('2d');

        // create chart
        let i = randomInt();
        charts[suite.Pkg] = new Chart(ctx, {
          type: 'line',
          data: {
            labels: runs.sort((a, b) => a.Date - b.Date).map(run => label(run)),
            datasets: suite.Benchmarks.map(bench => {
              i += 3;
              return {
                label: bench.Name,
                data: [newPoint(run, bench)],

                fill: false,
                backgroundColor: getColor(i),
                borderColor: getColor(i),
                pointRadius: 4,
              }
            }),
          },
          options: chartOptions(suite),
        });
      }
    });
    run++;
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

const newPoint = (run, bench) => ({
  t: new Date(run.Date*1000),
  y: bench.NsPerOp,
})
