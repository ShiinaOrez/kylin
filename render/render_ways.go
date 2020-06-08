package render

import (
	"context"
	"errors"
	"os"
)

func SaveAsFile(ctx context.Context, content string) error {
	var (
		id       string
		path     string
		filename string
		ok       bool
	)
	if path, ok = ctx.Value("path").(string); !ok {
		return errors.New("Invalid param `path`.")
	}
	if id, ok = ctx.Value("id").(string); !ok {
		return errors.New("Invalid param `id`.")
	}
	filename = id+"-result-data.txt"
	bytes := []byte(content)
	file, err := os.Create(path+"/"+filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func SendToEmail(ctx context.Context, content string) {
	// TODO: 发送邮件
}

func GeneratePDF(ctx context.Context, content string) {
	// TODO: 生成 PDF
}