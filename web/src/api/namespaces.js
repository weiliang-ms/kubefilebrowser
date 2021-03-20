import {get} from '@/lib/fetch.js'

export function GetNamespace(params) {
    return get('/k8s/namespace', params)
}