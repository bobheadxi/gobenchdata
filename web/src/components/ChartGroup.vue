<template>
  <div class="chart-group">
    <h3>{{ group.Name }}</h3>
    <p>{{ group.Description }}</p>

    <div v-for="c in group.Charts" :key="c.name">
      <Chart :config="c" :runs="runs" />
    </div>
  </div>
</template>

<style lang="scss">
.chart-group {
  margin-top: 2 * $gap;
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
