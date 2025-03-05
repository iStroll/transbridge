# TransBridge 配置详解

本文档详细说明了 TransBridge 的配置选项，以帮助您根据需求进行定制。

配置文件采用 YAML 格式，默认查找当前目录下的 `config.yml` 文件，也可以通过命令行参数 `-config` 指定配置文件路径。

## 配置结构

完整的配置结构如下：

```yaml
server:
  port: 8080

providers:
  - provider: "openai"
    api_url: "https://api.openai.com/v1/chat/completions"
    api_key: "your-api-key-1"
    timeout: 30
    is_default: true
    models:
      - name: "gpt-3.5-turbo"
        weight: 10
        max_tokens: 2000
        temperature: 0.3

  - provider: "zhipuai"
    api_url: "https://open.bigmodel.cn/api/paas/v4/chat/completions"
    api_key: "your-api-key-2"
    timeout: 30
    models:
      - name: "chatglm_turbo"
        weight: 5
        max_tokens: 2000
        temperature: 0.3

cache:
  enabled: true
  types: ["memory", "redis"]
  
  memory:
    ttl:
      value: "1h"
    max_size: 10000
  
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    ttl: 
      value: "7d"

prompt:
  template: "Translate the following text from %s to %s:\n\n%s"

auth:
  tokens:
    - "tr-xxxxxxxxxxxxxxxx"

log:
  enabled: true
  file_path: "logs/translation.log"
  max_size: 100
  max_age: 30
  max_backups: 10
  queue_size: 1000
```

## 服务器配置

```yaml
server:
  port: 8080  # 服务监听端口
```

## 提供商配置

`providers` 部分配置翻译服务提供商，支持多个提供商和模型。

```yaml
providers:
  - provider: "openai"                                       # 提供商标识
    api_url: "https://api.openai.com/v1/chat/completions"    # API 完整 URL
    api_key: "your-api-key-1"                                # API 密钥
    timeout: 30                                              # 请求超时时间(秒)
    is_default: true                                         # 是否为默认提供商
    models:                                                  # 模型配置
      - name: "gpt-3.5-turbo"                                # 模型名称
        weight: 10                                           # 负载均衡权重
        max_tokens: 2000                                     # 最大 token 数
        temperature: 0.3                                     # 温度参数
```

### 提供商参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider | 字符串 | 是 | 提供商标识，用于日志和引用 |
| api_url | 字符串 | 是 | API 完整 URL 地址 |
| api_key | 字符串 | 是 | API 认证密钥 |
| timeout | 整数 | 否 | 请求超时时间(秒)，默认 30 |
| is_default | 布尔值 | 否 | 是否为默认提供商，系统中需要有且仅有一个默认提供商 |
| models | 数组 | 是 | 模型配置数组 |

### 模型参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | 字符串 | 是 | 模型名称，需要符合提供商支持的模型 |
| weight | 整数 | 是 | 负载均衡权重，权重越大被选中概率越高 |
| max_tokens | 整数 | 否 | 最大输出 token 数量，默认 2000 |
| temperature | 浮点数 | 否 | 生成多样性参数，0-1 之间，默认 0.3 |

## 缓存配置

TransBridge 支持多级缓存，可同时启用内存缓存和 Redis 缓存。

```yaml
cache:
  enabled: true                   # 是否启用缓存
  types: ["memory", "redis"]      # 缓存类型，可选 memory 和 redis
  
  # 内存缓存配置
  memory:
    ttl:
      value: "1h"                 # 缓存有效期
    max_size: 10000               # 最大缓存条目数
  
  # Redis 缓存配置
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    ttl: 
      value: "7d"                 # Redis 缓存有效期
```

### 缓存 TTL 值格式

TTL (Time To Live) 支持以下格式：

| 值 | 说明 |
|------|------|
| "30s" | 30 秒 |
| "5m" | 5 分钟 |
| "2h" | 2 小时 |
| "1d" | 1 天 |
| "1w" | 1 周 |
| "permanent" | 永久存储 |
| "0" | 永久存储 |

## 提示模板配置

```yaml
prompt:
  template: "Translate the following text from %s to %s:\n\n%s"
```

提示模板用于构建发送给 AI 服务的提示，包含三个参数：源语言、目标语言和文本内容。

## 认证配置

```yaml
auth:
  tokens:                       # API 密钥列表
    - "tr-xxxxxxxxxxxxxxxx"
    - "tr-yyyyyyyyyyyyyyyy"
```

客户端请求翻译接口时，需要提供其中一个 API 密钥进行认证。

## 日志配置

```yaml
log:
  enabled: true                        # 是否启用日志
  file_path: "logs/translation.log"    # 日志文件路径
  max_size: 100                        # 单个日志文件最大大小(MB)
  max_age: 30                          # 保留日志文件的最大天数，0表示永久保留
  max_backups: 10                      # 最大备份文件数，0表示保留所有备份
  queue_size: 1000                     # 异步队列大小
```

### 日志参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| enabled | 布尔值 | 是 | 是否启用日志记录 |
| file_path | 字符串 | 是 | 日志文件路径 |
| max_size | 整数 | 否 | 单个日志文件最大大小(MB)，默认 100 |
| max_age | 整数 | 否 | 保留日志文件的最大天数，0表示永久保留 |
| max_backups | 整数 | 否 | 最大备份文件数，0表示保留所有备份 |
| queue_size | 整数 | 否 | 异步队列大小，默认 1000 |

## 完整配置示例

请参考项目根目录下的 `config.example.yml` 文件获取最新的配置示例。

## 环境变量支持

除了配置文件外，TransBridge 还支持通过环境变量覆盖配置。环境变量的命名规则是将配置路径转换为大写，并用下划线连接，例如：

- `SERVER_PORT=8080`
- `PROVIDERS_0_API_KEY=your-api-key`
- `CACHE_ENABLED=true`

环境变量优先级高于配置文件。