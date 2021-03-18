import {post} from '@/lib/fetch.js'

export function Upload(params) {
    return post('/k8s/upload', params)
}