import {get} from '@/lib/fetch.js'

export function Pods(params) {
    return get('/k8s/pods', params)
}