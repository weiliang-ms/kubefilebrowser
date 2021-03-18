import {get} from '@/lib/fetch.js'

export function Namespace(params) {
    return get('/k8s/namespace', params)
}