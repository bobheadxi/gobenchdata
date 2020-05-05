<template>
  <div>
    <h4>{{ config.Name }}</h4>
    <div v-if="error">{{ error }}</div>
    <div v-else v-for="c in generateCharts()" :key="c.metric">
      <h5>{{ c.metric }}</h5>
      <p>{{ c.options }}</p>
      <apexchart :options="c.options" :series="c.options.series"></apexchart>
    </div>
  </div>
</template>

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

        return Object.keys(seriesByMetric).map((m) => ({
          metric: m,
          options: {
            chart: {
              type: 'line',
              height: 200,
            },
            dataLabels: {
              enabled: false,
            },
            xaxis: {},
            series: seriesByMetric[m],
          },
        }));
      } catch (err) {
        console.error(err);
        this.error = err;
        return [];
      }
    },
  },
});
</script>

<style scoped lang="scss">

</style>
