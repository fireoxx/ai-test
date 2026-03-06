---
description: 完成总结工作流：生成接口变更、数据库改动、文件清单的完整总结，并输出 Apifox 接口文档
---

# 完成总结工作流 /summary

在实施完成后执行此工作流，生成变更总结和接口文档。

---

## 步骤 1：梳理本次变更内容

回顾本次实施的所有改动，整理以下信息：

**接口变更**
- 新增接口：路由路径、HTTP 方法、功能描述、路由分组（鉴权/公开）
- 改动接口：接口路径、改动点（新增参数/新增返回字段/逻辑变更）
- 删除接口：接口路径、原功能描述

**数据库改动**
- 新增表：表名、用途、Model 结构体定义
- 新增字段：表名、字段名、添加位置（Model 结构体末尾）
- 索引变更：新增或删除的索引

**文件清单**
- 新增文件列表（含路径）
- 改动文件列表（含改动概述）

---

## 步骤 2：生成总结文档

保存路径：`.windsurf/需求/{需求名称}-实现说明.md`

文档结构如下：

```markdown
# {需求名称} 实现说明

## 一、接口变更

### 新增接口
| 路由 | 方法 | 功能 | 路由分组 |
|-----|------|------|---------|
| /api/xxx/list | GET | xxx列表 | 鉴权 |

### 改动接口
| 路由 | 方法 | 改动点 |
|-----|------|-------|
| /admin/xxx/save | POST | 新增字段 xxx |

### 删除接口
| 路由 | 方法 | 原功能 |
|-----|------|-------|

## 二、数据库改动

### 新增表
（含 Model 结构体定义）

### 字段变更
（含 Model 结构体新增字段）

## 三、文件改动清单

### 新增文件
- `apps/server/api/v1/{module}/xxx.go`
- `apps/server/service/{module}/xxx.go`
- `apps/server/model/{module}/xxx.go`
- `apps/server/router/{module}/xxx.go`

### 改动文件
- `apps/server/api/v1/enter.go` — 注册新 API 结构体
- `apps/server/service/enter.go` — 注册新 Service 结构体
- `apps/server/router/enter.go` — 注册新 Router

## 四、部署注意事项
- AutoMigrate 是否已注册新表结构
- 需更新配置项：（如有）
- 其他注意事项：（如有）
```

---

## 步骤 3：调用 sync-apifox 技能

为本次**新增和改动**的接口生成 Apifox 导入文件。
保存路径：`.windsurf/需求/{需求名称}api.Apifox.json`

---

## 步骤 4：问题汇总与记忆落地（写入 PostgreSQL）暂时不需要主动写，用户自己手动执行

**project_name 自动从工作区根目录名获取**，session_id 格式：`{项目名}_{日期时间戳}`

```bash
PROJECT_NAME="ai-app"
SESSION_ID="${PROJECT_NAME}_$(date '+%Y%m%d_%H%M%S')"
echo "session_id=$SESSION_ID"
```

### 4.1 写入任务摘要（ai_session_records）

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
INSERT INTO ai_session_records (
    session_id, project_name, task_name, task_type,
    keywords, summary, key_decisions,
    files_changed, api_changed, db_changed, problem_count
) VALUES (
    '$SESSION_ID',
    '$PROJECT_NAME',
    '{需求名称}',
    '{feature/bugfix/refactor/config}',
    '{关键词1,关键词2,关键词3}',
    '{任务摘要100字内：核心方案+关键词拼接}',
    '["决策1","决策2"]',
    '["apps/server/api/v1/{module}/xxx.go","apps/server/service/{module}/xxx.go"]',
    '[{"method":"POST","path":"/api/xxx","desc":"xxx接口"}]',
    '[{"table":"xxx","op":"新增表","desc":"xxx功能表"}]',
    {问题数量}
);
SQL
```

**摘要压缩原则**：用关键词代替长句，如 `wechatGroup+fans_count排序+GORM+wechat_group_fans_num_day`

### 4.2 写入对话轮次记录（ai_conversation_logs）

按关键阶段写多条记录，还原完整对话脉络：

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
-- turn 0：用户原始需求
INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '{需求名称}', 0, 'user',
    '{用户需求摘要：核心意图，100字内}', '{关键词1,关键词2}');

-- turn 1：方案分析
INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '{需求名称}', 1, 'assistant',
    '{方案设计：涉及文件X/Y，数据库表A，接口路由B，核心逻辑C}', '{方案关键词}');

-- turn 2：实施过程
INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '{需求名称}', 2, 'assistant',
    '{实施：新增handler方法X；types结构体Y；路由注册Z}', '{实施关键词}');

-- turn 3：遇到问题（有问题时填，无则跳过）
-- INSERT INTO ai_conversation_logs ... turn_index=3, '问题：[描述] → 原因：[原因] → 修正：[方案]'

-- turn 4：任务总结
INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '{需求名称}', 4, 'assistant',
    '{总结：完成功能X，变更文件N个，接口M个，问题P个}', '{总结关键词}');
SQL
```

### 4.3 写入错误模式库（ai_error_patterns，有问题时执行）

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
INSERT INTO ai_error_patterns (
    session_id, project_name, task_name,
    error_type, error_category,
    description, root_cause, solution, prevention, keywords
) VALUES (
    '$SESSION_ID', '$PROJECT_NAME', '{需求名称}',
    '{logic/spec/perf/compile}', '{code/spec/tech-doc}',
    '{错误描述}', '{根本原因}', '{解决方案}', '{预防措施}',
    '{关键词1,关键词2}'
);
SQL
```

### 4.4 生成问题汇总文档（按条件触发）

满足以下条件时生成 `.windsurf/问题汇总/{需求名称}-问题汇总.md`：
- 本次问题数 **≥ 3**
- 遇到规范类错误（spec 类型）
- 用户明确要求

---

## 步骤 5：告知用户

输出完成后告知用户：
1. 总结文档路径
2. Apifox 文档路径（如有接口变更）
3. **需要手动执行的操作**（SQL 脚本、配置修改等），明确列出
4. 是否有后续待跟进事项
5. 本次任务遇到的问题数量及是否生成了问题汇总文档