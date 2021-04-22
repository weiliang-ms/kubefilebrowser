import {get} from '../lib/fetch.js'
import {post} from '../lib/fetch.js'

export function FileBrowserList(params) {
    return get('/file_browser/list', params)
}

export function FileBrowserOpen(params) {
    return get('/file_browser/open', params)
}

export function FileBrowserCreateFile(data, params) {
    return post('/file_browser/create_file', data, params, {"Content-Type":"multipart/form-data"})
}

export function FileBrowserCreateDir(params) {
    return post('/file_browser/create_dir', "", params, {"Content-Type":"multipart/form-data"})
}

export function FileBrowserRename(params) {
    return post('/file_browser/rename', "", params, {"Content-Type":"multipart/form-data"})
}

export function FileBrowserRemove(params) {
    return post('/file_browser/remove', "", params, {"Content-Type":"multipart/form-data"})
}