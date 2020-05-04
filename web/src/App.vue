<template>
  <div>
    <h1>Hello World</h1>
    <p v-if="loading">loading...</p>
    <P v-if="error"> {{ error }} </p>
    <p v-else> {{ config }} {{ benchmarks }} </p>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { Config, Run } from './types/generated';

type AppState = {
  loading: boolean;
  config: Config;
  benchmarks: Run[];
  error: any;
}

export default Vue.extend({
  name: 'App',
  data: (): AppState => ({
    loading: true,
    config: new Config(),
    benchmarks: [],
    error: undefined,
  }),
  created() {
    this.load();
  },
  methods: {
    async load() {
      try {
        const configResp = await fetch('/gobenchdata-web-config.json');
        if (configResp.status > 400) throw new Error(`${configResp.status}: failed to load config`);

        const config = new Config(await configResp.json());
        const benchmarksResp = await fetch(`/${config.BenchmarksFile || 'benchmarks.json'}`);
        if (benchmarksResp.status > 400) throw new Error(`${benchmarksResp.status}: failed to load benchmarks`);

        const runs = await benchmarksResp.json();
        this.benchmarks = runs.map((r: any) => new Run(r));
        this.config = config;
      } catch (err) {
        this.error = err;
      }
      this.loading = false;
    },
  },
});
</script>

<style scoped lang="scss">

</style>
