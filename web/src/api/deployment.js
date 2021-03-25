import {get} from '../lib/fetch.js'

export function GetDeployment(params) {
    return get('/k8s/deployment', params)
}