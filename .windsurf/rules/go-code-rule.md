---
trigger: model_decision
description: fansnewgo 项目 Go 代码开发规范，涵盖 handler/router/types/model/cron 各层约定
globs: internal/**/*.go
---

# fansnewgo 项目代码规范

---

## 一、Handler 层规范

### 文件命名
- PascalCase，如 `ShopTransferOrder.go`、`WechatWorkLink.go`
- 禁止使用 `handlers_*` 前缀

### 必须包含编译期接口断言
```go
var _ ShopYeePayHandler = (*shopYeePayHandler)(nil)
```

### Context 与 request_id
日志、model、队列等操作必须携带 request_id，便于上线后日志追踪：
```go
// 标准方式（从 gin.Context 提取）
ctx := context.WithValue(context.Background(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c))
// 简写方式（后台任务）
ctx := middleware.ACtxRequestIDField(context.Background())

// model 操作携带 ctx
insuranceFreightLog := &model.InsuranceFreightLog{}
model.GetDB().WithContext(ctx).Where("refund_sn = ?", afterSaleInfo.RefundSn).First(&insuranceFreightLog)
```

### 数据库操作
- **新增逻辑必须使用 GORM，禁止使用 DAO 层**
- 事务优先使用闭包形式，**事务内禁止出现耗时较长的操作（如三方 HTTP 请求）**：
```go
if err := model.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 业务处理（禁止在此调用外部 HTTP 接口）
    return nil
}); err != nil {
    // 错误处理
}
```

### 错误处理
- **所有接口禁止将 err.Error() 系统错误直接返回给用户**，必须记录日志后返回业务错误码：
```go
if err != nil {
    log.Error("查询失败", logger.Any("err", err))
    response.Error(c, ecode.ServerErr)
    return
}
```
- 使用 `global` 包中的错误码：
```go
response.Error(c, ecode.InvalidParams)
```

### 日志
- **禁止使用 `fmt.Println` 等控制台打印**，必须使用 logger 记录到日志：
```go
log := logger.WithFields(middleware.CtxRequestIDField(ctx))
log.Info("描述", logger.Any("key", value))
log.Error("操作失败", logger.Any("err", err))
```

---

## 二、Router 层规范

### 文件命名
- PascalCase，如 `Yeepay.go`、`WechatWorkLink.go`
- 禁止使用 `routers_*` 前缀

### 路由分组选择
| 场景 | 使用变量 |
|-----|---------|
| 商户端接口 | `apiV1RouterFnsNew` |
| 总后台接口 | `apiV2RouterFnsNew` |
| 商户端+总后台公共接口 | `apiV3RouterFnsNew` |

### 路由注册方式
```go
// ✅ 正确：直接在 group 上注册完整路径
group.POST("/yee-pay/create", h.Create)
group.GET("/yee-pay/list", h.GetList)

// ❌ 禁止：创建无必要的二级分组
subGroup := group.Group("/yee-pay")
subGroup.POST("/create", h.Create)
```
> 仅当需要为一组路由绑定独立中间件时，才允许使用二级分组

---

## 三、Types 层规范

### 文件命名
- 新建入参/出参文件放在 `internal/types/` 目录
- 文件名：`{模块名}_types.go`，如 `adLink_types.go`

### 列表接口
- 请求结构体嵌入 `types.PageReq`
- 响应结构体包含 `total` 及分页信息：
```go
type ListXxxResp struct {
    List     []*ListXxxItem `json:"list"`
    Total    int64          `json:"total"`
    types.PageReq
}
```

### 时间字段处理
- 列表返回的时间字段格式化为 `2006-01-02 15:04:05`
- 不返回 `deleted_at`（删除时间）
- 历史已有字段无需补充处理

---

## 四、Model 层规范

### 新增字段位置
- 新增字段在结构体**最后一行**追加，禁止插入中间位置

### 时间处理
- **所有时间相关操作使用 `carbonNew` 包处理，注意时区**：
```go
now := carbonNew.Now()
dateStr := carbonNew.Now().ToDateString()         // "2026-01-01"
datetimeStr := carbonNew.Now().ToDateTimeString() // "2026-01-01 10:01:01"
```

### 软删除字段规范
数据库删除标记字段统一使用 `is_deleted` tinyint 类型：
```sql
-- 字段定义
`is_deleted` tinyint NOT NULL DEFAULT '2' COMMENT '1.删除 2.未删除'
```
```go
// 程序中查询未删除记录
db.Where("is_deleted = ?", 2)

// 软删除操作
db.Model(&record).Update("is_deleted", 1)
```

### 权限字段规范
**业务数据表凡涉及权限控制，必须包含以下字段**：
```sql
`company_id` int NOT NULL DEFAULT '0' COMMENT '公司ID',
`admin_id`   int NOT NULL DEFAULT '0' COMMENT '管理员ID',
```
- `company_id`：用于数据隔离，按公司过滤
- `admin_id`：用于操作归属，记录创建/操作人
- 查询时必须携带 `company_id` 条件，防止跨公司数据泄漏

---

## 五、Cron 层规范

### 多机防重复执行
**多机部署的定时任务必须加分布式锁**：
```go
lockKey := "TaskName"
lockValue := "lock-value"
lock.GetInstance().TryLock(lockKey, lockValue, 10*time.Second, func() {
    // 任务逻辑
})
```

### 后台长期任务必须加 recover
```go
func RunWithRecover(ctx context.Context, worker func()) {
    go func() {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("worker panic", zap.Any("err", err), middleware.CtxRequestIDField(ctx), logger.Any("stack", debug.Stack()))
            }
        }()
        worker()
    }()
}
```

### 大数据任务分页处理
- **禁止一次查询大量数据**，必须分页处理：
```go
// ✅ 正确：分批查询
var offset int
for {
    var records []model.Xxx
    db.Offset(offset).Limit(500).Find(&records)
    if len(records) == 0 {
        break
    }
    // 处理 records
    offset += 500
}

// ❌ 禁止：全量查询
db.Find(&allRecords)
```

---

## 六、Global 层规范

### 优先使用全局枚举
- 状态码、类型标识等全局定义**优先使用 `global` 目录中的枚举**
- **禁止在业务代码中硬编码魔法数字**

### 枚举定义规范
非 1/0、true/false 类型的枚举必须定义常量及描述方法：
```go
// 定义自定义类型及枚举值
type ErpSource int

const (
    ErpActionPay    ErpSource = 1 // 支付成功
    ErpActionCancel ErpSource = 2 // 取消
    ErpActionEdit   ErpSource = 3 // 修改
    ErpActionAfter  ErpSource = 4 // 售后
)

func (s ErpSource) Int() int {
    return int(s)
}

func (s ErpSource) String() string {
    ErpSourceDesc := map[ErpSource]string{
        ErpActionPay:    "支付成功",
        ErpActionCancel: "取消订单",
        ErpActionEdit:   "编辑订单",
        ErpActionAfter:  "售后操作",
    }
    if value, ok := ErpSourceDesc[s]; ok {
        return value
    }
    return ""
}

// 使用方式
if order.PushErpSource == global.ErpActionCancel.Int() {
    // 逻辑操作
}
resp.PushErpSourceDesc = global.ErpActionCancel.String()
```

---

## 七、SQL 脚本规范

- 所有数据库变更脚本放在 `scripts/sql/` 目录
- 文件名格式：`add_{功能描述}.sql`，如 `add_wechat_user_remark_switch.sql`
- 脚本内容包含完整注释，说明字段用途
- 建表 DDL 必须包含 `is_deleted` tinyint 软删除字段（默认值 2=未删除，1=删除）
- 涉及权限的表必须包含 `company_id` 和 `admin_id` 字段
- **新建表索引设计必须覆盖常用查询条件字段，编写 SQL 时检查是否能使用上索引**

---

## 八、SQL 查询安全规范

### 禁止 fmt 拼接 SQL
**严禁使用 `fmt.Sprintf` 等字符串拼接方式构造 WHERE 条件**，必须使用 `?` 占位符：
```go
// ✅ 正确：使用占位符
db.Where("company_id = ? AND status = ?", companyID, status)

// ❌ 驳回：fmt 拼接，存在 SQL 注入风险
db.Where(fmt.Sprintf("company_id = %d AND status = %d", companyID, status))
```

### API 参数拼接 SQL 的处理
从接口获取的参数拼接 SQL 前，必须先校验参数有效性，使用 `[]string` 追加条件再 `AND` 拼接：
```go
var conditions []string
var args []interface{}

if req.Status > 0 {
    conditions = append(conditions, "status = ?")
    args = append(args, req.Status)
}
if req.CompanyID > 0 {
    conditions = append(conditions, "company_id = ?")
    args = append(args, req.CompanyID)
}

if len(conditions) > 0 {
    db = db.Where(strings.Join(conditions, " AND "), args...)
}
```

### 禁止循环中查询数据库
**禁止在循环中单条查询数据库，必须批量查询后在内存处理**：
```go
// ✅ 正确：批量查询
var users []model.User
db.Where("id IN ?", userIDs).Find(&users)
userMap := make(map[int64]*model.User)
for _, u := range users {
    userMap[u.ID] = u
}

// ❌ 禁止：循环查库
for _, id := range userIDs {
    var user model.User
    db.First(&user, id) // 每次循环查一次数据库
}
```
- 更详细的规则必须遵守rule ->  **.windsurf/rules/specification.md**
---

## 九、协程使用规范

### 禁止裸 goroutine
**禁止使用裸 `go func(){}`，必须使用项目封装的安全协程**：
```go
// ✅ 正确：使用安全协程
tools.RunRecover(ctx, func() {
    // 业务逻辑
})

// ❌ 禁止
go func() {
    // 业务逻辑
}()
```

### 协程数量控制
使用 `errgroup` 限制并发数量，防止内存泄漏：
```go
mutex := &sync.Mutex{}
g, _ := errgroup.WithContext(ctx)
g.SetLimit(10) // 限制最大并发数
for _, item := range items {
    item := item
    g.Go(func() error {
        // 并发写 map 必须加锁
        mutex.Lock()
        resultMap[item.ID] = item
        mutex.Unlock()
        return nil
    })
}
_ = g.Wait()
```

### API 异步协程使用新 Context
API 请求需要异步协程处理时，**必须使用新的 Context**，不能复用 gin.Context（请求结束后会被回收）：
```go
// ✅ 正确：使用新 context
newCtx := context.WithValue(context.Background(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c))
tools.RunRecover(newCtx, func() {
    // 异步处理
})

// ❌ 禁止：直接传 gin.Context 或复用 request ctx
go func() {
    doSomething(c) // c 可能已失效
}()
```

---

## 十、安全与稳定性规范

### 空指针判断
- 所有可能为 nil 的指针在使用前必须进行判断：
```go
if user == nil {
    return
}
// 再使用 user.xxx
```

### 生产环境禁止 panic
- **生产代码必须避免 panic**，发生错误时必须返回 error，由调用方决定处理方式：
```go
// ✅ 正确：返回 error
func doSomething() error {
    if err := someOp(); err != nil {
        return fmt.Errorf("doSomething failed: %w", err)
    }
    return nil
}

// ❌ 禁止
func doSomething() {
    if err := someOp(); err != nil {
        panic(err)
    }
}
```

### 禁止使用 os.Exit
**禁止在业务代码中使用 `os.Exit()`**，程序异常应通过返回 error 处理。

### defer 释放资源
文件句柄、锁等资源必须使用 `defer` 释放：
```go
f, err := os.Open(filename)
if err != nil {
    return err
}
defer f.Close()

mu.Lock()
defer mu.Unlock()
```

### sync.Mutex 禁止拷贝
`sync.Mutex` 是有状态的锁，**禁止以值拷贝方式使用**，必须用指针传递：
```go
// ✅ 正确：指针传递
type Service struct {
    mu *sync.Mutex
}

// ❌ 禁止：值拷贝
func doWork(mu sync.Mutex) { // 拷贝了 Mutex
}
```

---

## 十一、三方接口调用规范

### 必须加重试机制
调用外部三方接口必须加重试逻辑，避免因网络抖动导致失败：
```go
var resp *SomeResp
var err error
for i := 0; i < 3; i++ {
    resp, err = callThirdPartyAPI(ctx, req)
    if err == nil {
        break
    }
    time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
}
if err != nil {
    log.Error("三方接口调用失败", logger.Any("err", err))
    return err
}
```

---

## 十二、代码质量规范

### 复杂逻辑必须加注释
逻辑较复杂、非显而易见的代码块必须添加注释说明意图：
```go
// 当日数据存储在 Redis 中，第二天才会落库
// 如果查询时间范围包含当天，需要从 Redis 合并数据后再排序
if includeToday {
    // ...
}
```

### 保持整体代码风格一致
- 避免重复造轮子，**优先复用项目内已有的工具函数和封装**
- 新增代码风格必须与所在文件/模块保持一致
- 审查代码时发现风格严重不一致的，驳回修改

### go.mod 变更说明
- `go.mod` 文件发生变化时，**必须说明新增/更新了哪个依赖包及原因**

---

## 十三、通用开发原则

1. **禁止修改与当前任务无关的文件**
2. 如需改动其他文件，必须明确说明原因并等待确认
3. 新增接口必须同步更新路由注册
4. 接口设计参考现有同类接口的实现方式
5. 上线 menu 表权限数据时，不得在线上页面操作，须提供 INSERT 语句（含主键 ID）给 DBA 执行

---

## 十四、需求开发流程规范（强制执行）

**所有功能需求，必须先输出技术方案，等用户明确确认后，才能开始写代码。禁止在用户确认前进行任何代码改动。**

### 流程步骤

1. **分析阶段**：阅读代码，理解现有结构
2. **技术方案**：以 Markdown 文档形式输出到 `.windsurf/技术方案/` 目录，并在对话中告知用户，内容包含：
   - 需求背景与目标
   - 涉及文件清单（新增/修改）
   - 核心逻辑说明（数据结构、接口、Redis 设计等）
   - 关键决策说明
3. **等待确认**：用户明确回复"确认"或"OK"等确认语句后，才开始实施
   - 方案中存在**不确定或需要决策的点**（如：某字段用哪个、某功能是否需要实现、补数如何处理等），**必须在方案中明确列出，并逐一得到用户确认后**，才能开始实施，不可自行假设或跳过
   - 确认完成后，**必须将确认结论更新到方案文档**，再开始实施
4. **实施阶段**：严格按已确认的方案逐步执行，不可超出方案范围，不可自行修改已确认的设计
5. **完成总结**：说明实际改动文件、接口、数据库变更，并将方案状态更新为"已完成"

### 违规处理
- 若未经确认直接写代码，须立即回滚所有改动，重新输出方案后等待确认，确认后要修改方案状态
- 若实施内容与已确认方案不符（如遗漏功能点、修改了 Redis key 格式等），须立即说明差异并修正

---

## 十五、队列处理器开发规范

新增一个队列处理器时，**必须同时完成以下三步，缺一不可**：

### 第一步：定义队列名常量
在 `internal/queue/base.go` 中新增队列名常量：
```go
FansSkipAdDay = Prefix + "FansSkipAdDay" // 广告跳转数统计队列
```

### 第二步：注册队列并发权重（asynq.Config）⚠️ 极易遗漏
在 `internal/queue/processorqueue/base.go` 的 `asynq.Config.Queues` map 中添加对应条目：
```go
queue.QName + "FansSkipAdDay": 6, // 广告跳转数统计队列
```
> **未在 `Queues` 中注册的队列，asynq Server 不会消费任务，导致任务静默堆积且无任何报错提示。**

### 第三步：绑定处理器
在 `internal/queue/processorqueue/base.go` 的 `StartServer()` 方法中通过 `mux.setHandle` 绑定处理器：
```go
mux.setHandle(queue.FansSkipAdDay, NewFansSkipAdDay(background)) // 广告跳转数统计队列
```

### 处理器文件规范
- 文件放在 `internal/queue/processorqueue/` 目录，文件名使用 PascalCase，如 `FansSkipAdDay.go`
- 必须实现 `SendTask` 和 `ProcessTask` 两个方法
- `ProcessTask` 中必须重新创建 Context 并携带 requestID：
```go
ctx := context.WithValue(context.Background(), middleware.ContextRequestIDKey, requestID)
```
- Payload 中的 `Params` 经 JSON 反序列化后是 `map[string]interface{}`，需二次 Marshal/Unmarshal 转为具体类型：
```go
paramsBytes, _ := json.Marshal(p.Params)
json.Unmarshal(paramsBytes, &targetStruct)
```