
本项目为net work cloud的后端

运行项目前需要修改mysql 和 redis的配置
dao
-mysql.go  修改dsn
-redis.go  修改Addr和Password

可去 main.go文件内修改运行端口

配置完毕后输入：
1. go mod init 配置init文件
2. go mod tidy 更新项目依赖
3. go run main.go 运行项目

