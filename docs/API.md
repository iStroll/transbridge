# TransBridge API 文档

TransBridge 提供以下 API 接口：

## 翻译接口（DeepL 兼容）

### 请求

```
POST /v2/translate
```

#### 请求头

| 名称 | 必填 | 描述 |
|------|------|------|
| Authorization | 是 | Bearer 认证，格式为 `Bearer YOUR_API_KEY` |
| Content-Type | 是 | 固定为 `application/json` |

#### 请求体

```json
{
  "text": "Hello world",
  "source_lang": "EN",
  "target_lang": "ZH",
  "provider": "openai",
  "model": "gpt-3.5-turbo"
}
```

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| text | 字符串 | 是 | 要翻译的文本 |
| source_lang | 字符串 | 是 | 源语言代码，例如 "EN", "ZH" |
| target_lang | 字符串 | 是 | 目标语言代码，例如 "EN", "ZH" |
| provider | 字符串 | 否 | 指定服务提供商，不填则随机选择 |
| model | 字符串 | 否 | 指定模型名称，不填则随机选择 |

#### 支持的语言代码

| 代码 | 语言 |
|------|------|
| EN | 英语 |
| ZH | 中文 |
| JA | 日语 |
| KO | 韩语 |
| ES | 西班牙语 |
| FR | 法语 |
| DE | 德语 |
| IT | 意大利语 |
| RU | 俄语 |
| PT | 葡萄牙语 |
| NL | 荷兰语 |
| PL | 波兰语 |
| ... | 其他语言 |

### 响应

```json
{
  "code": 200,
  "data": "你好世界",
  "source_lang": "EN",
  "target_lang": "ZH"
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| code | 数字 | 状态码，200 表示成功 |
| data | 字符串 | 翻译后的文本 |
| source_lang | 字符串 | 源语言代码 |
| target_lang | 字符串 | 目标语言代码 |

### 错误响应

```json
{
  "code": 401,
  "data": "Invalid API key"
}
```

| 状态码 | 描述 |
|------|------|
| 400 | 请求参数错误 |
| 401 | 未授权（API 密钥无效） |
| 500 | 服务器内部错误 |

### 示例

#### cURL

```bash
curl -X POST "http://localhost:8080/v2/translate" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello world",
    "source_lang": "EN",
    "target_lang": "ZH"
  }'
```

#### Python

```python
import requests

url = "http://localhost:8080/v2/translate"
headers = {
    "Authorization": "Bearer your-api-key",
    "Content-Type": "application/json"
}
data = {
    "text": "Hello world",
    "source_lang": "EN",
    "target_lang": "ZH"
}

response = requests.post(url, headers=headers, json=data)
print(response.json())
```

#### JavaScript

```javascript
fetch("http://localhost:8080/v2/translate", {
  method: "POST",
  headers: {
    "Authorization": "Bearer your-api-key",
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    text: "Hello world",
    source_lang: "EN",
    target_lang: "ZH"
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

## OpenAI 兼容接口

TransBridge 还提供与 OpenAI API 兼容的接口，可以直接替代 OpenAI 的聊天完成接口。

### 请求

```
POST /v1/chat/completions
```

#### 请求头

| 名称 | 必填 | 描述 |
|------|------|------|
| Authorization | 是 | Bearer 认证，格式为 `Bearer YOUR_API_KEY` |
| Content-Type | 是 | 固定为 `application/json` |

#### 请求体

```json
{
  "model": "openai/gpt-3.5-turbo",
  "messages": [
    {
      "role": "system",
      "content": "You are a helpful assistant."
    },
    {
      "role": "user",
      "content": "Hello!"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 100
}
```

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| model | 字符串 | 是 | 模型名称，格式为 "provider/model" |
| messages | 数组 | 是 | 消息数组，包含多个消息对象 |
| temperature | 浮点数 | 否 | 温度参数，控制生成文本的随机性 |
| max_tokens | 整数 | 否 | 最大输出 token 数量 |

### 响应

响应格式与 OpenAI API 保持一致：

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677858242,
  "model": "openai/gpt-3.5-turbo",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ]
}
```

## 健康检查接口

### 请求

```
GET /health
```

### 响应

```
OK
```

状态码 200 表示服务正常运行。