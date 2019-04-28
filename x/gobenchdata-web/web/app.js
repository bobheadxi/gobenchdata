// Generate one chart per suite
export async function generateCharts({
  div,  // div to populate with charts 
  json, // path to JSON database
}) {
  const runs = await readJSON(json);
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
            labels: runs.map(run => label(run)),
            datasets: suite.Benchmarks.map(bench => {
              return {
                label: bench.Name,
                data: [newPoint(run, bench)]
              }
            }),
          },
          options: {
            tooltips: {
              callbacks: {
                label: function(tt, data) {
                  let label = data.datasets[tt.datasetIndex].label || '';
                  if (label) {
                    label += ': ';
                  }
                  label += tt.yLabel;
                  return label;
                }
              },
            }
          }
        });
      }
    });
  })
}

async function readJSON(path) { return (await fetch(path)).json(); }

function newPoint(run, bench) {
  return {
    x: run.Version,
    t: new Date(run.Date*1000),
    y: bench.NsPerOp,
  }
}

function label(run) {
  const d = new Date(run.Date*1000)
  return `${run.Version.substring(0, 7)} (${d.toDateString()}, ${d.toLocaleTimeString()})`
}
