---
description: 技术方案工作流：分析需求并生成技术文档，等待用户确认后再实施，不做任何代码改动
---

# 技术方案工作流 /plan

只做分析和文档输出，**不改动任何代码**。输出文档后等待用户确认。

---

## 步骤 0：记忆召回（先于一切执行）

通过 MCP pgsql 查询历史相关任务和错误模式：

```sql
-- 历史相关任务
SELECT task_name, keywords, summary, problem_count, created_at
FROM ai_session_records
WHERE project_name = 'ai-app'
  AND (keywords ILIKE '%{关键词}%' OR task_name ILIKE '%{关键词}%')
ORDER BY created_at DESC LIMIT 3;

-- 相关错误模式（重点关注）
SELECT error_type, description, prevention, occurrence_count
FROM ai_error_patterns
WHERE project_name = 'ai-app'
  AND keywords ILIKE '%{关键词}%'
ORDER BY occurrence_count DESC LIMIT 5;
```

如查到相关**错误模式**，在方案中主动标注，实施时重点注意。

---

## 步骤 1：理解需求

- 仔细阅读用户的需求描述
- 如果需求存在歧义或缺失关键信息，先提问明确，再继续

---

## 步骤 2：探索代码库（并行读取）

重点查找：

- 相关 Handler（`apps/server/api/v1/`）
- 相关路由注册（`apps/server/router/`）
- 相关 Service（`apps/server/service/`）
- 相关模型（`apps/server/model/`）
- 全局变量（`apps/server/global/`）
- 前端页面（`apps/web/src/view/`、`apps/H5Drift/src/views/`）

目标：充分理解现有代码结构和实现方式，确保方案与项目规范一致。

---

## 步骤 3：分析影响范围

梳理以下内容：

**文件层面**
- 需要新增的文件（handler/types/router/model/SQL）
- 需要改动的文件（改动点描述）

**接口层面**
- 新增接口：路由路径、HTTP 方法、所属分组（商户端/总后台/公共）
- 改动接口：改动参数或返回字段
- 删除接口：废弃的接口

**数据库层面**
- 新增数据库表（Model 结构体定义 + AutoMigrate 注册，不生成独立 SQL 文件）
- 新增或修改字段（在 Model 结构体末尾追加）
- 需要索引的字段（在 GORM tag 中标注 `index` 或 `index:idx_name`）
- 通过 MCP pgsql 确认现有表结构（`\d table_name`）

---

## 步骤 4：调用 tech-doc 技能

调用 `tech-doc` 技能，生成标准格式技术方案文档。
保存路径：`.windsurf/技术方案/{需求名称}.md`

---

## 步骤 5：等待确认

输出文档后，告知用户：
1. 文档保存路径
2. 需要特别关注的设计决策或风险点
3. 明确说明：**等待您确认后再实施，当前不会改动任何代码**
