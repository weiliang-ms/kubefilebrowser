import Vue from 'vue'
import Router from 'vue-router'
import i18n from '@/lang'

Vue.use(Router)

const _import = file => () => import('../view/' + file + '.vue')

const routes = [
    {
        path: '/',
        component: _import('layer'),
        name: 'console',
        title: i18n.t('console'),
        meta: {
            title: i18n.t('console'),
            icon: 'icon-console',
        },
        redirect: { name: 'dashboard' },
        children: [
            {
                path: 'dashboard',
                name: 'dashboard',
                meta: {
                    title: i18n.t('dashboard'),
                    icon: 'icon-dashboard',
                    single: true,
                },
                component: _import('dashboard'),
            },
            {
                path: 'terminal',
                name: 'terminal',
                meta: {
                    title: i18n.t('terminal'),
                    icon: 'icon-terminal',
                },
                component: _import('terminal'),
            },
            {
                path: 'filebrowser',
                name: 'filebrowser',
                meta: {
                    title: i18n.t('filebrowser'),
                    icon: 'icon-filebrowser',
                },
                component: _import('filebrowser'),
            },
        ],
    },
    {
        path: '/copy',
        name: 'copy',
        title: i18n.t('copy'),
        component: _import('layer'),
        meta: {
            title: i18n.t('copy'),
            icon: 'icon-copy',
        },
        children: [
            {
                path: 'upload',
                name: 'upload',
                meta: {
                    title: i18n.t('upload'),
                },
                component: _import('upload'),
            },
            {
                path: 'multiupload',
                name: 'multiupload',
                meta: {
                    title: i18n.t('multiupload'),
                },
                component: _import('multiupload'),
            },
            {
                path: 'download',
                name: 'download',
                meta: {
                    title: i18n.t('download'),
                },
                component: _import('download'),
            },
        ],
    }
]

const router = new Router({
    routes: routes,
    base: __dirname,
    scrollBehavior: () => ({ y: 0 }),
    mode: 'history',
})
export { routes }
export default router