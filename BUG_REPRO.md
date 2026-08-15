# BUG_REPRO

## Bug 是什么
并发处理任务时，处理完成计数器在多 worker 并发自增时没有加锁，发生数据竞争并丢失更新，最终 processed 计数小于实际处理的任务数。

## 如何触发
提交 500 个任务、用 16 个 worker 并发处理，处理结束后读取 Store.Processed()；或运行评分测试 `go test -race ./internal/worker -run TestVerify -count=20 .`。

## 错误信息
- `go test -race` 报 data race（Store.processed 被并发读写）。
- 计数断言失败：`processed=499 want 500`（实际数字每次略有不同，通常小于期望值）。
