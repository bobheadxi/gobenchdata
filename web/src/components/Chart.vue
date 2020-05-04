<template>
  <div>
    <h4>{{ chart.Name }}</h4>

    {{ chart.Benchmarks }}
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

import { Run, ConfigChartGroupChart } from '../generated';

export default Vue.extend({
  name: 'Chart',
  props: {
    chart: {
      type: Object as () => ConfigChartGroupChart,
      required: true,
    },
    runs: {
      type: Array as () => Run[],
      required: true,
    },
  },
});
</script>

<style scoped lang="scss">

</style>
