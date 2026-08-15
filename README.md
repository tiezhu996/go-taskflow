# taskflow

一个用 Go 写的内存任务队列示例服务，演示 `handler → service → repository → model` 分层与并发 worker 池的常见写法。

## 功能

- 提交 / 查询任务，按提交顺序返回任务 ID 列表
- 并发 worker 池批量处理任务，支持 context 取消
- 内存存储，读写锁保护

## 目录结构

```
cmd/taskflow/        程序入口
internal/config/     环境变量配置
internal/model/      模型定义与纯工具函数
internal/store/      内存存储（任务 map + 插入顺序 + 完成计数）
internal/service/    业务逻辑（提交 / 查询 / 处理）
internal/worker/     并发 worker 池
```

## 运行与测试

```bash
go build ./...          # 编译
go test ./...           # 全量测试
go run ./cmd/taskflow   # 启动
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `TASKFLOW_WORKERS` | worker 数量 | `2` |

## 技术栈

- Go 1.22
