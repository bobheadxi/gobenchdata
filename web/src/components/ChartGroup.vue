<template>
  <div class="chart-group">
    <h2>{{ group.Name }}</h2>
    <p>{{ group.Description }}</p>

    <div v-for="c in group.Charts" :key="c.name" class="chart">
      <Chart :config="c" :runs="runs" />
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
