package jsrun

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/dop251/goja"
	utils "github.com/skyfox2000/nect-utils"
	"github.com/skyfox2000/nect-utils/ants"
	"github.com/skyfox2000/nect-utils/async"
	"github.com/skyfox2000/nect-utils/jsmodule"
	"github.com/skyfox2000/nect-utils/json"
	"github.com/skyfox2000/nect-utils/underscore"
)

// JSRun 对应的结构体
var JSRun = &jsrunStruct{}
var vm *goja.Runtime

type jsrunStruct struct {
}

func (p *jsrunStruct) Compile(jsName, jscodeStr string, utilsTool utils.UtilsTool) (*goja.Program, error) {
	jsCode := jscodeStr
	if vm == nil {
		vm = goja.New()
	}

	prepareCode := fmt.Sprintf("\"use strict\";\n(async function(){\n%s\n})();", jsCode)

	if underscore.Underscore.Contains(utilsTool.Debug, "Debug") {
		utilsTool.Logger.Debug("["+utilsTool.Name+"] ", jsName, ", jsCode: \n", prepareCode)
	}

	prog, err := goja.Compile(jsName+".js", prepareCode, false)
	if err != nil {
		lines := strings.Split(prepareCode, "\n")
		for index, line := range lines {
			lines[index] = "line " + strconv.Itoa(index+1) + ":  " + line
		}
		utilsTool.Logger.Error("["+utilsTool.Name+"] error compiling script: \n", strings.Join(lines, "\n"))
		return nil, err
	}

	return prog, nil
}

func (p *jsrunStruct) Run(
	ctx *context.Context,
	prog *goja.Program,
	utilsTool utils.UtilsTool,
	data map[string]interface{},
	keyMutexes map[string]*sync.RWMutex,
	cacheFlag bool,
	concurrent int,
	timeout *int) (interface{}, error) {

	newVm := goja.New()
	newVm.Set("require", func(call goja.FunctionCall) goja.Value {
		return p.require(call, utilsTool)
	})

	for k, v := range data {
		newVm.Set("$"+k, v)
	}

	// 注册 console 对象
	console := map[string]func(args ...interface{}){
		"log": func(args ...interface{}) {
			p.consoleLog("log", utilsTool, args...)
		},
		"info": func(args ...interface{}) {
			p.consoleLog("info", utilsTool, args...)
		},
		"warn": func(args ...interface{}) {
			p.consoleLog("warn", utilsTool, args...)
		},
		"debug": func(args ...interface{}) {
			p.consoleLog("debug", utilsTool, args...)
		},
		"error": func(args ...interface{}) {
			p.consoleLog("error", utilsTool, args...)
		},
	}
	newVm.Set("console", console)

	for _, mutex := range keyMutexes {
		mutex.Lock()
	}

	result, ex := async.Async.AsyncRun(func() (interface{}, error) {
		r, e := newVm.RunProgram(prog)
		return r, e
	}, utilsTool.Name, concurrent, timeout)

	for _, mutex := range keyMutexes {
		mutex.Unlock()
	}

	// 等待结果或错误
	finalResult := result.(goja.Value)
	finalErr := ex

	for k := range data {
		newVm.GlobalObject().Delete("$" + k)
	}

	if finalErr != nil {
		return nil, finalErr
	}

	// 判断是否是 Promise
	if promise, ok := finalResult.Export().(*goja.Promise); ok {
		next := make(chan bool)
		var promiseResult interface{}
		chanResult := make(chan interface{}, 1)
		ants.Ants.Submit("JSRun.Promise", func() {
			defer close(next)
			// 如果是 Promise，等待异步操作完成，并获取结果
			var pr interface{}
			pr = promise.Result()

			if ok := promise.Result() == nil; ok {
				pr = nil
			} else if errStr, ok := p.isError(promise.Result()); ok {
				pr = utils.NewError(3001, errStr)
			} else {
				pr = promise.Result().Export()
			}
			chanResult <- pr
		}, concurrent)
		<-next

		promiseResult = <-chanResult
		if errResult, ok := promiseResult.(*utils.CustomError); ok {
			return nil, errResult
		}

		if errResult, ok := promiseResult.(map[string]interface{}); ok {
			errno, ok1 := errResult["errno"].(int64)
			msg, ok2 := errResult["msg"].(string)
			if ok1 && ok2 {
				err := utils.NewError(int(errno), msg)
				return nil, err
			}
		}

		// JS不能修改系统数据

		return promiseResult, nil
	}

	return nil, nil
}

func (p *jsrunStruct) isError(result goja.Value) (string, bool) {
	errStr := strings.Split(result.String(), ":")
	if strings.Index(errStr[0], "Error") > 0 {
		return result.String(), true
	}
	return "", false
}

func (p *jsrunStruct) consoleLog(logLevel string, utilsTool utils.UtilsTool, args ...interface{}) {
	var message string
	message += "[" + utilsTool.Name + "] "

	for _, arg := range args {
		switch t := arg.(type) {
		case string:
			message += t
		default:
			jsonStr := json.JSON.Log(t)
			if message != "" {
				message += " "
			}
			message += jsonStr.(string)
		}
	}

	switch logLevel {
	case "log":
		utilsTool.Logger.Info(message)
	case "info":
		utilsTool.Logger.Info(message)
	case "warn":
		utilsTool.Logger.Warn(message)
	case "debug":
		utilsTool.Logger.Debug(message)
	case "error":
		utilsTool.Logger.Error(message)
	}
}

// 模块加载
func (p *jsrunStruct) require(call goja.FunctionCall, utilsTool utils.UtilsTool) goja.Value {
	moduleName := call.Argument(0).String()

	modules := jsmodule.JSModules
	module, ok := modules[moduleName]
	if !ok {
		err := utils.NewError(3001, "require module not found: "+moduleName)
		utilsTool.Logger.Error(err.Error())
		return vm.ToValue(err)
	}

	return vm.ToValue(module)
}
