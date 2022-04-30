<template>
  <div class="app">
    <div v-if="loading">loading...</div>
    <div v-if="error">
      {{ error }}
    </div>

    <!-- okay state -->
    <div v-if="!loading && !error">
      <h1>{{ config.title }}</h1>
      <h3 v-html="config.description"></h3>

      <div v-for="g in groups" :key="g.name">
        <ChartGroup :group="g" :runs="benchmarks" :repo="config.repository" />
      </div>
    </div>

    <div class="footer">
      <p>
        generated using
        <a href="https://bobheadxi.dev/r/gobenchdata">gobenchdata</a>
      </p>
    </div>
  </div>
</template>

<style lang="scss">
.app {
  padding-left: $gap;
  padding-right: $gap;

  font-family: "Fira Code", monospace;
  text-align: center;

  .footer {
    margin-top: 2 * $gap;
    margin-bottom: 2 * $gap;
  }
}
</style>

<script lang="ts">
import { defineComponent } from "vue";
import { parse } from "yaml";
import ChartGroup from "@/components/ChartGroup.vue";

import { iterateSuites } from "@/lib/series";
import {
  Config,
  Run,
  ChartGroup as ConfigChartGroup,
  Chart,
} from "@/generated";

type AppState = {
  loading: boolean;
  config: Config;
  benchmarks: Run[];
  error: unknown;
};

export default defineComponent({
  name: "App",
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
    groups(): ConfigChartGroup[] {
      if (this.config.chartGroups && this.config.chartGroups.length > 0)
        return this.config.chartGroups;

      // group by package by default
      const suites: { [pkg: string]: boolean } = {};
      iterateSuites(this.benchmarks, (s) => {
        suites[s.Pkg] = true;
      });
      const packages = Object.keys(suites).sort();
      console.log("no chart groups configured - detected packages", packages);
      return [
        new ConfigChartGroup({
          name: "Benchmarks by package",
          description: "All detected benchmarks, grouped by Package",
          charts: packages.map(
            (pkg) =>
              new Chart({
                name: pkg,
                package: `^${pkg}$`,
                benchmarks: [],
              })
          ),
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
        // load config
        const configResp = await fetch("./gobenchdata-web.yml");
        if (configResp.status > 400) {
          console.error(configResp);
          throw new Error(`${configResp.status}: failed to load config`);
        }
        const raw = await configResp.text();
        console.log(raw);
        const config = new Config(parse(raw));

        // load benchmark runs
        const benchmarksResp = await fetch(
          `./${config.benchmarksFile || "benchmarks.json"}`
        );
        if (benchmarksResp.status > 400) {
          console.error(benchmarksResp);
          throw new Error(
            `${benchmarksResp.status}: failed to load benchmarks`
          );
        }
        const runs = await benchmarksResp.json();

        // update state
        this.benchmarks = runs.map((r: unknown) => new Run(r));
        this.config = config;
        console.log("data loaded", { config, runs: this.benchmarks.length });
      } catch (err) {
        this.error = err;
      }
      this.loading = false;
    },
  },
});
</script>
