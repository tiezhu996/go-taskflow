# BUG_REPRO

## Bug 是什么
worker 池在 context 被取消后没有停止：生产者和消费者的循环都不检查 ctx.Done()，继续投递并处理全部剩余任务。

## 如何触发
提交 100 个任务、4 个 worker，在第一个任务处理时取消 context，观察处理计数。

## 错误信息
`pool kept processing after cancel: calls=100`
