import {post} from '../lib/fetch.js'

export function FileOrDirUpload(data, params, headers) {
    return post('/k8s/upload', data, params, headers)
}