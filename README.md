# 客服工单系统（Ticketing）

一个纯 Go 标准库实现的客服工单管理 REST API 服务，采用标准 Go 工程目录结构，内存存储，零第三方依赖。

## 目录结构

```
origin/
├── cmd/server/          # 程序入口
├── internal/
│   ├── app/             # 依赖装配
│   ├── config/          # 配置加载
│   ├── model/           # 领域模型与校验
│   ├── store/           # 数据访问接口 + 内存实现
│   ├── service/         # 业务逻辑层
│   └── handler/         # HTTP 处理器层
└── pkg/
    ├── httpx/           # HTTP 响应工具
    ├── idgen/           # ID 生成
    └── logger/          # 分级日志
```

## 运行

```bash
go run ./cmd/server
PORT=8081 go run ./cmd/server
```

默认监听 `:8080`。

## 测试

```bash
go test ./...
```

## 工单状态机

`open → processing → resolved → closed`，`resolved/closed` 可重新 `open`。

## API 接口

### 分类

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/categories` | 创建分类 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/categories/{id}` | 分类详情 |
| PUT | `/api/categories/{id}` | 更新分类 |
| DELETE | `/api/categories/{id}` | 删除分类 |

### 客服

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/agents` | 创建客服 |
| GET | `/api/agents` | 客服列表 |
| GET | `/api/agents/{id}` | 客服详情 |
| PUT | `/api/agents/{id}` | 更新客服 |
| PATCH | `/api/agents/{id}/status` | 调整状态 `{"status":"busy"}` |
| DELETE | `/api/agents/{id}` | 删除客服 |

### 工单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/tickets` | 创建工单 |
| GET | `/api/tickets?status=&priority=&category_id=&assignee_id=&keyword=&page=&size=` | 工单列表 |
| GET | `/api/tickets/{id}` | 工单详情 |
| PUT | `/api/tickets/{id}` | 更新工单 |
| DELETE | `/api/tickets/{id}` | 删除工单 |
| PATCH | `/api/tickets/{id}/assign` | 分配客服 `{"assignee_id":"..."}` |
| PATCH | `/api/tickets/{id}/status` | 状态流转 `{"status":"resolved"}` |

### 留言与统计

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/tickets/{id}/comments` | 添加留言 |
| GET | `/api/tickets/{id}/comments` | 留言列表 |
| DELETE | `/api/comments/{id}` | 删除留言 |
| GET | `/api/stats` | 工单统计 |
| GET | `/api/stats/agents` | 客服工作量 |
