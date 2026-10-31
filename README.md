 [English](./README_US.md) | 中文


## 我的极客时间

不仅是下载器，更是在线文档

 * 支持极客时间VIP账号一次缓存数据，永久观看
 * 支持一键发布整个课程为在线文档
 * 支持一键下载整个课程音视频资源到本地目录



### [在线体验](http://8.141.6.243:8090)


### 环境要求

- Go 1.25+
- Node.js 20+ (用于前端构建)
- Docker & Docker Compose (可选,用于容器化部署)
- MySQL/PostgreSQL/SQLite (数据库,默认使用 SQLite)


### 安装方式

#### 方式一:下载 Release 版本(推荐)

最简单快捷的方式,直接从 GitHub Releases 下载预编译的二进制文件。

1. 访问 [Releases 页面](https://github.com/zkep/my-geektime/releases)
2. 下载对应操作系统的安装包:
   - macOS: `my-geektime-darwin-amd64.tar.gz` 或 `my-geektime-darwin-arm64.tar.gz`
   - Linux: `my-geektime-linux-amd64.tar.gz`
   - Windows: `my-geektime-windows-amd64.zip`

3. 解压并生成配置文件
```shell
# macOS/Linux
tar -xzf my-geektime-*.tar.gz
./my-geektime cli config
./my-geektime server --config=config_templete.yml

# Windows
# 解压 zip 文件后
my-geektime.exe cli config
my-geektime.exe server --config=config_templete.yml
```

4. 访问应用
浏览器打开: http://127.0.0.1:8090



#### 方式二:Docker Compose 部署

最简单快捷的部署方式,适合生产环境使用。

1. 克隆项目
```shell
git clone https://github.com/zkep/my-geektime.git
cd my-geektime/docker
```

2. 启动服务
```shell
docker-compose up -d
```

3. 访问应用
浏览器打开: http://127.0.0.1:8090


#### 方式三:Docker 镜像部署

1. 拉取官方镜像
```shell
# AMD64 架构
docker pull --platform=linux/amd64 zkep/mygeektime:latest

# ARM64 架构(如 Apple Silicon Mac)
docker pull --platform=linux/arm64 zkep/mygeektime:latest
```

2. 准备配置文件
创建 `config.yml` 配置文件(参考项目根目录的 `config.yml` 模板)

3. 运行容器
```shell
docker run -d \
  --name my-geektime \
  -p 8090:8090 \
  -v $(pwd)/config.yml:/config.yml \
  -v $(pwd)/repo:/repo \
  zkep/mygeektime:latest \
  server --config=/config.yml
```


#### 方式四:从源码编译安装

适合开发者和需要自定义构建的用户。

##### 前置条件
- 安装 Go 1.25+
- 安装 Node.js 20+ 和 npm

##### 编译步骤

1. 克隆项目
```shell
git clone https://github.com/zkep/my-geektime.git
cd my-geektime
```

2. 构建前端资源
```shell
make web
```

3. 编译后端
```shell
make build
```

4. 准备配置文件
复制并修改配置模板:
```shell
cp config.yml config.local.yml
# 编辑 config.local.yml 根据需要修改配置
```

5. 运行服务
```shell
./my-geektime server --config=config.local.yml
```

6. 访问应用
浏览器打开: http://127.0.0.1:8090


#### 方式五:直接运行(开发模式)

适合快速测试和开发调试。

1. 克隆项目并进入目录
```shell
git clone https://github.com/zkep/my-geektime.git
cd my-geektime
```

2. 安装依赖并运行
```shell
# 后端依赖
go mod download

# 前端依赖
cd frontend && npm install && cd ..

# 启动开发服务器(前端)
cd frontend && npm run dev

# 在另一个终端启动后端服务
go run main.go server --config=config.yml
```

3. 访问应用
- 前端开发服务器: http://localhost:3000
- 后端 API 服务: http://127.0.0.1:8090


### 配置说明

主要配置文件 `config.yml` 关键配置项:

```yaml
server:
  http_port: 8090          # 服务端口
  
database:
  driver: sqlite           # 数据库类型: sqlite|mysql|postgres
  source: mygeektime.db    # 数据库连接字符串

```

更多配置项请参考 `config.yml` 文件中的注释说明。


### MCP 服务

服务端在 `/mcp` 暴露本地课程、正文和评论的只读访问能力。登录后可在「MCP 服务 → Token 管理」生成不过期的长期 Token，认证方式与客户端接入说明见 [MCP 认证配置](docs/mcp-auth.md)。


### 数据库支持

项目支持多种数据库:

- **SQLite** (默认): 无需额外配置,适合单机部署
- **MySQL**: 适合生产环境,性能更好
- **PostgreSQL**: 备选方案,功能强大

使用 MySQL 示例配置:
```yaml
database:
  driver: mysql
  source: root:password@tcp(127.0.0.1:3306)/mygeektime?charset=utf8&parseTime=True&loc=Local
```


### 故障排查

1. **端口被占用**
   - 修改 `config.yml` 中的 `http_port` 配置
   - 或停止占用 8090 端口的进程

2. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证配置文件中的连接字符串是否正确
   - 确认数据库权限设置

3. **前端资源加载失败**
   - 确保已执行 `make web` 构建前端
   - 检查 `web/` 目录是否存在且包含构建产物

4. **Docker 部署问题**
   - 确认 Docker 和 Docker Compose 版本
   - 检查卷挂载路径权限
   - 查看容器日志: `docker-compose logs -f`


### 微信赞赏

如若有用,不吝赞赏👏

赞赏并留言 <b>邮箱账号</b>,回赠<b> 数据库 *课程参考 [在线体验](http://8.141.6.243:8090)* </b>

<picture>
  <img
    alt="sponsor"
    src="web/sponsor.jpg"
    width="356px"
  />
</picture>


### 许可证

本项目采用 MIT 许可证,详见 [LICENSE](LICENSE) 文件。



