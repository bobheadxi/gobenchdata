<template>
  <div class="chart-group">
    <h2>{{ group.name }}</h2>
    <p>{{ group.description }}</p>

    <div v-for="c in group.charts" :key="c.name" class="chart">
      <Chart :config="c" :runs="runs" :repo="repo" />
    </div>
  </div>
</template>

<style lang="scss">
.chart-group {
  margin-top: 2 * $gap;

  .chart {
    margin-top: $gap;
    margin-bottom: $gap;
  }
}
</style>

<script lang="ts">
import Vue from 'vue';
import Chart from '@/components/Chart.vue';
import { Run, ConfigChartGroup } from '@/generated';

export default Vue.extend({
  name: 'ChartGroup',
  components: {
    Chart,
  },
  props: {
    repo: { type: String, required: true },
    group: {
      type: Object as () => ConfigChartGroup,
      required: true,
    },
    runs: {
      type: Array as () => Run[],
      required: true,
    },
  },
});
</script>
