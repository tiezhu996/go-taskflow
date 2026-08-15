# BUG_REPRO

## Bug 是什么
查询不存在任务时，service.Get 包装 store 错误时用 `%v` 丢掉了错误链，导致 errors.Is(err, store.ErrNotFound) 为 false。

## 如何触发
对一个空 store 调用 Service.Get("missing")，用 errors.Is 判断哨兵错误。

## 错误信息
`errors.Is(err, ErrNotFound)=false, err=get task nope: task not found`
