package files

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	file       *os.File
	goodFile   *os.File
	badFile    *os.File
	Reader     *csv.Reader
	GoodWriter *csv.Writer
	BadWriter  *csv.Writer
	sync.Once
}

// NewFileProcessor 初始化文件处理器
func NewFileProcessor(maxFileSize int64) (*FileProcessor, error) {
	fp := &FileProcessor{}
	if err := fp.openData(); err != nil {
		return nil, err
	}
	if err := fp.createGoodFile(); err != nil {
		return nil, err
	}
	if err := fp.createBadFile(); err != nil {
		return nil, err
	}
	return fp, nil
}

// FlushClose 刷新并关闭文件
func (f *FileProcessor) FlushClose() {
	f.Once.Do(func() {
		_ = f.file.Close()
		f.GoodWriter.Flush()
		_ = f.goodFile.Close()
		f.BadWriter.Flush()
		_ = f.badFile.Close()
	})
}
func (f *FileProcessor) openData() (err error) {
	f.file, err = os.OpenFile("./demo/mullinkcheck/files/data.csv", os.O_RDONLY, 0666)
	//f.file, err = os.OpenFile("./data.csv", os.O_RDONLY, 0666)
	if err != nil {
		return
	}

	f.Reader = csv.NewReader(f.file)
	return
}
func (f *FileProcessor) createGoodFile() (err error) {
	f.goodFile, err = os.OpenFile("./demo/mullinkcheck/files/good.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//f.goodFile, err = os.OpenFile("./good.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	bufWriter := bufio.NewWriterSize(f.goodFile, 4096)
	f.GoodWriter = csv.NewWriter(bufWriter)
	return
}
func (f *FileProcessor) createBadFile() (err error) {
	f.badFile, err = os.OpenFile("./demo/mullinkcheck/files/bad.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//f.badFile, err = os.OpenFile("./bad.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	bufWriter := bufio.NewWriterSize(f.badFile, 4096)
	f.BadWriter = csv.NewWriter(bufWriter)
	return
}

// WriteRow 写入一行数据
func (f *FileProcessor) WriteRow(row []string, isValid bool) {
	var writer *csv.Writer
	if isValid {
		writer = f.GoodWriter
	} else {
		writer = f.BadWriter
	}
	err := writer.Write(row)
	if err != nil {
		fmt.Println(err)
	}
}
