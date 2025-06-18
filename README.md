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
| `AUTH` | 认证令牌，用于验证请求 | 是 |
| `SMTP_HOST` | SMTP服务器地址 | 是 |
| `SMTP_PORT` | SMTP服务器端口 | 是 |
| `SMTP_USER` | SMTP登录用户名 | 是 |
| `SMTP_PASS` | SMTP登录密码 | 是 |
| `MAIL_FROM` | 发件人邮箱地址 | 是 |
| `MAIL_TO` | 收件人邮箱地址 | 是 |
| `SMTP_ENCRYPTION` | SMTP加密类型：`ssl`/`tls`/`starttls`/`none` | 否，默认为`starttls` |

## 不同语言使用示例

### cURL

```bash
curl -X POST \
  -H "Authorization: Bearer your_token_here" \
  -F "title=我的重要备份" \
  -F "file=@/path/to/your/file.zip" \
  https://your-vercel-deployment-url.vercel.app/api/backup
```

### JavaScript (Fetch API)

```javascript
async function backupFile(file, title) {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('title', title);

  try {
    const response = await fetch('https://your-vercel-deployment-url.vercel.app/api/backup', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer your_token_here'
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error(`备份失败: ${response.status} ${response.statusText}`);
    }

    const result = await response.text();
    console.log('备份结果:', result);
    return result;
  } catch (error) {
    console.error('备份出错:', error);
    throw error;
  }
}

// 使用示例
const fileInput = document.getElementById('fileInput');
const titleInput = document.getElementById('titleInput');
const backupButton = document.getElementById('backupButton');

backupButton.addEventListener('click', async () => {
  if (fileInput.files.length > 0) {
    try {
      await backupFile(fileInput.files[0], titleInput.value);
      alert('备份成功!');
    } catch (error) {
      alert('备份失败: ' + error.message);
    }
  } else {
    alert('请选择要备份的文件');
  }
});
```

### PHP

```php
<?php
function backupFile($filePath, $title = '') {
    $url = 'https://your-vercel-deployment-url.vercel.app/api/backup';
    $token = 'your_token_here';
    
    if (!file_exists($filePath)) throw new Exception('文件不存在: ' . $filePath);
    
    $postFields = [
        'title' => $title,
        'file' => new CURLFile($filePath)
    ];
    
    $ch = curl_init();
    
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $postFields);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $token
    ]);
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    
    if (curl_errno($ch)) {
        throw new Exception('备份请求失败: ' . curl_error($ch));
    }
    
    curl_close($ch);
    
    if ($httpCode !== 200) {
        throw new Exception('备份失败，HTTP状态码: ' . $httpCode . ', 响应: ' . $response);
    }
    
    return $response;
}

// 使用示例
try {
    $result = backupFile('/path/to/your/file.zip', '重要文件备份');
    echo "备份成功: $result\n";
} catch (Exception $e) {
    echo "备份出错: " . $e->getMessage() . "\n";
}
```

### Python

```python
import requests
import os

def backup_file(file_path, title=''):
    """将指定文件上传到备份服务"""
    url = 'https://your-vercel-deployment-url.vercel.app/api/backup'
    token = 'your_token_here'
    
    if not os.path.exists(file_path):
        raise FileNotFoundError(f'文件不存在: {file_path}')
    
    headers = {
        'Authorization': f'Bearer {token}'
    }
    
    files = {
        'file': (os.path.basename(file_path), open(file_path, 'rb'))
    }
    
    data = {}
    if title:
        data['title'] = title
    
    try:
        response = requests.post(url, headers=headers, files=files, data=data)
        response.raise_for_status()  # 抛出HTTP错误
        return response.text
    except requests.exceptions.RequestException as e:
        raise Exception(f'备份请求失败: {e}')
    finally:
        files['file'][1].close()  # 确保文件被关闭

# 使用示例
try:
    result = backup_file('/path/to/your/file.zip', '每日数据备份')
    print(f'备份成功: {result}')
except Exception as e:
    print(f'备份出错: {e}')
```

### Go

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func backupFile(filePath, title, token string) (string, error) {
	url := "https://your-vercel-deployment-url.vercel.app/api/backup"
	
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()
	
	// 创建表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// 添加title字段
	if title != "" {
		if err := writer.WriteField("title", title); err != nil {
			return "", fmt.Errorf("添加title字段失败: %w", err)
		}
	}
	
	// 添加文件
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %w", err)
	}
	
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("复制文件内容失败: %w", err)
	}
	
	// 关闭writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("关闭表单writer失败: %w", err)
	}
	
	// 创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("备份失败，HTTP状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}
	
	return string(respBody), nil
}

func main() {
	filePath := "/path/to/your/file.zip"
	title := "重要数据备份"
	token := "your_token_here"
	
	result, err := backupFile(filePath, title, token)
	if err != nil {
		fmt.Printf("备份出错: %v\n", err)
		return
	}
	
	fmt.Printf("备份成功: %s\n", result)
}
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
