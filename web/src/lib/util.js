import moment from 'moment'

export default {
    MessageSuccess(cb){
        this.$message({
            message: this.$t('operate_success'),
            type: 'success',
            duration: 1200,
            onClose: cb,
        });
    },

    PageInit() {
        this.$root.PageSize = 7
        this.$root.Page = 1
        this.$root.Total = 0
    },

    PageReset() {
        this.$root.Total--
        let maxPage = Math.ceil(this.$root.Total / this.$root.PageSize)
        if (this.$root.Page > maxPage) {
            this.$root.Page = maxPage
        }
        if (this.$root.Page < 1) {
            this.$root.Page = 1
        }
    },

    PageOffset() {
        return this.$root.PageSize * (this.$root.Page - 1)
    },

    ConfirmDelete(cb, title) {
        if (!title) {
            title = '此操作将永久删除该数据, 是否继续?'
        }
        this.$confirm(title, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        }).then(() => {
            cb()
        }).catch((err) => {
            console.log(err)
        })
    },

    FormatDateTime(unixtime, format) {
        if (!unixtime) {
            return '--'
        }
        if (!format) {
            format = 'YYYY-MM-DD HH:mm:ss'
        }
        return moment.unix(unixtime).format(format)
    },

    FormatDateDuration(num) {
        return moment.duration(num).humanize(false)
    },

    FormatDateFromNow(unixtime) {
        if (!unixtime) {
            return '--'
        }
        return moment.unix(unixtime).fromNow()
    },

    Substr(str, len) {
        if (Object.prototype.toString.call(str) != '[object String]') {
            return ''
        }
        let postfix = ''
        if (str.length > len) {
            postfix = "..."
        }
        return str.substr(0, len) + postfix
    }
}