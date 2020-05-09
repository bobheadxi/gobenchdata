<template>
  <div class="chart">
    <h3>{{ config.name }}</h3>
    <p v-html="description"></p>
    <div v-if="error">{{ error }}</div>
    <div v-else class="chart-set">
      <div v-for="c in generateCharts()" :key="c.metric" class="metric">
        <h5>{{ c.metric }}</h5>
        <div class="chart-container">
          <apexchart
            :options="c.options"
            :series="c.options.series"
          ></apexchart>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.chart {
  .chart-set {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;

    .metric {
      width: 33%;
      @media (max-width: $desktop) {
        width: 50%
      }
      @media (max-width: $touch) {
        width: 100%
      }
    }
  }
}
</style>

<script lang="ts">
import Vue, { PropType } from 'vue';
import { ApexOptions } from 'apexcharts';
import { Run, ConfigChartGroupChart, ParseDate } from '@/generated';
import { generateSeries } from '@/lib/series';

export default Vue.extend({
  name: 'Chart',
  props: {
    repo: { type: String, required: true },
    config: {
      type: Object as PropType<ConfigChartGroupChart>,
      required: true,
    },
    runs: {
      type: Array as () => Run[],
      required: true,
    },
  },
  data: () => ({ error: undefined }),
  computed: {
    description(): string {
      return this.config.description
        || `"Package": "${this.config.package}", "Benchmarks": ${JSON.stringify(this.config.benchmarks)}`;
    },
  },
  methods: {
    generateCharts(): { metric: string; options: ApexOptions }[] {
      try {
        const pkgMatcher = new RegExp(this.config.package || '.');
        const benchMatchers = this.config.benchmarks.map((b) => new RegExp(b || '.'));
        if (benchMatchers.length === 0) benchMatchers.push(new RegExp('.'));
        const series = generateSeries(this.runs, pkgMatcher, benchMatchers, this.config.metrics);
        console.log(`chart ${this.config.name}`, series);

        return Object.keys(series.charts).map((m): { metric: string; options: ApexOptions } => ({
          metric: m,
          options: {
            chart: {
              type: 'line',
              height: 700,
              events: {
                click: (event, chartContext, config) => {
                  const { dataPointIndex: x } = config;
                  const d = ParseDate(series.xaxis[m][x]);
                  const r = this.runs.find(r => {
                    return ParseDate(r.Date).valueOf() === d.valueOf();
                  });
                  if (r) window.open(`${this.repo}/commit/${r.Version}`);
                },
              },
              toolbar: {
                tools: {
                  zoom: false,
                  selection: false,
                },
              },
            },
            responsive: [{
              breakpoint: 769, // "tablet"
              options: {
                chart: {
                  height: 300,
                },
              },
            }],
            dataLabels: {
              enabled: false,
            },
            xaxis: {
              type: 'category',
              categories: series.xaxis[m],
              sorted: true,
              tooltip: { enabled: false },
              labels: {
                show: false,
                formatter: (date): string => {
                  const d = ParseDate(date);
                  const r = this.runs.find(r => {
                    return ParseDate(r.Date).valueOf() === d.valueOf();
                  });
                  // Tue May 05 2020 hh:mm
                  const formatted = `${d.toDateString()} ${('0' + d.getHours()).slice(-2)}:${('0' + d.getMinutes()).slice(-2)}`;
                  return r && r.Version ? `${formatted} (${r.Version.substring(0, 9)})` : formatted;
                },
              },
            },
            tooltip: {
              enabled: true,
              shared: true,
              onDatasetHover: {
                highlightDataSeries: true,
              },
              fixed: {
                enabled: true,
                position: 'topLeft',
              },
            },
            // flatten since we are using categories
            series: series.charts[m].map((s) => {
              s.data = s.data.map(p => p.y) as any;
              return s;
            }),
          },
        })) || [];
      } catch (err) {
        console.error(`chart ${this.config.name}`, this.config, err);
        this.error = err;
        return [];
      }
    },
  },
});
</script>
