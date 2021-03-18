import {get} from '@/lib/fetch.js'

export function Terminal(params) {
    return get('/k8s/terminal', params)
}