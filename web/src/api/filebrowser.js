import {get} from '../lib/fetch.js'

export function FileBrowser(params) {
    return get('/k8s/file_browser', params)
}