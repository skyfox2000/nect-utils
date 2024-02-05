package async

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/skyfox2000/nect-utils/ants"
)

// Ants 对应的结构体
var Async = &asyncStruct{}

type asyncStruct struct{}

type ResultInfo struct {
	Key    string
	Result interface{}
	Err    error
}

func (p *asyncStruct) AsyncRun(
	execFunc func() (interface{}, error),
	logName string,
	concurrent int,
	timeout *int) (interface{}, error) {
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	// 默认15秒
	timerDuration := 15 * time.Second
	if timeout != nil && *timeout > 0 {
		timerDuration = time.Duration(*timeout) * time.Second
	}

	// 创建定时context
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timerDuration)
	defer cancel()

	// 异步调用 plugin.Execute，并获取结果
	done := make(chan bool) // 添加一个信号通道用来表示协程是否完成

	// 协程库调用
	ants.Ants.Submit(logName, func() {
		defer close(done) // 在结束时关闭done通道
		r, e := execFunc()
		resultChan <- r
		errChan <- e
	}, concurrent)

	// 普通调用
	// go func() {
	// 	defer close(done) // 在结束时关闭done通道
	// 	r, e := execFunc()
	// 	resultChan <- r
	// 	errChan <- e
	// }()

	// 等待结果或错误
	var finalResult interface{}
	var finalErr error
	select {
	// 等待协程执行完毕
	case <-done:
		// defer close(resultChan)
		// defer close(errChan)
		finalResult = <-resultChan
		finalErr = <-errChan
	case <-timeoutCtx.Done():
		// 如果超时，执行相应的处理逻辑
		finalErr = errors.New("执行超时，Timeout: " + timerDuration.String())
	}

	if finalErr != nil {
		return nil, finalErr
	}

	result := finalResult

	return result, nil
}

func (p *asyncStruct) ConcurrentRun(
	execFunc func(index int, dataRow interface{}) (string, interface{}, error),
	dataRows []interface{},
	batchSize, delay int, logName string, ignErr bool) map[string]interface{} {
	mutex := &sync.Mutex{}
	results := make(map[string]interface{})

	if len(dataRows) > 0 {
		brResults := make(map[string]ResultInfo)
		resultChan := make(chan map[string]ResultInfo, 1)

		for index, dataRow := range dataRows {
			ants.Ants.Submit(logName+".ConcurrentRun", func(i int, row interface{}) func() {
				return func() {
					// 提交任务到 ants 池
					ke, rs, ex := execFunc(i, row)
					resultInfo := ResultInfo{
						Key:    ke,
						Result: rs,
						Err:    ex,
					}
					mutex.Lock()
					brResults[ke] = resultInfo
					mutex.Unlock()
					if len(brResults) == len(dataRows) {
						resultChan <- brResults
					}
				}
			}(index, dataRow), batchSize)
		}

		resultList := <-resultChan
		for _, resultInfo := range resultList {
			if ignErr {
				if resultInfo.Err == nil {
					results[resultInfo.Key] = resultInfo.Result
				}
			} else {
				results[resultInfo.Key] = resultInfo.Result
			}
		}
	}

	return results
}
