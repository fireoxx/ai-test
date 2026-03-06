---
description: 记忆管理工作流：查询历史经验、写入新记忆、查看错误模式库
---

# 记忆管理工作流 /memory

用于管理 AI 的本地 PostgreSQL 记忆库，支持跨项目隔离存储。

**数据库配置**：host=127.0.0.1 port=5432 user=root password=root dbname=test

---

## 项目名获取规则

每次执行记忆操作前，`PROJECT_NAME` 固定取工作区根目录名：

> 当前项目：`ai-app`（路径 `/Users/fengxiaobo/docker/ai-app`）

长期记忆若适用于所有项目，`project_name` 填 `__global__`

---

## 子命令说明

| 子命令 | 说明 |
|--------|------|
| `recall` | 查询历史经验（任务开始前执行） |
| `save` | 写入本次任务记忆（任务结束后执行） |
| `errors` | 查看错误模式库 |
| `learn` | 写入长期记忆（规范/经验/最佳实践） |
| `stats` | 查看记忆统计 |

---

## /memory recall — 任务开始前召回相关记忆

**执行前，我（AI）先从用户描述中提取 2~4 个关键词，替换下面 SQL 中的 `关键词` 变量后执行。**

```bash
# 将 KW 替换为实际关键词（如：wechatGroup、fans_count、排序）
KW="关键词"

psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
-- 1. 历史相关任务摘要
SELECT task_name, keywords, summary, problem_count, created_at
FROM ai_session_records
WHERE project_name = 'ai-app'
  AND (keywords ILIKE '%$KW%' OR task_name ILIKE '%$KW%')
ORDER BY created_at DESC LIMIT 5;

-- 2. 相关错误模式（重点关注，避免重蹈覆辙）
SELECT error_type, error_category, description, prevention, occurrence_count
FROM ai_error_patterns
WHERE project_name = 'ai-app'
  AND keywords ILIKE '%$KW%'
ORDER BY occurrence_count DESC LIMIT 10;

-- 3. 相关对话历史（按关键词搜索对话摘要）
SELECT c.session_id, c.task_name, c.turn_index, c.role, c.content_summary, c.created_at
FROM ai_conversation_logs c
WHERE c.project_name = 'ai-app'
  AND (c.keywords ILIKE '%$KW%' OR c.content_summary ILIKE '%$KW%')
ORDER BY c.created_at DESC LIMIT 8;

-- 4. 全局+本项目长期记忆（更新 recall_count）
UPDATE ai_long_term_memory
SET recall_count = recall_count + 1, last_recalled_at = NOW()
WHERE project_name IN ('ai-app', '__global__')
  AND (tags ILIKE '%$KW%' OR title ILIKE '%$KW%');

SELECT category, title, content, tags, importance
FROM ai_long_term_memory
WHERE project_name IN ('ai-app', '__global__')
  AND (tags ILIKE '%$KW%' OR title ILIKE '%$KW%')
ORDER BY importance DESC, recall_count DESC LIMIT 5;
SQL
```

**召回后行动**：
- 查到**错误模式**：任务开始前主动告知用户，实施时重点注意
- 查到**历史方案**：参考已有实现，避免重复设计

---

## /memory save — 任务结束后写入记忆

**执行方式：AI 整理内容后生成本地脚本，再静默执行，对话中不展示大量 SQL，减少 token 消耗。**

### 执行流程（AI 操作步骤）

1. 根据本次对话，整理以下内容（仅在 AI 内部组织，不输出到对话）：
   - `task_name`、`task_type`、`keywords`（2~5个）、`summary`（100字内）
   - `files_changed`、`api_changed`、`db_changed`
   - `problem_count` 及各问题的 `description/root_cause/solution/prevention`
   - 各对话阶段摘要（turn 0~4）

2. 使用 `write_to_file` 工具生成脚本到 `.windsurf/memory_scripts/ai_memory_save_{SESSION_ID}.sh`

3. 使用 `run_command` 执行脚本，对话中只回复：
   ```
   记忆已写入，session_id=ai-app_20260305_150023
   ```

### 脚本模板（AI 填充后写入 .windsurf/memory_scripts/ai_memory_save_{SESSION_ID}.sh）

```bash
#!/bin/bash
set -e
PROJECT_NAME="ai-app"
SESSION_ID="${PROJECT_NAME}_$(date '+%Y%m%d_%H%M%S')"
TASK_NAME="【AI填充】任务名称"

psql -h 127.0.0.1 -p 5432 -U root -d test << SQL

-- ① 任务摘要
INSERT INTO ai_session_records (
    session_id, project_name, task_name, task_type,
    keywords, summary, key_decisions,
    files_changed, api_changed, db_changed, problem_count
) VALUES (
    '$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME',
    '【AI填充】feature/bugfix/refactor/config',
    '【AI填充】关键词1,关键词2,关键词3',
    '【AI填充】100字内摘要',
    '【AI填充】["决策1"]',
    '【AI填充】["handler/Xxx.go"]',
    '【AI填充】[{"method":"POST","path":"/xxx","desc":"xxx"}]',
    '【AI填充】[{"table":"xxx","op":"新增","desc":"xxx"}]',
    0
);

-- ② 对话轮次（turn 0=用户需求 1=方案 2=实施 3=问题 4=总结）
INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME', 0, 'user', '【AI填充】用户需求摘要', '【AI填充】关键词');

INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME', 1, 'assistant', '【AI填充】方案摘要', '【AI填充】关键词');

INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME', 2, 'assistant', '【AI填充】实施摘要', '【AI填充】关键词');

-- turn 3 有问题时取消注释
-- INSERT INTO ai_conversation_logs ... 3, 'assistant', '问题→原因→修正', '关键词';

INSERT INTO ai_conversation_logs (session_id, project_name, task_name, turn_index, role, content_summary, keywords)
VALUES ('$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME', 4, 'assistant', '【AI填充】总结摘要', '【AI填充】关键词');

-- ③ 错误记录（有问题时取消注释）
-- INSERT INTO ai_error_patterns (session_id, project_name, task_name, error_type, error_category,
--     description, root_cause, solution, prevention, keywords)
-- VALUES ('$SESSION_ID', '$PROJECT_NAME', '$TASK_NAME',
--     '【类型】', '【分类】', '【描述】', '【原因】', '【方案】', '【预防】', '【关键词】');
--
-- UPDATE ai_error_patterns SET occurrence_count = occurrence_count + 1, updated_at = NOW()
-- WHERE project_name = '$PROJECT_NAME' AND error_category = '【分类】'
--   AND keywords ILIKE '%【关键词】%' AND session_id != '$SESSION_ID';

SQL

echo "session_id=$SESSION_ID"
```

---

## /memory errors — 查看错误模式库

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
-- 本项目高频错误（按出现次数排序）
SELECT error_type, error_category, description, occurrence_count, prevention
FROM ai_error_patterns
WHERE project_name = 'ai-app'
ORDER BY occurrence_count DESC LIMIT 20;

-- 规范类错误（最需要注意）
SELECT description, root_cause, prevention, occurrence_count
FROM ai_error_patterns
WHERE project_name = 'ai-app' AND error_category = 'spec'
ORDER BY occurrence_count DESC;
SQL
```

---

## /memory learn — 写入长期记忆

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
-- 项目特有规范用 ai-app，全局通用规范用 __global__
INSERT INTO ai_long_term_memory (
    project_name, category, title, content, tags, importance
) VALUES (
    'ai-app',
    'rule',
    '规范标题',
    '规范内容，详细说明',
    '标签1,标签2,Go,handler',
    3
);
SQL
```

**category**：`rule`（规范）/ `pattern`（设计模式）/ `solution`（解决方案）/ `spec`（规格说明）

**importance**：1（低）~ 5（关键必须遵守）

---

## /memory stats — 统计概览

```bash
psql -h 127.0.0.1 -p 5432 -U root -d test << SQL
SELECT '对话留存' AS type, COUNT(*) AS total, MAX(created_at) AS last_updated
FROM ai_session_records WHERE project_name = 'ai-app'
UNION ALL
SELECT '对话轮次', COUNT(*), MAX(created_at) FROM ai_conversation_logs WHERE project_name = 'ai-app'
UNION ALL
SELECT '长期记忆', COUNT(*), MAX(updated_at) FROM ai_long_term_memory WHERE project_name IN ('ai-app','__global__')
UNION ALL
SELECT '错误模式', COUNT(*), MAX(updated_at) FROM ai_error_patterns WHERE project_name = 'ai-app';

SELECT error_type, error_category, COUNT(*) AS cnt, SUM(occurrence_count) AS total_occ
FROM ai_error_patterns WHERE project_name = 'ai-app'
GROUP BY error_type, error_category ORDER BY total_occ DESC;
SQL
```

---

## 查询完整对话记录

```bash
# 查某次任务完整对话
psql -h 127.0.0.1 -p 5432 -U root -d test -c "
SELECT turn_index, role, content_summary, keywords
FROM ai_conversation_logs
WHERE session_id = 'ai-app_20260305_150023'
ORDER BY turn_index ASC;"

# 跨任务搜索对话内容
psql -h 127.0.0.1 -p 5432 -U root -d test -c "
SELECT s.session_id, s.task_name, s.created_at, l.turn_index, l.role, l.content_summary
FROM ai_session_records s
JOIN ai_conversation_logs l ON s.session_id = l.session_id
WHERE s.project_name = 'ai-app'
  AND (l.content_summary ILIKE '%关键词%' OR l.keywords ILIKE '%关键词%')
ORDER BY s.created_at DESC, l.turn_index ASC;"
```

---

## 表结构速查

| 表名 | 用途 | 核心字段 |
|------|------|---------|
| `ai_session_records` | 任务摘要（短期记忆） | project_name, task_name, keywords, summary, problem_count |
| `ai_conversation_logs` | 完整对话轮次记录 | session_id, project_name, turn_index, role, content_summary |
| `ai_long_term_memory` | 长期记忆（规范经验） | project_name, category, title, content, importance, recall_count |
| `ai_error_patterns` | 错误模式库（自学习） | project_name, error_type, description, prevention, occurrence_count |
