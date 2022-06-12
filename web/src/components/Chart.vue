<template>
  <div class="chart">
    <h3>{{ config.name }}</h3>
    <p v-html="description"></p>
    <div v-if="error">{{ error }}</div>
    <div v-else class="chart-set">
      <div
        v-for="c in generateCharts()"
        class="metric"
        :key="c.metric"
        :class="{ 'full-width': config.display && config.display.fullWidth }"
      >
        <h5>{{ c.metric }}</h5>
        <div class="chart-container">
          <v-chart
            class="chart"
            :option="c.options"
            :autoresize="true"
            :height="config.display && config.display.fullWidth ? 300 : 'auto'"
          />
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
        width: 50%;
      }
      @media (max-width: $touch) {
        width: 100%;
      }

      &.full-width {
        width: 80%;
        @media (max-width: $touch) {
          width: 100%;
        }
      }
    }
  }
}
</style>

<script lang="ts">
import { defineComponent, PropType, ref, Ref } from "vue";
import type { ECBasicOption } from "echarts/types/dist/shared";
import VChart, { THEME_KEY } from "vue-echarts";

import { ParseDate } from "../generated";
import type { Run, Chart as ConfigChart } from "../generated";
import { generateSeries } from "../lib/series";

type ChartOption = { metric: string; options: ECBasicOption };

type LineSeries = {
  name: string;
  coords: number[][];
};

export default defineComponent({
  name: "LineChart",
  props: {
    repo: { type: String, required: true },
    config: {
      type: Object as PropType<ConfigChart>,
      required: true,
    },
    runs: {
      type: Array as () => Run[],
      required: true,
    },
  },
  components: {
    VChart,
  },
  provide: {
    [THEME_KEY]: "dark",
  },
  data: () => ({ error: undefined }),
  computed: {
    description(): string {
      return (
        this.config.description ||
        `"Package": "${this.config.package}", "Benchmarks": ${JSON.stringify(
          this.config.benchmarks
        )}`
      );
    },
  },
  methods: {
    generateCharts(): ChartOption[] {
      try {
        const pkgMatcher = new RegExp(this.config.package || ".");
        const benchMatchers =
          this.config.benchmarks?.map((b) => new RegExp(b || ".")) || [];
        if (benchMatchers.length === 0) benchMatchers.push(new RegExp("."));
        const series = generateSeries(
          this.runs,
          pkgMatcher,
          benchMatchers,
          this.config.metrics
        );
        console.log(`chart ${this.config.name}`, series);

        return Object.keys(series.charts).map(
          (m): ChartOption => ({
            metric: m,
            options: {
              series: series.charts[m].map((s) => {
                return {
                  type: "line",
                  name: s.name || "",
                  data: s.data,
                };
              }),
            },
          })
        );
      } catch (err) {
        console.error(`chart ${this.config.name}`, this.config, err);
        return [];
      }
    },
  },
});
</script>

<!-- {
                chart: {
                  type: "line",
                  events: {
                    markerClick: (event, chartContext, config) => {
                      if (!this.repo) return;
                      const { dataPointIndex: x } = config;
                      const d = ParseDate(series.xaxis[m][x]);
                      const r = this.runs.find((r) => {
                        return ParseDate(r.Date).valueOf() === d.valueOf();
                      });
                      if (r) window.open(`${this.repo}/commit/${r.Version}`);
                    },
                  },
                  toolbar: {
                    tools: {
                      zoom: false,
                      selection: false,
                    },
                  },
                },
                responsive: [
                  {
                    breakpoint: 769, // "tablet"
                    options: {
                      chart: {
                        height: 300,
                      },
                    },
                  },
                ],
                dataLabels: {
                  enabled: false,
                },
                xaxis: {
                  type: "category",
                  categories: series.xaxis[m],
                  sorted: true,
                  tooltip: { enabled: false },
                  labels: {
                    show: false,
                    formatter: (date): string => {
                      const d = ParseDate(date);
                      const r = this.runs.find((r) => {
                        return ParseDate(r.Date).valueOf() === d.valueOf();
                      });
                      // Tue May 05 2020 hh:mm
                      const formatted = `${d.toDateString()} ${(
                        "0" + d.getHours()
                      ).slice(-2)}:${("0" + d.getMinutes()).slice(-2)}`;
                      return r && r.Version
                        ? `${formatted} (${r.Version.substring(0, 9)})`
                        : formatted;
                    },
                  },
                },
                tooltip: {
                  enabled: true,
                  shared: true,
                  onDatasetHover: {
                    highlightDataSeries: true,
                  },
                  fixed: {
                    enabled: true,
                    position: "topLeft",
                  },
                },
                // flatten since we are using categories
                series: series.charts[m].map((s) => {
                  // TODO remove eslint disable, I forget why we do this
                  s.data = s.data.map((p) => p.y) as any; // eslint-disable-line @typescript-eslint/no-explicit-any
                  return s;
                }),
              } -->
