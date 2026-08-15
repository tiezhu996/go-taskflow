# BUG_REPRO

## Bug 是什么
Store.IDs() 把内部保存插入顺序的切片直接返回给调用方，调用方拿到列表后原地排序会写穿共享底层数组，污染后续列表顺序。

## 如何触发
依次提交 b、a、c、d 四个任务，拿到 ID 列表后 sort.Strings 排序，再次调用 ListIDs() 观察顺序。

## 错误信息
`ListIDs after=[a b c] want [b a c] (returned slice aliases internal state)`
