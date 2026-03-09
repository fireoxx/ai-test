---
description: 实施工作流：根据已确认的技术方案，按计划逐步实施代码变更
---

# 实施工作流 /impl

在用户确认技术方案文档后执行此工作流。

---

## 步骤 0：记忆召回（先于一切执行）

通过 MCP pgsql 查询历史相关错误模式，避免重复踩坑：

```sql
SELECT error_type, description, prevention, occurrence_count
FROM ai_error_patterns
WHERE project_name = 'ai-app'
  AND keywords ILIKE '%{关键词}%'
ORDER BY occurrence_count DESC LIMIT 5;
```

如查到相关**错误模式**，实施对应步骤时重点注意。

---

## 步骤 1：读取技术方案

- 读取 `.windsurf/技术方案/` 目录下对应的方案文档
- 确认实施范围（文件清单、接口列表、数据库改动）
- 如需澄清，先提问再实施

---

## 步骤 2：创建实施计划

使用 `todo_list` 工具，按以下顺序拆解任务：

```
[ ] 1. MCP pgsql 确认数据库表结构
[ ] 2. Model 层（apps/server/model/{module}/）
[ ] 3. AutoMigrate 注册（initialize/gorm.go）
[ ] 4. 请求/响应结构体（model/{module}/request/ 和 response/）
[ ] 5. Service 层（apps/server/service/{module}/）
[ ] 6. Handler 层（apps/server/api/v1/{module}/）
[ ] 7. Router 层（apps/server/router/{module}/）
[ ] 8. enter.go 注册（api/v1/enter.go、service/enter.go、router/enter.go）
[ ] 9. 前端接口封装 + 页面开发（如涉及）
```

---

## 步骤 3：逐步实施

按任务顺序实施，每完成一项更新 todo 状态为 completed。

### 通用规范检查（每个文件实施前确认）

**Handler 规范**（`apps/server/api/v1/{module}/`）
- 文件名 `snake_case.go`，如 `sys_user.go`
- 只做参数绑定/校验/调用 service/返回响应，禁止直接操作数据库
- 参数绑定：`c.ShouldBindJSON` / `c.ShouldBindQuery`
- 成功响应：`response.OkWithData(data, c)` / `response.OkWithMessage(msg, c)`
- 失败响应：`response.FailWithMessage(err.Error(), c)`
- 日志：`global.GVA_LOG.Error("失败", zap.Error(err))`
- 每个函数必须写完整 Swagger 注解（@Tags / @Summary / @Param / @Success / @Router）

**Service 规范**（`apps/server/service/{module}/`）
- 所有业务逻辑放 service 层
- 数据库操作：`global.GVA_DB`，事务：`global.GVA_DB.Transaction(...)`
- 日志：`global.GVA_LOG.Error/Info/Warn`
- `gorm.ErrRecordNotFound` 单独处理，不视为系统错误

**Router 规范**（`apps/server/router/{module}/`）
- 需要鉴权的路由挂在带有 `JWTAuth()` 中间件的路由组下
- 路由命名使用 camelCase，路径使用 lowerCamelCase
- 注册后在 `router/enter.go` 中引用

**Model 规范**（`apps/server/model/{module}/`）
- 嵌入 `global.GVA_MODEL` 获取 ID/CreatedAt/UpdatedAt
- 字段必须带完整 tag：`json` + `gorm` + 可选 `binding`
- 新增字段追加到结构体**最后一行**，禁止插入中间
- 请求结构体放 `request/`，响应结构体放 `response/`

**数据库规范**
- 使用 MCP pgsql 工具操作 PostgreSQL，不生成独立 SQL 文件
- 表结构变更通过 GORM AutoMigrate
- 所有字段 NOT NULL DEFAULT，text/jsonb 除外；禁止外键约束；禁止触发器
- 常用查询条件字段在 GORM tag 加 `index` 标注
- 禁止循环中单条查库，必须批量查询后在内存处理
- `gorm.ErrRecordNotFound` 单独处理，不视为系统错误，不打 Error 日志

---

## 步骤 4：自检

实施完成后逐项检查：

- [ ] 所有新增 Handler 是否注册了路由
- [ ] 路由文件是否在 `router/enter.go` 中被引用
- [ ] `api/v1/enter.go`、`service/enter.go` 是否注册了新结构体
- [ ] 如有新路由组，`initialize/router.go` 是否已挂载
- [ ] **GET 接口入参 tag 是否用 `form`，POST 接口入参 tag 是否用 `json`**
- [ ] 列表接口是否返回了 `total`、`page`、`pageSize`
- [ ] 是否有未处理的 error 返回
- [ ] 每个 Handler 函数是否有完整 Swagger 注解（@Tags/@Summary/@Security/@Param/@Success/@Router）
- [ ] 是否通过 MCP pgsql 确认了数据库表结构正确
- [ ] AutoMigrate 是否已注册新表
- [ ] `gorm.ErrRecordNotFound` 是否单独处理，没有错误当成系统错误打 Error 日志
- [ ] 是否有循环查库的情况（应改为批量查询）
- [ ] 异步协程是否使用了新的 Context，未复用 gin.Context

---

## 步骤 5：实施完成，调用 /summary 工作流
