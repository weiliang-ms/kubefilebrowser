import {get} from '@/lib/fetch.js'

export function Status(params) {
    return get('/k8s/status', params)
}