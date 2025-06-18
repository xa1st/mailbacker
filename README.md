# Mailbacker

一个轻量级的邮件备份服务，基于Vercel Serverless Functions。通过简单的HTTP POST请求，可以将文件作为附件发送到指定的邮箱，实现快速备份。

## 功能特点

- 部署在Vercel上的Serverless服务
- 支持文件上传并通过邮件附件发送
- 支持Token认证，确保API安全
- 支持多种SMTP加密方式（SSL/TLS, STARTTLS, 无加密）
- 完整的日志记录，使用Go的slog库

## 环境变量配置

部署前，需要在Vercel平台配置以下环境变量：

| 变量名 | 描述 | 是否必须 |
|--------|------|----------|
| `TOKEN` | 认证令牌，用于验证请求 | 是 |
| `MAIL_HOST` | SMTP服务器地址 | 是 |
| `MAIL_PORT` | SMTP服务器端口 | 是 |
| `MAIL_MAIL` | SMTP登录用户名 | 是 |
| `MAIL_PASS` | SMTP登录密码 | 是 |
| `MAIL_FROM` | 发件人邮箱地址 | 是 |
| `MAIL_TO` | 收件人邮箱地址 | 是 |
| `MAIL_SSL` | SMTP加密类型：`ssl`/`tls`/`starttls`/`none` | 否，默认为`starttls` |

## API响应格式

所有API响应均以JSON格式返回，格式如下：

```json
{
  "code": 200,
  "data": {},
  "msg": "备份成功"
}
```

### 成功响应示例

```json
{
  "code": 200,
  "data": {
    "filename": "backup.zip",
    "title": "重要数据备份",
    "size": 1024567,
    "timestamp": "2025-06-18T09:00:00Z"
  },
  "msg": "备份成功"
}
```

### 错误响应示例

```json
{
  "code": 401,
  "data": {
    "error": "无效的认证令牌"
  },
  "msg": "认证失败"
}
```

## 不同语言使用示例

### cURL

```bash
curl -X POST \
  -H "Authorization: Bearer your_token_here" \
  -F "title=我的重要备份" \
  -F "file=@/path/to/your/file.zip" \
  https://your-vercel-deployment-url.vercel.app/api/backup
```

## 部署到Vercel

1. 创建一个GitHub仓库，并将代码推送到该仓库
2. 登录[Vercel](https://vercel.com/)并导入该仓库
3. 在部署设置中配置所有必要的环境变量
4. 点击部署按钮

## 限制

- 文件大小限制：5MB
- 执行时间：最长60秒（可在vercel.json中配置）

## 贡献

欢迎提交Issue和Pull Request。

## 许可证

MIT
