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
  runs.forEach(run => {
    run.Suites.forEach(suite => {
      if (charts[suite.Pkg]) {
        // if the chart div was already set up, append data to chart
        suite.Benchmarks.forEach(bench => {
          charts[suite.Pkg].data.datasets
            .find(e => (e.label === bench.Name)).data
            .push(newPoint(run, bench));
        });
      } else {
        const canvas = document.createElement('canvas');
        div.appendChild(canvas);
        canvas.id = suite.Pkg;
        const ctx = canvas.getContext('2d');
        charts[suite.Pkg] = new Chart(ctx, {
          type: 'line',
          data: {
            labels: runs.sort((a, b) => a.Date - b.Date).map(run => label(run)),
            datasets: suite.Benchmarks.map(bench => {
              return {
                label: bench.Name,
                data: [newPoint(run, bench)],
              }
            }),
          },
          options: { },
        });
      }
    });
  })
}

async function readJSON(path) {
  return (await fetch(path)).json();
}

function newPoint(run, bench) {
  return {
    t: new Date(run.Date*1000),
    y: bench.NsPerOp,
  }
}

function label(run) {
  const d = new Date(run.Date*1000)
  return `${run.Version.substring(0, 7)} (${d.toDateString()}, ${d.toLocaleTimeString()})`
}
