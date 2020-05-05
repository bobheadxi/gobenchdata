<template>
  <div class="app">
    <div v-if="loading">
      loading...
    </div>
    <div v-if="error">
      {{ error }}
    </div>

    <!-- okay state -->
    <div v-else>
      <h1>{{ config.Title }}</h1>
      <h2>{{ config.Description }}</h2>

      <div v-for="g in chartGroups" :key="g.name">
        <ChartGroup :group="g" :runs="benchmarks" />
      </div>
    </div>

    <hr />

    <div>
      <p>generated with <a href="https://bobheadxi.dev/gobenchdata">gobenchdata</a></p>
    </div>
  </div>
</template>

<style lang="scss">
.app {
  padding-left: $gap;
  padding-right: $gap;

  font-family: 'Fira Code', monospace;
  text-align: center;

  hr {
    margin-top: 2 * $gap;
    margin-bottom: 2 * $gap;
  }
}
</style>

<script lang="ts">
import Vue from 'vue';
import ChartGroup from '@/components/ChartGroup.vue';

import { Config, Run, ConfigChartGroup, ConfigChartGroupChart } from '@/generated';

type AppState = {
  loading: boolean;
  config: Config;
  benchmarks: Run[];
  error: any;
}

export default Vue.extend({
  name: 'App',
  components: {
    ChartGroup,
  },
  data: (): AppState => ({
    loading: true,
    config: new Config(),
    benchmarks: [],
    error: undefined,
  }),
  computed: {
    chartGroups(): ConfigChartGroup[] {
      if (this.config.ChartGroups && this.config.ChartGroups.length > 0) return this.config.ChartGroups;
      // TODO: generate chart groups from benchmark runs
      return [
        new ConfigChartGroup({
          Name: 'test group',
          Description: 'a test group',
          Charts: [
            new ConfigChartGroupChart({
              Name: 'test chart',
              Description: 'a test chart',
              Package: '**',
              Benchmarks: ['**'],
            }),
          ],
        }),
      ];
    },
  },
  created() {
    this.load();
  },
  methods: {
    async load() {
      try {
        const configResp = await fetch('./gobenchdata-web.json');
        if (configResp.status > 400) {
          console.error(configResp);
          throw new Error(`${configResp.status}: failed to load config`);
        }

        const config = new Config(await configResp.json());
        const benchmarksResp = await fetch(`./${config.BenchmarksFile || 'benchmarks.json'}`);
        if (benchmarksResp.status > 400) {
          console.error(benchmarksResp);
          throw new Error(`${benchmarksResp.status}: failed to load benchmarks`);
        }

        const runs = await benchmarksResp.json();
        this.benchmarks = runs.map((r: any) => new Run(r));
        this.config = config;
        console.log('config loaded', { config });
      } catch (err) {
        this.error = err;
      }
      this.loading = false;
    },
  },
});
</script>
