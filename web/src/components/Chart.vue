<template>
  <div class="chart">
    <h3>{{ config.Name }}</h3>
    <p>{{ description }}</p>
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

      .chart-container {
        height: 520px;
      }
    }
  }
}
</style>

<script lang="ts">
import Vue, { PropType } from 'vue';
import { ApexOptions } from 'apexcharts';
import { Run, ConfigChartGroupChart } from '@/generated';
import { generateSeries } from '@/lib/series';

export default Vue.extend({
  name: 'Chart',
  props: {
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
      return this.config.Description
        || `"Package": "${this.config.Package}", "Benchmarks": ${JSON.stringify(this.config.Benchmarks)}`;
    },
  },
  methods: {
    generateCharts(): { metric: string; options: ApexOptions }[] {
      try {
        const pkgMatcher = new RegExp(this.config.Package || '.');
        const benchMatchers = this.config.Benchmarks.map((b) => new RegExp(b || '.'));
        if (benchMatchers.length === 0) benchMatchers.push(new RegExp('.'));
        const seriesByMetric = generateSeries(this.runs, pkgMatcher, benchMatchers, this.config.Metrics);

        const generatedCharts = Object.keys(seriesByMetric).map((m): { metric: string; options: ApexOptions } => ({
          metric: m,
          options: {
            chart: {
              type: 'line',
              height: 500,
            },
            markers: {
            },
            dataLabels: {
              enabled: false,
            },
            xaxis: {},
            series: seriesByMetric[m],
          },
        }));
        console.log(`chart ${this.config.Name}`, generatedCharts);
        return generatedCharts || [];
      } catch (err) {
        console.error(`chart ${this.config.Name}`, this.config, err);
        this.error = err;
        return [];
      }
    },
  },
});
</script>
