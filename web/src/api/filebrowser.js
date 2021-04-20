import {get} from '../lib/fetch.js'
import {post} from '../lib/fetch.js'

export function FileBrowserList(params) {
    return get('/file_browser/list', params)
}

export function FileBrowserOpen(params) {
    return get('/file_browser/open', params)
}

export function FileBrowserCreateFile(params) {
    return post('/file_browser/create_file', params)
}

export function FileBrowserCreateDir(params) {
    return post('/file_browser/create_dir', params)
}

export function FileBrowserRename(params) {
    return post('/file_browser/rename', params)
}

export function FileBrowserRemove(params) {
    return post('/file_browser/remove', params)
}