# 搭建基本的gin框架

## 文件夹结构
config：存储基本的配置信息
gxmodule：存储模型，目前只有用户模型
    user.go：定义了用户模型和一些方法（认证、插入、修改）
middle：存储gin中间件
    authMiddle：定义一个鉴权中间件，为以后各种接口设置统一入口
service：定义业务行为