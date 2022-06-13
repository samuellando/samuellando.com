import Vue from 'vue'
import App from './components/view/App.vue'
import { BootstrapVue } from 'bootstrap-vue'

// Import Bootstrap an BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue)

Vue.config.productionTip = false

const routes = {
  '/': App,
}

new Vue({
  render: h => h(routes["/"+window.location.pathname.split("/")[1]])
}).$mount('#app')

