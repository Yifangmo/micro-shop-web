# micro-shop-web
## 概述
这是一个实现简单功能的商城Web层服务，使用Gin框架来实现Restful API，与[其下层 gRPC 服务层](https://github.com/Yifangmo/micro-shop-services)组成一个商城后端系统，其代码需要依赖 gRPC 服务层所定义的接口。系统架构如下：
![商城后端系统架构.png](https://s2.loli.net/2023/11/07/FvBnQXA8Ob2DzWk.png)

## 功能模块
| 子目录 | 服务 | 描述 |
| ---  | ---- | ---- |
| user | 用户服务 | 包含用户登录注册、收货地址管理、商品收藏、留言 API |
| goods | 商品服务 | 包含商品及其分类、品牌及其分类、轮播图 API |
| oss | 对象存储服务 | 包含阿里云对象存储服务 token 及对象的获取、上传回调 API|
| order  | 订单服务 | 包含对订单的新增查询、支付通知回调和购物车增删改查操作 API |

## 目录结构
    |-- go.mod
    |-- user   对应一个服务模块
        |-- main.go
        |-- run.sh   简单的运行命令
        |-- configs  存放配置文件、Nacos和服务的配置结构体
        |-- global   存放全局使用的变量，如服务器配置、gRPC客户端、是否为debug等
        |-- handlers 存放所有的 gin.HandlerFunc
            |-- server.go 服务器构造，使用 gin.Engine 作为 http.Handler
            |-- user 存放某个资源对应的路径下的 gin.HandlerFunc 和路由组注册函数
            |-- ... 其他路径(同上)
        |-- initialize 存放全局变量的初始化或功能注册逻辑
        |-- middlewares 存放 gin 框架使用的中间件
        |-- models 存放请求绑定的表单模型和响应模型
        |-- utils 存放通用的工具函数，如错误处理或转换、服务注册
        |-- validators 存放表单中各种类型的校验规则
    |-- ... 其他服务模块(同上)

## 服务启动
1. 在 `${service}/configs/dev.yaml` 配置好 Nacos 服务器信息，并在 Nacos 中配置好 `ServerConfig`
2. 在 `go.mod` 同级目录运行 `go generate`
3. 进入各服务的子目录，如 `cd user/`，运行 `./run.sh`
