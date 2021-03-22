import Vue from 'vue'
import Router from 'vue-router'
import i18n from '@/lang'

Vue.use(Router)

const _import = file => () => import('../view/' + file + '.vue')

const routes = [
    {
        path: '/',
        component: _import('layer'),
        name: 'dashboard',
        title: i18n.t('dashboard'),
        redirect: { name: 'dashboard' },
        children: [
            {
                path: 'dashboard',
                name: 'dashboard',
                meta: {
                    title: i18n.t('dashboard'),
                    icon: 'el-icon-monitor',
                    single: true,
                },
                component: _import('dashboard'),
            }
        ],
    },
    {
        path: '/copy',
        name: 'copy',
        title: i18n.t('copy'),
        component: _import('layer'),
        meta: {
            title: i18n.t('copy'),
            icon: 'el-icon-upload',
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