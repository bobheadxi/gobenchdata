<template>
  <div class="chart">
    <h4>{{ config.Name }}</h4>
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
import { Minimatch } from 'minimatch';
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
  methods: {
    generateCharts(): { metric: string; options: ApexOptions }[] {
      try {
        // use https://github.com/isaacs/minimatch syntax
        const pkgMatcher = new Minimatch(this.config.Package).makeRe();
        const benchMatchers = this.config.Benchmarks.map((b) => new Minimatch(b).makeRe());
        const seriesByMetric = generateSeries(this.runs, pkgMatcher, benchMatchers, this.config.Metrics);

        const results = Object.keys(seriesByMetric).map((m): { metric: string; options: ApexOptions } => ({
          metric: m,
          options: {
            chart: {
              type: 'line',
              height: 400,
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
        console.log(results);
        return results;
      } catch (err) {
        this.error = err;
        return [];
      }
    },
  },
});
</script>
