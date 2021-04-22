import {post} from '../lib/fetch.js'

export function MultiUpload(data, params, headers) {
    return post('/k8s/multi_upload', data, params, headers)
}