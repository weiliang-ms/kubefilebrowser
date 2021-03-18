import Vue from 'vue'

Vue.directive("loadmore", {
    bind(el, binding, vnode) {
        const SELECTWRAP = el.querySelector(
            ".el-select-dropdown .el-select-dropdown__wrap"
        );
        SELECTWRAP.addEventListener("scroll", function () {
            const CONDITION = this.scrollHeight - this.scrollTop <= this.clientHeight;
            if (CONDITION) {
                binding.value();
            }
        });
    },
});