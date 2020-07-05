const routes = [
  {
    path: '/app',
    component: () => import('layouts/AppLayout.vue'),
    children: [{ path: '', component: () => import('pages/Dashboard.vue') }],
    meta: {
      authRequired: true
    }
  },
  {
    path: '/',
    component: () => import('layouts/AuthLayout.vue'),
    children: [
      { path: '', component: () => import('pages/Homepage.vue') },
      { path: '/login', component: () => import('pages/Login.vue') },
      { path: '/register', component: () => import('pages/Login.vue') }
    ]
  },
  {
    path: '*',
    component: () => import('pages/Error404.vue')
  }
];

export default routes;
