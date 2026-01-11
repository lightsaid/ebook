package main

// 提供给swagger文档生成通用的出参类型定义
type ApiResponse struct {
	BizCode   string `json:"bizCode"`             // 业务编码
	Message   string `json:"message"`             // 客户消息
	Data      any    `json:"data"`                // 任意数据
	Version   string `json:"version,omitempty"`   // 版本信息
	RequestID string `json:"requestId,omitempty"` // 请求Id，做简单的链路追踪
}
