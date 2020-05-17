import Vue from 'vue';
import VueApexCharts from 'vue-apexcharts';
import App from './App.vue';

Vue.config.productionTip = false;
Vue.component('apexchart', VueApexCharts);

new Vue({
  render: h => h(App),
}).$mount('#app');
