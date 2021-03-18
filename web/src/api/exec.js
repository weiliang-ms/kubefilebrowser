import {get} from '@/lib/fetch.js'

export function Exec(params) {
    return get('/k8s/exec', params)
}