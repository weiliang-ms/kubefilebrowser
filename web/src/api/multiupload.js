import {post} from '@/lib/fetch.js'

export function MultiUpload(params) {
    return post('/k8s/multi_upload', params)
}