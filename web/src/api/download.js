import {get} from '@/lib/fetch.js'

export function Download(params) {
    return get('/k8s/download', params)
}