# BUG_REPRO

## Bug 是什么
Store 构造时没有初始化 tasks map，第一次 Put 写入 nil map 触发 panic。

## 如何触发
创建 store 后调用 service.Submit（内部走 store.Put），或运行 `go test ./internal/store -run TestPutGetList`。

## 错误信息
`panic: assignment to entry in nil map`（栈顶指向 Store.Put）
