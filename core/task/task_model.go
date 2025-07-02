package task

import "time"

// DownloadTask 下载任务模型
type DownloadTask struct {
	ID           string    `json:"id"`           // 任务ID
	URL          string    `json:"url"`          // 下载URL
	Status       string    `json:"status"`       // 状态: pending, downloading, completed, failed
	SavePath     string    `json:"savePath"`     // 保存路径
	StartTime    time.Time `json:"startTime"`    // 开始时间
	CompleteTime time.Time `json:"completeTime"` // 完成时间
	UpdatedAt    time.Time `json:"updatedAt"`    // 更新时间
	Error        string    `json:"error"`        // 错误信息
	Name         string    `json:"name"`         // 任务名
	Progress     struct {
		Current int `json:"current"` // 当前已下载项目数
		Total   int `json:"total"`   // 总项目数
	} `json:"progress"` // 下载进度
}

// DownloadStatus 表示下载任务状态
type DownloadStatus string

const (
	StatusPending     DownloadStatus = "pending"     // 等待下载
	StatusDownloading DownloadStatus = "downloading" // 下载中
	StatusCompleted   DownloadStatus = "completed"   // 下载完成
	StatusFailed      DownloadStatus = "failed"      // 下载失败
	StatusCancelled   DownloadStatus = "cancelled"   // 已取消
)
