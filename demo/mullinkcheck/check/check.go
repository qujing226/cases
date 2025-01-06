package check

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qujing226/cases/demo/mullinkcheck/files"
	"net/http"
	"strings"
	"sync"
	"time"
)

type result struct {
	Seq     int
	IsValid bool
}

var (
	total   = 0
	valid   = 0
	invalid = 0
	records [][]string
	printMu sync.Mutex
)

func ProcessLinks() error {
	var (
		wg        sync.WaitGroup
		startTime time.Time
		endTime   time.Time
	)
	startTime = time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fp, err := files.NewFileProcessor(3 * 1024 * 1024)
	if err != nil {
		return errors.Wrapf(err, "failed to create file processor")
	}
	records, err = fp.Reader.ReadAll()
	if err != nil {
		return errors.Wrapf(err, "failed to read data file")
	}

	total = len(records)
	resultCh := make(chan result, len(records))
	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(resultCh)
		done <- struct{}{}
	}()

	for index, record := range records {
		url := record[len(record)-1]
		wg.Add(1)
		go checkLink(ctx, url, index, resultCh, &wg)
	}

	index := 0 // 控制进度条打印进度
	for res := range resultCh {
		index++
		wg.Add(1)
		go printProgress(index, res, fp, &wg)
	}
	<-done
	endTime = time.Now()
	fmt.Printf("\nTotal time taken:%v\n", endTime.Sub(startTime))
	fp.FlushClose()
	close(done)
	return nil
}

func checkLink(ctx context.Context, url string, seq int, ch chan result, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Range", "bytes=0-0")
	resp, err := client.Do(req)
	isValid := err == nil && resp.StatusCode == http.StatusOK
	ch <- result{
		Seq:     seq,
		IsValid: isValid,
	}
	if err == nil {
		resp.Body.Close()
	}
}

func printProgress(i int, res result, fp *files.FileProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	if res.IsValid {
		printMu.Lock()
		valid++
		fp.WriteRow(records[res.Seq], true)
		printMu.Unlock()
	} else {
		printMu.Lock()
		invalid++
		fp.WriteRow(records[res.Seq], false)
		printMu.Unlock()
	}
	// 每 10 次以进度条的形式打印一下输出,打印时清除上次的信息
	if i%10 == 0 {
		printMu.Lock()
		progressPercent := float64(i) / float64(total) * 100
		successRate := float64(valid) / float64(total) * 100
		progressBar := strings.Repeat("=", int(progressPercent/2)) + ">" + strings.Repeat("·", 50-int(progressPercent/2))
		printMu.Unlock()
		fmt.Printf("\r[%s] %.2f%% success: %.2f%%", progressBar, progressPercent, successRate)
	}
}
