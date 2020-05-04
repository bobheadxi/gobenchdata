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

type AppState = {
  loading: boolean;
  config: {};
  benchmarks: {}[];
  error: any;
}

export default Vue.extend({
  name: 'App',
  data: (): AppState => ({
    loading: true,
    config: {},
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

        const config = await configResp.json();
        const benchmarksResp = await fetch(`/${config.benchmarksFile || 'benchmarks.json'}`);
        if (benchmarksResp.status > 400) throw new Error(`${benchmarksResp.status}: failed to load benchmarks`);

        this.benchmarks = await benchmarksResp.json();
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
