<template>
    <el-card>
      <el-form class="apply-form first-form" :model="formData" :rules="rule" ref="form">
        <el-form-item :label="$t('namespace')" :prop="namespace">
          <el-select v-model="selected" @click.native="loadMore()">
            <el-option
                v-for="option in NamespceOptions"
                :label="option"
                :value="option"
                :key="option"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
        <el-table
                class="app-table"
                size="medium">
            <el-table-column prop="Pods" label="Pods"></el-table-column>
            <el-table-column prop="Container" label="Container"></el-table-column>
            <el-table-column prop="Image" label="Image"></el-table-column>
            <el-table-column prop="Tag" label="Tag"></el-table-column>
            <el-table-column prop="PullSecret" label="PullSecret"></el-table-column>
            <el-table-column prop="State" label="State"></el-table-column>
            <el-table-column prop="CPU" label="CPU"></el-table-column>
            <el-table-column prop="RAM" label="RAM"></el-table-column>
            <el-table-column prop="Operate" label="Operate"></el-table-column>
        </el-table>
    </el-card>
</template>


<script>
import { Status } from '@/api/status'
import { Namespace } from '@/api/namespaces'
export default {
  data() {
    return {
      selected: "",
      NamespceOptions:[],
      tableLoading: false,
      tableData: [],
    }
  },
  methods: {
    loadMore() {
      Namespace().then(res =>{
        if (res) {
          this.NamespceOptions=[]
          var data = res.items
          for(var key in data){
            this.NamespceOptions.push(data[key].metadata.name)
          };
          console.log(this.NamespceOptions)
        }
      })
    },
  },
}
</script>