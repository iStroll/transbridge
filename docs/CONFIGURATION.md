# TransBridge 配置指南 ⚙️

## 目录
- [配置文件概述](#配置文件概述)
- [服务器配置](#服务器配置)
- [提供商配置](#提供商配置)
- [缓存配置](#缓存配置)
- [认证配置](#认证配置)
- [日志配置](#日志配置)
- [完整配置示例](#完整配置示例)

## 配置文件概述

TransBridge 使用 YAML 格式的配置文件。默认配置文件名为 `config.yml`，您也可以通过命令行参数指定其他配置文件：

```bash
./transbridge -config custom_config.yml
```

## 服务器配置

服务器配置部分控制 TransBridge 的基本运行参数。

```yaml
server:
  port: 8080           # 服务监听端口
  host: "0.0.0.0"      # 服务监听地址，默认所有地址
```

## 提供商配置

提供商配置是最核心的部分，支持配置多个翻译服务提供商。

### OpenAI 配置示例
```yaml
providers:
  - provider: "openai"
    api_url: "https://api.openai.com/v1/chat/completions"  # API 地址
    api_key: "your-api-key"                                # API 密钥
    timeout: 30                                            # 请求超时时间（秒）
    is_default: true                                       # 是否为默认提供商
    models:
      - name: "gpt-3.5-turbo"                             # 模型名称
        weight: 10                                         # 负载均衡权重
        max_tokens: 2000                                   # 最大 token 数
        temperature: 0.3                                   # 温度参数
```

### ChatGLM 配置示例
```yaml
providers:
  - provider: "chatglm"
    api_url: "http://localhost:8000/v1/chat/completions"
    timeout: 30
    is_default: false
    models:
      - name: "chatglm3-6b"
        weight: 5
        max_tokens: 2000
        temperature: 0.3
```

### Ollama 配置示例
```yaml
providers:
  - provider: "ollama"
    api_url: "http://localhost:11434/api/chat"
    timeout: 30
    is_default: false
    models:
      - name: "llama2"
        weight: 5
        max_tokens: 2000
        temperature: 0.3
```

### 提供商配置参数说明

| 参数 | 说明 | 默认值 | 是否必填 |
|------|------|--------|----------|
| provider | 提供商类型 | - | 是 |
| api_url | API 接口地址 | - | 是 |
| api_key | API 密钥 | - | 部分必填 |
| timeout | 请求超时时间（秒） | 30 | 否 |
| is_default | 是否为默认提供商 | false | 否 |

### 模型配置参数说明

| 参数 | 说明 | 默认值 | 是否必填 |
|------|------|--------|----------|
| name | 模型名称 | - | 是 |
| weight | 负载均衡权重 | 1 | 否 |
| max_tokens | 最大生成 token 数 | 2000 | 否 |
| temperature | 采样温度 | 0.3 | 否 |

## 缓存配置

缓存配置支持内存缓存和 Redis 缓存两种方式，可以同时启用。

### 内存缓存配置
```yaml
cache:
  enabled: true                # 是否启用缓存
  types: ["memory"]           # 缓存类型：memory 或 redis
  memory:
    ttl:
      value: "1h"            # 缓存过期时间（支持：30s, 5m, 2h, 1d, 1w, permanent）
    max_size: 10000          # 最大缓存条目数
```

### Redis 缓存配置
```yaml
cache:
  enabled: true
  types: ["redis"]
  redis:
    host: "localhost"         # Redis 服务器地址
    port: 6379               # Redis 端口
    password: ""             # Redis 密码
    db: 0                    # Redis 数据库编号
    ttl:
      value: "24h"          # 缓存过期时间
```

### 缓存配置参数说明

| 参数 | 说明 | 默认值 | 是否必填 |
|------|------|--------|----------|
| enabled | 是否启用缓存 | false | 是 |
| types | 缓存类型列表 | [] | 是 |
| ttl.value | 缓存过期时间 | "1h" | 否 |
| max_size | 最大缓存条目数 | 10000 | 否 |

types 可以取值 ["memory"] ["redis"] 和["memory", "redis"]

## 认证配置

配置 API 访问认证信息。

```yaml
transapi:
  tokens:                    # API 密钥列表
    - "your-api-key-1"
    - "your-api-key-2"
```

## 日志配置

配置日志记录相关参数。

```yaml
log:
  enabled: true                        # 是否启用日志
  file_path: "logs/translation.log"    # 日志文件路径
  max_size: 100                        # 单个日志文件最大大小（MB）
  max_age: 30                         # 日志文件保留天数
  max_backups: 10                     # 最大备份文件数
  queue_size: 1000                    # 异步日志队列大小
```

## 完整配置示例

下面是一个包含所有主要配置项的完整示例：

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30
  write_timeout: 30

providers:
  - provider: "openai"
    api_url: "https://api.openai.com/v1/chat/completions"
    api_key: "your-openai-key"
    timeout: 30
    is_default: true
    models:
      - name: "gpt-3.5-turbo"
        weight: 10
        max_tokens: 2000
        temperature: 0.3

  - provider: "ollama"
    api_url: "http://localhost:11434/api/chat"
    timeout: 30
    is_default: false
    models:
      - name: "llama2"
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
      value: "24h"

transapi:
  tokens:
    - "your-api-key-1"
    - "your-api-key-2"

log:
  enabled: true
  file_path: "logs/translation.log"
  max_size: 100
  max_age: 30
  max_backups: 10
  queue_size: 1000

prompt:
  template: "Translate the following {{source_lang}} content to {{target_lang}}: {{input}}"
```