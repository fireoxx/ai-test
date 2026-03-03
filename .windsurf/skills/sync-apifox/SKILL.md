---
name: sync-apifox
description: 扫描 Go 项目中的 API 接口并同步到 API Fox。当用户需要将项目 API 同步到 API Fox、更新 API 文档、或导出 OpenAPI 格式文档时使用此 Skill。触发词包括：同步 API、导出接口、更新 API Fox、生成 API 文档。
---

# API Fox 同步

扫描项目中带有 Swagger 注解的 Go 文件，生成 OpenAPI 3.0 格式文档，并同步到 API Fox。

## 配置

在执行前确认以下配置：
- **API Fox Token**: 从 API Fox 账号设置 -> API 访问令牌获取
- **项目 ID**: 从 API Fox 项目设置 -> 基本设置获取

当前配置：
- Token: `afxp_c14a7eZ0biqoZRHG01xGYMvTMpSNnqTaDfHb`
- 项目 ID: `7868592`

> 脚本启动时自动从本文件读取以上配置，无需修改脚本源码。

## 执行方式

**必须使用项目 venv 环境执行**，不可使用系统 python3（缺少 requests 模块）：

```bash
# 在 monorepo 根目录执行
/Users/fengxiaobo/docker/ai-app/.venv/bin/python3 .windsurf/skills/sync-apifox/scripts/sync_apifox.py
```

> ⚠️ 禁止使用 `python3 - <<'EOF'` here-doc 写法，会在 IDE run_command 中卡住。  
> 如需临时分析脚本，先写入 `/tmp/xxx.py` 再用 venv python3 执行。

## 执行流程

### 1. 扫描范围

脚本自动扫描 `apps/server/` 下以下目录：

| 类型 | 扫描目录 |
|------|---------|
| API Handler | `api/v1/` 下所有子目录（自动遍历，含 system、driftbottle 等） |
| Request 模型 | `model/*/request/` |
| Response 模型 | `model/*/response/`、`model/*/` |

新增业务模块时**无需手动配置**，只要在 `api/v1/<module>/` 下有 Swagger 注解即可自动识别。

### 2. 运行脚本

脚本会：
1. 扫描所有 API Handler 文件提取接口信息
2. 按规则生成目录结构和中文名称
3. 生成 `apps/server/apifox.json`
4. 调用 API Fox Import API 同步

### 3. 验证同步结果

同步完成后检查 API Fox 响应中的 `endpointFailed` 是否为 0，并在 API Fox 项目中抽查目录结构。

## 目录结构规范

API Fox 目录按**项目维度**分两级：

```
漂流瓶-H5端/
  └── 漂流瓶（H5 公开接口，路径前缀 /driftBottle，不含 /admin）

管理系统/
  ├── 漂流瓶管理（路径前缀 /driftBottle/admin）
  ├── 系统用户
  ├── 角色权限
  ├── 菜单管理
  └── ...（其余管理端子目录）
```

**规则说明**：
- 路径含 `/driftBottle/admin` → `管理系统/漂流瓶管理`
- 路径含 `/driftBottle`（不含 admin）→ `漂流瓶-H5端/漂流瓶`
- 其余 system/example 类 Tag → `管理系统/<中文名>`
- 新增业务模块的管理端接口应归入 `管理系统`，H5 接口单独一个一级目录

## Tag 英文→中文映射

脚本内置完整映射表（`_TAG_CN_MAP`），新增 Tag 时在脚本对应位置补充即可。  
同时需将新 Tag 加入 `_SYSTEM_TAGS` 集合（管理端接口）或 `_PROJECT_DIR_RULES`（路径规则接口）。

## ⚠️ 重要注意事项（踩坑记录）

### 1. 禁止开启 deleteUnmatchedResources
```python
"deleteUnmatchedResources": False  # 必须保持 False！
```
**原因**：设为 `True` 时，如果本次只同步部分模块（如只有漂流瓶），会把 API Fox 中其他所有接口全部删除。  
**历史事故**：2026-03-03 曾因此删除全部管理端接口，导致需要重新全量同步恢复。

### 2. 扫描目录必须指向 apps/server
脚本的 `root_dir` 必须是 `apps/server/`，不是 monorepo 根目录。  
当前脚本已正确配置：`repo_root / "apps" / "server"`。

### 3. 新增模块后需补充 Tag 映射
若新模块的 `@Tags` 值不在 `_TAG_CN_MAP` 中，接口会归入「其他」目录（英文名目录）。  
每次新增模块后需同步更新 `_TAG_CN_MAP` 和 `_SYSTEM_TAGS`。

### 4. API Fox 开放 API 说明
- 查询目录列表：`GET /api/v1/projects/{id}/api-folders`（需 `X-Apifox-Version` header）
- 删除目录：`DELETE /api/v1/projects/{id}/api-folders/{folderId}`
- `/v1/` 路径（非 `/api/v1/`）大部分接口返回 302 重定向，**不可用**
- 空目录无法通过 `deleteUnmatchedResources` 删除，需手动调用删除接口

## 手动同步

如需手动导入，将生成的 `apps/server/apifox.json` 通过 API Fox 的「导入 -> OpenAPI」功能导入。
