import { createApp } from "vue";

import ECharts from "vue-echarts";
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LinesChart } from "echarts/charts";
import { GridComponent, TooltipComponent } from "echarts/components";

import App from "./App.vue";

use([CanvasRenderer, LinesChart, GridComponent, TooltipComponent]);

const app = createApp(App);
app.component("v-chart", ECharts);
app.mount("#app");
