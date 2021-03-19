import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import moment from 'moment'
import App from './app.vue'
import router from './router'
import i18n from './lang'
import util from './lib/util.js'
import data from './lib/data.js'
import "xterm/dist/xterm.css";
import "xterm/dist/xterm.js";
import { Terminal } from 'xterm';
import './scss/app.scss'

let localeLang
if (global.navigator.language) {
    localeLang = global.navigator.language
    localeLang = localeLang.toLowerCase()
}
if (localeLang.indexOf('en') !== 0) {
    localeLang = 'zh-cn'
}
moment.locale(localeLang)
// Vue.config.debug = true;
Vue.config.productionTip = false
Vue.use(ElementUI, Terminal);

new Vue({
    i18n,
    router,
    methods: util,
    data: data,
    render: h => h(App)
}).$mount('#app')