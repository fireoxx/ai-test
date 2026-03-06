---
description: 完整需求开发工作流：需求分析 → 技术文档（待确认）→ 确认后实施 → 完成总结
---

# 完整需求开发工作流 /dev

## 概述
端到端需求开发工作流，覆盖从需求分析到交付完成的全流程。
工作流分为四个阶段，**阶段二结束后必须等待用户确认，不得提前实施代码**。

---

## 阶段零：记忆召回（任务开始前）

在分析需求前，先查询本地 PG 记忆库，召回相关历史经验和错误模式：

```bash
PROJECT_NAME="ai-app"

psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
-- 历史相关任务
SELECT task_name, keywords, summary, problem_count, created_at
FROM ai_session_records
WHERE project_name = '$PROJECT_NAME'
  AND (keywords ILIKE '%{关键词}%' OR task_name ILIKE '%{关键词}%')
ORDER BY created_at DESC LIMIT 3;

-- 相关错误模式（重点关注！）
SELECT error_type, description, prevention, occurrence_count
FROM ai_error_patterns
WHERE project_name = '$PROJECT_NAME'
  AND keywords ILIKE '%{关键词}%'
ORDER BY occurrence_count DESC LIMIT 5;
SQL
```

如查到相关**错误模式**，在实施前主动告知用户，并在对应步骤重点注意。

---

## 阶段一：需求分析

1. 仔细理解用户描述的需求，提取以下关键信息：
   - 功能描述与业务目标
   - 涉及的子项目（`apps/server` 后端 / `apps/web` 管理后台 / `apps/H5Drift` H5端）
   - 是否涉及新增接口 / 改动已有接口 / 删除接口
   - 是否涉及数据库表结构变更（PostgreSQL）
   - 是否涉及 Redis / 定时任务（`task/`）

2. 探索相关代码，理解现有实现：
   - `apps/server/router/` — 路由注册
   - `apps/server/api/v1/` — Handler 层
   - `apps/server/service/` — 业务逻辑层
   - `apps/server/model/` — 数据库模型及请求/响应结构体
   - `apps/server/global/` — 全局变量（GVA_DB、GVA_LOG、GVA_CONFIG、GVA_REDIS）
   - `apps/web/src/view/` — 管理后台页面
   - `apps/H5Drift/src/views/` — H5 页面

3. 如需求存在歧义或缺漏，**先向用户提问明确后再继续**。

---

## 阶段二：输出技术文档（待确认模式）

调用 `tech-doc` 技能，输出标准技术方案文档。
保存路径：`.windsurf/技术方案/{需求名称}.md`

文档包含以下章节：
- **需求概述**：功能描述、核心目标、特殊约束
- **影响范围**：涉及文件清单（新增/改动）、涉及数据库表
- **数据库设计**：Model 结构体定义、AutoMigrate 注册（不生成独立 SQL 文件）
- **接口设计**：新增接口（路由+方法+Swagger注解+入参+出参）、改动接口、删除接口
- **核心实现思路**：关键逻辑说明、数据流、注意事项
- **风险点**：潜在风险及应对措施

**⚠️ 文档输出后告知用户路径，明确等待确认，不得进入实施阶段。**

---

## 阶段三：确认后制定实施计划并执行

用户确认技术方案后：

1. 使用 `todo_list` 工具拆解任务，建议按以下顺序：
   - 通过 MCP pgsql 确认数据库表结构
   - Model 定义（`apps/server/model/{module}/`）
   - AutoMigrate 注册
   - 请求/响应结构体（`model/{module}/request/` 和 `response/`）
   - Service 业务逻辑（`apps/server/service/{module}/`）
   - Handler 层（`apps/server/api/v1/{module}/`）
   - 路由注册（`apps/server/router/{module}/`）
   - 在各 `enter.go` 中注册新结构体
   - 前端接口封装 + 页面开发（如涉及）

2. 按任务逐步实施，每完成一项更新 todo 状态。

3. 严格遵循 `AGENTS.md` 及 `go-code-rule.md` 中的代码规范：
   - Handler 只做参数绑定/校验/调用service/返回响应
   - 使用 `response.OkWithData` / `response.FailWithMessage`
   - 日志使用 `global.GVA_LOG`，携带 `zap.Error(err)`
   - 每个 Handler 函数必须写完整 Swagger 注解

4. 实施过程中若发现与方案有偏差，**暂停并告知用户**后再继续。

---

## 阶段四：完成总结

实施完成后，调用 `/summary` 工作流：
- 生成变更总结文档
- 生成 Apifox 接口文档（调用 `sync-apifox` 技能）
- 告知用户需要手动执行的操作（配置等）
