import { createRouter, createMemoryHistory, createWebHistory, createWebHashHistory } from 'vue-router'
import routes from './routes'

export default function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,

    // Leave this as is and make changes in quasar.config.js instead!
    // quasar.config.js -> build -> vueRouterMode
    // quasar.config.js -> build -> publicPath
    history: createHistory(process.env.MODE === 'ssr' ? void 0 : process.env.VUE_ROUTER_BASE)
  })

  Router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.authRequired)) {
      if (JSON.parse(localStorage.getItem('auth')) === true) {
        next();
      } else {
        next({
          path: '/login',
          params: { nextUrl: to.fullPath }
        });
      }
    } else {
      next();
    }
  });

  return Router
}