(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-c2fc56e0"],{ad8f:function(t,e,n){"use strict";n.d(e,"a",(function(){return i}));var a=n("b775");function i(t){return Object(a["a"])({url:"/vue-admin-template/table/list",method:"get",params:t})}},b3ca:function(t,e,n){"use strict";n("d65e")},d65e:function(t,e,n){},f30a:function(t,e,n){"use strict";n.r(e);var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"app-container"},[n("div",{staticClass:"top"}),n("div",{staticClass:"base-table"},[n("el-table",{attrs:{data:t.tableData,"element-loading-text":"Loading",border:"",fit:"","highlight-current-row":""}},[n("el-table-column",{attrs:{align:"center",label:"序号",width:"95"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.$index)+" ")]}}])}),n("el-table-column",{attrs:{align:"center",label:"姓名",width:"95"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.row.name)+" ")]}}])}),n("el-table-column",{attrs:{align:"center",label:"性别",width:"95"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.row.sex)+" ")]}}])}),n("el-table-column",{attrs:{align:"center",label:"学号",width:"120"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.row.schoolNumber)+" ")]}}])}),n("el-table-column",{attrs:{label:"专业",prop:"date",align:"center",width:"200"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.row.major)+" ")]}}])}),n("el-table-column",{attrs:{label:"邮箱",align:"center"},scopedSlots:t._u([{key:"default",fn:function(e){return[t._v(" "+t._s(e.row.email)+" ")]}}])}),n("el-table-column",{attrs:{label:"操作",align:"center"},scopedSlots:t._u([{key:"default",fn:function(e){return[n("el-button",{attrs:{type:"white",size:"mini"},on:{click:function(n){return t.upgradePer(e.row)}}},[t._v("升级管理员")]),n("el-button",{attrs:{type:"white",size:"mini"},on:{click:function(n){return t.operate(e.row,"repass")}}},[t._v("重置密码")]),n("el-button",{attrs:{type:"danger",size:"mini"},on:{click:function(n){return t.operate(e.row,"delete")}}},[t._v("删除")])]}}])})],1),n("el-pagination",{staticClass:"pagination",attrs:{background:"","current-page":t.pageNum,"page-size":10,layout:"total,prev, pager, next",total:t.totalNum},on:{"current-change":t.handleCurrentChange}})],1),n("el-dialog",{attrs:{title:"学员申请",visible:t.centerDialogVisible,width:"30%",center:""},on:{"update:visible":function(e){t.centerDialogVisible=e}}},[n("span",[t._v("需要注意的是内容是默认不居中的")]),n("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[n("el-button",{on:{click:function(e){t.centerDialogVisible=!1}}},[t._v("取 消")]),n("el-button",{attrs:{type:"primary"},on:{click:function(e){t.centerDialogVisible=!1}}},[t._v("确 定")])],1)])],1)},i=[],o=n("bc3a"),l=n.n(o),s=(n("ad8f"),{data:function(){return{pageNum:1,totalNum:0,centerDialogVisible:!1,active:0,data:[{id:0,title:"xx"},{id:1,title:"qq"},{id:2,title:"cc"},{id:3,title:"aa"}],currentPage1:5,currentPage2:5,currentPage3:5,currentPage4:4,tableData:[],form:{value:"",input:""},list:null,pickerOptions:{shortcuts:[{text:"最近一周",onClick:function(t){var e=new Date,n=new Date;n.setTime(n.getTime()-6048e5),t.$emit("pick",[n,e])}},{text:"最近一个月",onClick:function(t){var e=new Date,n=new Date;n.setTime(n.getTime()-2592e6),t.$emit("pick",[n,e])}},{text:"最近三个月",onClick:function(t){var e=new Date,n=new Date;n.setTime(n.getTime()-7776e6),t.$emit("pick",[n,e])}}]},value1:[new Date(2e3,10,10,10,10),new Date(2e3,10,11,10,10)],value2:""}},created:function(){this.getStudentList()},methods:{upgradePer:function(t){var e=this;l()({url:"/dev/upIdentity",method:"post",data:{username:t.schoolNumber}}).then((function(t){e.$message.success(t.data.msg),e.getStudentList(),console.log(t)}))},operate:function(t,e){var n=this;l()({url:"/dev/remake_pass_delete",method:"post",data:{activity:e,schoolNumber:t.schoolNumber}}).then((function(t){n.$message.success(t.data.msg),n.getStudentList(),console.log(t)}))},getStudentList:function(){var t=this,e=JSON.parse(sessionStorage.getItem("adminInfo")).data||"",n=e.token;console.log(n),l()({url:"/dev/query_student",method:"post",data:{token:n,pageNum:this.pageNum}}).then((function(e){t.tableData=e.data.data,t.totalNum=e.data.total,console.log(e),"202"==e.data.status&&(t.$message.warning(e.data.msg),sessionStorage.clear(),setTimeout((function(){t.$router.push("/login")}),3e3))}))},change:function(t,e){this.active=e},handleSizeChange:function(t){console.log("每页 ".concat(t," 条"))},handleCurrentChange:function(t){this.pageNum=t,this.getStudentList()}}}),r=s,u=(n("b3ca"),n("2877")),c=Object(u["a"])(r,a,i,!1,null,null,null);e["default"]=c.exports}}]);