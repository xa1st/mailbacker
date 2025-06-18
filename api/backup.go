package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// 全局日志器
var logger *slog.Logger

func init() {
	// 初始化日志器
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// 取出令牌
	auth := os.Getenv("TOKEN")
	if auth == "" {
		logger.Error("环境变量AUTH未设置")
		http.Error(w, "服务器配置错误", http.StatusInternalServerError)
		return
	}

	// 获取请求中的令牌
	authorization := r.Header.Get("Authorization")
	expectedAuth := "Bearer " + auth

	// 检查是否包含正确的Bearer token
	if authorization != expectedAuth {
		logger.Warn("令牌验证失败", "received", authorization)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// 这里获取发送的标题
	title := strings.TrimSpace(r.FormValue("title"))
	// 如果为空给个默认值
	if title == "" {
		title = "[数据备份][未命名数据]" + time.Now().Format("2006-01-02") + "备份"
	}

	// 限制文件上传的大小
	r.ParseMultipartForm(5 << 20) // 最大5M

	// 获取上传的文件
	file, handler, err := r.FormFile("file")

	if err != nil {
		logger.Error("文件获取错误", "error", err)
		http.Error(w, "错误的文件", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		logger.Error("文件读取错误", "error", err)
		http.Error(w, "文件读取错误", http.StatusBadRequest)
		return
	}

	// 这里发送邮件
	logger.Info("准备发送邮件", "filename", handler.Filename, "title", title)
	err = sendMail(title, handler.Filename, fileBytes)
	if err != nil {
		logger.Error("邮件发送错误", "error", err)
		http.Error(w, "邮件发送错误", http.StatusInternalServerError)
		return
	}
	logger.Info("邮件发送成功")
	// 设置内容类型为text/plain
	w.Header().Set("Content-Type", "text/plain")
	// 返回200
	w.WriteHeader(http.StatusOK)
	// 返回成功响应
	fmt.Fprint(w, "备份成功!")
}

func sendMail(title, filename string, fileBytes []byte) error {

	// 获取并验证所有必需的环境变量
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPortStr := os.Getenv("MAIL_PORT")
	smtpUser := os.Getenv("MAIL_MAIL")
	smtpPass := os.Getenv("MAIL_PASS")
	mailFrom := os.Getenv("MAIL_FROM")
	mailTo := os.Getenv("MAIL_TO")
	smtpEncryption := os.Getenv("MAIL_SSL")

	// 验证关键环境变量
	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" || mailFrom == "" || mailTo == "" {
		logger.Error("SMTP配置不完整，缺少必要的环境变量")
		return fmt.Errorf("邮件服务配置不完整")
	}

	// 转换SMTP端口为整数
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		logger.Error("SMTP端口格式错误", "error", err)
		return fmt.Errorf("SMTP端口格式错误: %w", err)
	}

	// 创建SMTP客户端
	smtpClient := mail.NewSMTPClient()
	smtpClient.Host = smtpHost
	smtpClient.Port = smtpPort
	smtpClient.Username = smtpUser
	smtpClient.Password = smtpPass

	// 设置加密类型
	switch strings.ToLower(smtpEncryption) {
	case "ssl", "tls":
		logger.Info("使用SSL/TLS加密")
		smtpClient.Encryption = mail.EncryptionSSLTLS
	case "none", "":
		logger.Info("不使用加密")
		smtpClient.Encryption = mail.EncryptionNone
	default: // 默认使用STARTTLS
		logger.Info("使用STARTTLS加密")
		smtpClient.Encryption = mail.EncryptionSTARTTLS
	}
	// 连接smtp服务器
	logger.Info("正在连接SMTP服务器", "host", smtpHost, "port", smtpPort, "encryption", smtpEncryption)
	conn, err := smtpClient.Connect()
	if err != nil {
		logger.Error("SMTP连接失败", "error", err)
		return fmt.Errorf("SMTP连接失败: %w", err)
	}
	defer conn.Close()

	// 创建邮件
	email := mail.NewMSG()
	// 发件人
	email.SetFrom(mailFrom)
	// 收件人
	email.AddTo(mailTo)
	// 邮件标题
	email.SetSubject(title)
	// 设置邮件正文
	email.SetBody(mail.TextHTML, fmt.Sprintf("<h1>文件备份</h1><p>备份时间: %s</p><p>文件名: %s</p>",
		time.Now().Format("2006-01-02 15:04:05"),
		filename))
	// 添加附件
	email.Attach(&mail.File{Name: filename, Data: fileBytes})
	// 发送邮件
	logger.Info("正在发送邮件", "to", mailTo, "subject", title)
	err = email.Send(conn)
	if err != nil {
		logger.Error("邮件发送失败", "error", err)
		return fmt.Errorf("邮件发送失败: %w", err)
	}
	logger.Info("邮件发送成功")
	return nil
}
