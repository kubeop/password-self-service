# 1. 项目介绍
password-self-service 是一个基于[Gin](https://gin-gonic.com)开发的基于微软AD域控的密码自助平台，帮助企业员工快速重置密码和解锁账号，减少对IT运维的依赖。

# 2. 技术选型
- 后端：用 [Gin](https://gin-gonic.com/) 快速搭建基础restful风格API，[Gin](https://gin-gonic.com/) 是一个go语言编写的Web框架。
- 缓存：使用`Redis`实现记录当前活跃用户的`jwt`令牌并实现多点登录限制。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用 [viper](https://github.com/spf13/viper) 实现`yaml`格式的配置文件。
- 日志：使用 [zap](https://github.com/uber-go/zap) 实现日志记录。



# 3. 功能

- [x] 重置密码
- [x] 解锁账户
- [x] 密码过期通知
- [x] 支持邮件方式发送消息
- [x] 支持阿里云短信方式发送消息
- [x] 支持腾讯云短信方式发送消息
- [ ] 支持钉钉应用方式发送消息
- [ ] 支持企业微信方式发送消息



# 4. 部署服务

使用[docker-compose](manifest/docker/docker-compose.yml)部署。

使用[kubernetes](manifest/k8s)部署。


# 5. 二次开发

```shell
# 拉取代码
git clone https://github.com/kubeop/password-self-service.git

# 安装swag
go install github.com/swaggo/swag/cmd/swag@latest

# 下载依赖并生成swagger文档
make init

# 启动服务
make run
```
