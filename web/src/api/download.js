import {get} from '../lib/fetch.js'

export function Download(params) {
    const headers = {"responseType":"blob"}
    return get('/k8s/download', params,headers)
}