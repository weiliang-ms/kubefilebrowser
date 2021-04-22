import {get} from '../lib/fetch.js'

export function GetStatus(params) {
    return get('/k8s/status', params)
}