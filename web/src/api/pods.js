import {get} from '../lib/fetch.js'

export function GetPods(params) {
    return get('/k8s/pods', params)
}