package util

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SaveBase64Image(base64Data, savePath string) error {
	// 清理Base64数据（移除Data URL前缀）
	cleanData := cleanBase64Data(base64Data)
	if cleanData == "" {
		return fmt.Errorf("无效的Base64数据")
	}

	// 解码Base64
	imageData, err := base64.StdEncoding.DecodeString(cleanData)
	if err != nil {
		return fmt.Errorf("Base64解码失败: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(savePath, imageData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func Base64ToVideo(base64Str, outputPath string) error {
	// 解码 Base64 字符串
	videoData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("Base64 解码失败: %v", err)
	}

	// 写入视频文件
	err = os.WriteFile(outputPath, videoData, 0644)
	if err != nil {
		return fmt.Errorf("文件写入失败: %v", err)
	}

	fmt.Printf("视频文件已创建: %s (%d bytes)\n", outputPath, len(videoData))
	return nil
}

// 清理Base64数据，移除Data URL前缀
func cleanBase64Data(data string) string {
	if strings.Contains(data, "base64,") {
		parts := strings.Split(data, "base64,")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return data
}
