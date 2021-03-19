module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:9999/',
                changeOrigin: true,
            }
        }
    },
    pages: {
        index: {
            // page 的入口
            entry: 'src/main.js',
            // 模板来源
            template: 'public/index.html',
            // 当使用 title 选项时，
            // template 中的 title 标签需要是 <title><%= htmlWebpackPlugin.options.title %></title>
            title: 'KubernetesFileBrowser',
        }
    },
    publicPath: '/static/'
}
