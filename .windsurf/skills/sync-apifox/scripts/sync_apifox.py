#!/usr/bin/env python3
"""
扫描 Go 项目 API 接口并同步到 API Fox
"""

import os
import re
import json
import requests
from pathlib import Path

def _load_skill_config() -> dict:
    """从 SKILL.md frontmatter 读取配置"""
    skill_md = Path(__file__).parent.parent / "SKILL.md"
    if not skill_md.exists():
        return {}
    content = skill_md.read_text(encoding="utf-8")
    # 解析 YAML frontmatter (--- ... ---)
    fm_match = re.search(r'^---\s*\n(.*?)\n---', content, re.DOTALL)
    if not fm_match:
        return {}
    config = {}
    for line in fm_match.group(1).splitlines():
        kv = line.split(':', 1)
        if len(kv) == 2:
            config[kv[0].strip()] = kv[1].strip()
    # 从正文提取 Token 和 项目 ID（格式: - Token: `xxx` 或 - 项目 ID: `xxx`）
    token_match = re.search(r'Token[：:][^`]*`([^`]+)`', content)
    project_match = re.search(r'项目\s*ID[：:][^`]*`([^`]+)`', content)
    if token_match:
        config['token'] = token_match.group(1)
    if project_match:
        config['project_id'] = project_match.group(1)
    return config

_SKILL_CONFIG = _load_skill_config()

# API Fox 配置（优先从 SKILL.md 读取，避免硬编码敏感信息）
APIFOX_TOKEN = _SKILL_CONFIG.get('token', '')
APIFOX_PROJECT_ID = _SKILL_CONFIG.get('project_id', '')
APIFOX_BASE_URL = "https://api.apifox.com/v1"

# 存储解析的 request 结构体
REQUEST_SCHEMAS = {}

# 存储解析的 response 结构体
RESPONSE_SCHEMAS = {}

# Tag -> [summary列表]，用于自动推断目录名
_TAG_SUMMARIES: dict = {}

# Tag 英文名 -> 中文名 映射（精确匹配优先）
_TAG_CN_MAP = {
    # 漂流瓶
    "DriftBottle":          "漂流瓶",
    # 系统用户
    "SysUser":              "系统用户",
    "Base":                 "登录认证",
    # 角色权限
    "Authority":            "角色权限",
    "Casbin":               "权限策略",
    # 菜单
    "Menu":                 "菜单管理",
    "AuthorityBtn":         "按钮权限",
    # 字典
    "SysDictionary":        "字典管理",
    "SysDictionaryDetail":  "字典详情",
    # API管理
    "SysApi":               "接口管理",
    "SysApiToken":          "接口Token",
    # 系统配置
    "System":               "系统配置",
    "SysParams":            "系统参数",
    # 操作日志
    "SysOperationRecord":   "操作日志",
    "SysLoginLog":          "登录日志",
    # JWT
    "Jwt":                  "JWT管理",
    # 自动化代码
    "AutoCode":             "自动化代码",
    "AutoCodeTemplate":     "代码模板",
    "AutoCodePackage":      "代码包管理",
    "AutoCodePlugin":       "插件管理",
    "AddFunc":              "方法扩展",
    # 版本管理
    "SysVersion":           "版本管理",
    # 导出模板
    "SysExportTemplate":    "导出模板",
    # 验证码
    "SysCaptcha":           "验证码",
    # 示例
    "ExaCustomer":          "示例-客户",
    "ExaFileUploadAndDownload": "示例-文件上传",
    "ExaAttachmentCategory":    "示例-附件分类",
    "SysSkills":            "技能管理",
    "SysError":             "错误管理",
    # 菜单权限（AuthorityMenu 归菜单管理）
    "AuthorityMenu":        "菜单管理",
    # 附件分类
    "AddCategory":          "附件分类",
    "DeleteCategory":       "附件分类",
    "GetCategoryList":      "附件分类",
    "ExaAttachmentCategory": "附件分类",
    # 数据库初始化
    "InitDB":               "数据库初始化",
    "CheckDB":              "数据库初始化",
    # 导出导入（归入导出模板）
    "ExportExcelByToken":   "导出模板",
    "ExportTemplateByToken":"导出模板",
    "SysImportTemplate":    "导出模板",
    # 接口管理相关
    "IgnoreApi":            "接口管理",
    # MCP 工具
    "mcp":                  "MCP工具",
    "AutoCodeMcp":          "MCP工具",
}

# 接口路径前缀 -> (一级目录, 二级目录)
# 规则：路径包含关键字时归入对应目录，None 表示使用 tag 中文名作为二级目录
_PROJECT_DIR_RULES = [
    ("/driftBottle/admin", ("管理系统", "漂流瓶管理")),
    ("/driftBottle",       ("漂流瓶-H5端", "漂流瓶")),
]

# 系统管理类 tag 归入「管理系统」目录
_SYSTEM_TAGS = {
    "SysUser", "Base", "Authority", "Casbin", "Menu", "AuthorityBtn",
    "SysDictionary", "SysDictionaryDetail", "SysApi", "SysApiToken",
    "System", "SysParams", "SysOperationRecord", "SysLoginLog", "Jwt",
    "AutoCode", "AutoCodeTemplate", "AutoCodePackage", "AutoCodePlugin",
    "AddFunc", "SysVersion", "SysExportTemplate", "SysCaptcha",
    "ExaCustomer", "ExaFileUploadAndDownload", "ExaAttachmentCategory",
    "SysSkills", "SysError",
    # 补充
    "AuthorityMenu",
    "AddCategory", "DeleteCategory", "GetCategoryList",
    "InitDB", "CheckDB",
    "ExportExcelByToken", "ExportTemplateByToken", "SysImportTemplate",
    "IgnoreApi",
    "mcp", "AutoCodeMcp",
}

def _get_project_folder(tag: str, path: str, tag_cn: str) -> tuple:
    """根据路径和 tag 返回 (一级目录, 二级目录) 元组"""
    for prefix, (l1, l2) in _PROJECT_DIR_RULES:
        if path.startswith(prefix):
            return (l1, l2)
    if tag in _SYSTEM_TAGS:
        return ("管理系统", tag_cn)
    return ("其他", tag_cn)

def _get_tag_cn(tag: str) -> str:
    """返回 tag 的中文名，找不到则返回原 tag"""
    return _TAG_CN_MAP.get(tag, tag)

# summary 中英文词 -> 中文替换表
_SUMMARY_REPLACE = {
    "SysDictionaryDetail": "字典详情",
    "SysDictionary":       "字典",
    "SysOperationRecord":  "操作日志",
    "SysExportTemplate":   "导出模板",
    "SysVersion":          "版本",
    "SysParams":           "系统参数",
    "SysError":            "错误日志",
    "McpTool":             "MCP工具",
    "API":                 "接口",
    "URL":                 "链接",
    "SQL":                 "SQL语句",
    "JSON":                "JSON",
    "Id":                  "ID",
}

def _fix_summary(summary: str) -> str:
    """替换 summary 中的英文技术词为中文"""
    for en, cn in _SUMMARY_REPLACE.items():
        summary = summary.replace(en, cn)
    return summary

# 动词前缀列表（推断目录名时需去掉）
_VERB_PREFIXES = [
    '获取', '查询', '搜索', '分页', '创建', '新增', '添加', '更新', '修改', '编辑',
    '删除', '设置', '重置', '上传', '下载', '导入', '导出', '同步', '发送', '校验',
    '验证', '绑定', '解绑', '激活', '禁用', '启用', '切换', '更改', '获得', '查看',
    '注册', '登录', '登出', '刷新', '初始化', '生成', '批量',
]

def _infer_tag_name(tag: str) -> str:
    """根据该 Tag 下收集的 @Summary 自动推断中文目录名"""
    summaries = _TAG_SUMMARIES.get(tag, [])
    if not summaries:
        return tag

    # 从每条 summary 中提取中文片段（去除动词前缀）
    chinese_words = []
    for s in summaries:
        # 仅保留中文字符
        zh = re.sub(r'[^\u4e00-\u9fff]', '', s)
        if not zh:
            continue
        # 去掉常见动词前缀
        for verb in _VERB_PREFIXES:
            if zh.startswith(verb):
                zh = zh[len(verb):]
                break
        if zh:
            chinese_words.append(zh)

    if not chinese_words:
        return tag

    # 统计最高频的 2~4 字前缀词
    counter: dict = {}
    for word in chinese_words:
        key = word[:4]  # 取前4个字作为候选
        for length in range(2, len(key) + 1):
            candidate = key[:length]
            counter[candidate] = counter.get(candidate, 0) + 1

    if not counter:
        return tag

    # 选出频率最高且长度最短的词
    best = max(counter.items(), key=lambda x: (x[1], -len(x[0])))
    # 必须至少出现在一半以上的 summary 中才采用
    if best[1] >= max(1, len(summaries) // 2):
        return best[0]
    return tag

def scan_request_models(root_dir: str) -> None:
    """扫描 request 模型文件，提取结构体字段"""
    model_dirs = [
        Path(root_dir) / "model" / "wansnap" / "request",
        Path(root_dir) / "model" / "secretspace" / "request",
        Path(root_dir) / "model" / "system" / "request",
        Path(root_dir) / "model" / "common" / "request",
        Path(root_dir) / "model" / "driftbottle" / "request",
    ]
    
    for model_dir in model_dirs:
        if not model_dir.exists():
            continue
        for go_file in model_dir.glob("*.go"):
            content = go_file.read_text(encoding="utf-8")
            parse_request_structs(content)

def scan_response_models(root_dir: str) -> None:
    """扫描 response 模型文件，提取结构体字段"""
    model_dirs = [
        Path(root_dir) / "model" / "wansnap" / "response",
        Path(root_dir) / "model" / "wansnap",
        Path(root_dir) / "model" / "secretspace" / "response",
        Path(root_dir) / "model" / "secretspace",
        Path(root_dir) / "model" / "system" / "response",
        Path(root_dir) / "model" / "common" / "response",
        Path(root_dir) / "service" / "wansnap",
        Path(root_dir) / "service" / "secretspace",
        Path(root_dir) / "model" / "driftbottle" / "response",
        Path(root_dir) / "model" / "driftbottle",
    ]
    
    for model_dir in model_dirs:
        if not model_dir.exists():
            continue
        for go_file in model_dir.glob("*.go"):
            content = go_file.read_text(encoding="utf-8")
            parse_response_structs(content)

def parse_response_structs(content: str) -> None:
    """解析 response 结构体"""
    # 找到所有结构体定义的开始位置
    struct_starts = list(re.finditer(r'type\s+(\w+)\s+struct\s*\{', content))
    matches = []
    
    for match in struct_starts:
        struct_name = match.group(1)
        start_pos = match.end()  # { 后面的位置
        
        # 通过计数大括号找到结构体结束位置
        brace_count = 1
        pos = start_pos
        while pos < len(content) and brace_count > 0:
            if content[pos] == '{':
                brace_count += 1
            elif content[pos] == '}':
                brace_count -= 1
            pos += 1
        
        if brace_count == 0:
            struct_body = content[start_pos:pos-1]  # 不包含最后的 }
            matches.append((struct_name, struct_body))
    
    for struct_name, struct_body in matches:
        fields = []
        # 检查是否嵌入了 GVA_MODEL
        if 'global.GVA_MODEL' in struct_body or 'GVA_MODEL' in struct_body:
            fields.extend([
                {"name": "ID", "type": "integer", "description": "主键ID", "go_type": "uint"},
                {"name": "createdAt", "type": "string", "description": "创建时间", "go_type": "time.Time"},
                {"name": "updatedAt", "type": "string", "description": "更新时间", "go_type": "time.Time"}
            ])
        
        # 匹配字段: FieldName Type `json:"name"...` // 注释
        # 支持多种格式的 tag
        lines = struct_body.split('\n')
        for line in lines:
            line = line.strip()
            if not line or line.startswith('//') or 'GVA_MODEL' in line:
                continue
            
            # 匹配字段定义，支持 interface{} 等类型
            field_match = re.match(r'(\w+)\s+([\w\[\]\*\.{}]+)\s+`([^`]+)`\s*(?://\s*(.+))?', line)
            if field_match:
                field_name, field_type, tags, comment = field_match.groups()
                
                # 从 tags 中提取 json 名称
                json_match = re.search(r'json:"([^"]+)"', tags)
                if json_match:
                    json_name = json_match.group(1).split(',')[0]  # 去掉 omitempty 等
                    if json_name == '-':
                        continue
                else:
                    json_name = field_name
                
                # 从 tags 中提取 example
                example_match = re.search(r'example:"([^"]+)"', tags)
                example = example_match.group(1) if example_match else None
                
                # 从 tags 中提取 enums
                enums_match = re.search(r'enums:"([^"]+)"', tags)
                enums = enums_match.group(1).split(',') if enums_match else None
                
                # 提取描述：优先使用行尾注释，其次使用 gorm comment
                description = comment.strip() if comment else None
                if not description:
                    comment_match = re.search(r'comment:([^;"]+)', tags)
                    if comment_match:
                        description = comment_match.group(1).strip()
                if not description:
                    description = json_name
                
                # 检查是否是嵌套结构体类型
                is_nested = False
                is_nested_array = False
                nested_struct_name = None
                
                # 去掉指针和数组符号
                clean_type = field_type.lstrip('*').lstrip('[]')
                if '.' in clean_type:
                    # 如 wansnap.Category
                    nested_struct_name = clean_type.split('.')[-1]
                    is_nested = True
                    is_nested_array = field_type.startswith('[]')
                elif clean_type[0:1].isupper() and clean_type not in ['JSON', 'UUID', 'Time']:
                    # 首字母大写且不是已知类型，可能是结构体
                    nested_struct_name = clean_type
                    is_nested = True
                    is_nested_array = field_type.startswith('[]')
                
                field_info = {
                    "name": json_name,
                    "type": go_type_to_openapi(field_type),
                    "go_type": field_type,
                    "description": description
                }
                
                if is_nested and nested_struct_name:
                    field_info["nested_struct"] = nested_struct_name
                    field_info["is_nested_array"] = is_nested_array
                    field_info["type"] = "array" if is_nested_array else "object"
                
                if example:
                    field_info["example"] = example
                if enums:
                    field_info["enum"] = enums
                
                fields.append(field_info)
        
        RESPONSE_SCHEMAS[struct_name] = fields

def parse_request_structs(content: str) -> None:
    """解析 request 结构体"""
    # 匹配结构体定义
    struct_pattern = r'type\s+(\w+)\s+struct\s*\{([^}]+)\}'
    matches = re.findall(struct_pattern, content, re.DOTALL)
    
    for struct_name, struct_body in matches:
        fields = []
        # 匹配字段: FieldName Type `json:"name" form:"formName"` // 注释
        # 支持 json 和 form 标签
        field_pattern = r'(\w+)\s+(\S+)\s+`[^`]*(?:json:"([^"]+)"|form:"([^"]+)")[^`]*`\s*(?://\s*(.+))?'
        field_matches = re.findall(field_pattern, struct_body)
        
        for field_match in field_matches:
            field_name, field_type, json_name, form_name, comment = field_match
            # 优先使用 form 标签，其次使用 json 标签
            param_name = form_name if form_name else json_name
            # 去掉 form 标签中的 [] 后缀
            param_name = param_name.rstrip('[]')
            if not param_name:
                continue
            required = 'required' in struct_body.split(field_name)[1].split('\n')[0] if field_name in struct_body else False
            fields.append({
                "name": param_name,
                "type": go_type_to_openapi(field_type),
                "description": comment.strip() if comment else param_name,
                "required": required
            })
        
        # 检查是否嵌入了 PageInfo
        if 'request.PageInfo' in struct_body or 'PageInfo' in struct_body:
            fields.extend([
                {"name": "page", "type": "integer", "description": "页码", "required": False},
                {"name": "pageSize", "type": "integer", "description": "每页数量", "required": False}
            ])
        
        REQUEST_SCHEMAS[struct_name] = fields

def go_type_to_openapi(go_type: str) -> str:
    """将 Go 类型转换为 OpenAPI 类型"""
    # 去掉指针符号
    go_type = go_type.lstrip('*')
    
    type_map = {
        "string": "string",
        "int": "integer",
        "int64": "integer",
        "int32": "integer",
        "uint": "integer",
        "uint64": "integer",
        "float64": "number",
        "float32": "number",
        "bool": "boolean",
        "[]string": "array",
        "[]int": "array",
        "[]uint": "array",
        "time.Time": "string",
        "datatypes.JSON": "object",
        "uuid.UUID": "string",
        "interface{}": "array",  # 分页结果中的 list 字段
    }
    return type_map.get(go_type, "string")

def scan_go_files(root_dir: str) -> list:
    """扫描 Go 文件提取 API 信息"""
    apis = []
    # 扫描 api/v1 下所有子目录（含 system、example、driftbottle 等）
    api_v1 = Path(root_dir) / "api" / "v1"
    api_dirs = [Path(root_dir) / "api"]
    if api_v1.exists():
        for sub in api_v1.iterdir():
            if sub.is_dir():
                api_dirs.append(sub)
    
    for api_dir in api_dirs:
        if not api_dir.exists():
            continue
        for go_file in api_dir.glob("*.go"):
            content = go_file.read_text(encoding="utf-8")
            apis.extend(parse_go_file(content))
    
    return apis

def parse_go_file(content: str) -> list:
    """解析 Go 文件中的 Swagger 注解"""
    apis = []
    
    # 匹配完整的注解块（从 @Tags 到 @Router）
    block_pattern = r'(// @Tags[\s\S]*?// @Router\s+\S+\s+\[\w+\])'
    blocks = re.findall(block_pattern, content)
    
    for block in blocks:
        api = parse_api_block(block)
        if api:
            apis.append(api)
    
    return apis

def parse_api_block(block: str) -> dict:
    """解析单个 API 注解块"""
    # 提取基本信息
    tags_match = re.search(r'@Tags\s+(\S+)', block)
    summary_match = re.search(r'@Summary\s+(.+?)$', block, re.MULTILINE)
    router_match = re.search(r'@Router\s+(\S+)\s+\[(\w+)\]', block)
    
    if not all([tags_match, summary_match, router_match]):
        return None

    tag_value = tags_match.group(1).strip()
    summary_value = summary_match.group(1).strip()
    # 收集 Tag -> Summary 映射，用于后续自动推断目录名
    _TAG_SUMMARIES.setdefault(tag_value, []).append(summary_value)

    api = {
        "tags": tag_value,
        "summary": summary_match.group(1).strip(),
        "path": router_match.group(1).strip(),
        "method": router_match.group(2).lower().strip(),
        "params": [],
        "request_body": None,
        "response_body": None
    }
    
    # 提取 @Success 注解中的响应结构体
    # 格式: @Success 200 {object} response.Response{data=xxx} 或 response.Response{data=[]xxx}
    success_match = re.search(r'@Success\s+\d+\s+\{[^}]+\}\s+(\S+)', block)
    if success_match:
        response_type = success_match.group(1).strip()
        # 提取 data= 后面的类型
        data_match = re.search(r'data=([^},]+)', response_type)
        if data_match:
            data_type = data_match.group(1).strip()
            # 检查是否是数组类型
            is_array = data_type.startswith('[]')
            if is_array:
                data_type = data_type[2:]  # 去掉 []
            # 去掉包名前缀
            struct_name = data_type.split(".")[-1] if "." in data_type else data_type
            api["response_body"] = struct_name
            api["response_is_array"] = is_array
        else:
            # 尝试匹配简单类型，如 response.Response{data=object}
            api["response_body"] = None
            api["response_is_array"] = False
    
    # 提取 @Param 注解
    # 格式: @Param name location type required "description"
    param_pattern = r'@Param\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+"([^"]+)"'
    param_matches = re.findall(param_pattern, block)
    
    for param_match in param_matches:
        name, location, param_type, required, description = param_match
        
        if location == "body":
            # body 参数，提取 request 结构体名
            struct_name = param_type.split(".")[-1] if "." in param_type else param_type
            api["request_body"] = {
                "struct_name": struct_name,
                "description": description
            }
        elif location == "formData":
            # formData 参数（文件上传）
            if "formData_params" not in api:
                api["formData_params"] = []
            api["formData_params"].append({
                "name": name,
                "type": param_type,  # file 类型
                "required": required.lower() == "true",
                "description": description
            })
        elif location == "query":
            # query 参数，检查是否是结构体类型
            struct_name = param_type.split(".")[-1] if "." in param_type else param_type
            # 如果是结构体类型（首字母大写），展开其字段
            if struct_name and struct_name[0].isupper() and struct_name in REQUEST_SCHEMAS:
                fields = REQUEST_SCHEMAS[struct_name]
                for field in fields:
                    api["params"].append({
                        "name": field["name"],
                        "in": "query",
                        "type": field["type"],
                        "required": field.get("required", False),
                        "description": field["description"]
                    })
            elif struct_name and struct_name[0].isupper():
                # 结构体类型但未解析到字段，跳过不添加 data 参数
                pass
            else:
                # 普通参数（如 uint, string 等基本类型）
                api["params"].append({
                    "name": name,
                    "in": location,
                    "type": go_type_to_openapi(param_type),
                    "required": required.lower() == "true",
                    "description": description
                })
        else:
            api["params"].append({
                "name": name,
                "in": location,  # path, header
                "type": go_type_to_openapi(param_type),
                "required": required.lower() == "true",
                "description": description
            })
    
    return api

def build_field_schema(field: dict, depth: int = 0) -> dict:
    """构建字段的 OpenAPI schema，支持嵌套结构体"""
    if depth > 3:  # 防止无限递归
        return {"type": "object", "description": field.get("description", "")}
    
    # 检查是否是嵌套结构体
    if field.get("nested_struct"):
        nested_name = field["nested_struct"]
        nested_fields = RESPONSE_SCHEMAS.get(nested_name, [])
        
        if nested_fields:
            nested_properties = {}
            for nf in nested_fields:
                nested_properties[nf["name"]] = build_field_schema(nf, depth + 1)
            
            if field.get("is_nested_array"):
                return {
                    "type": "array",
                    "description": field.get("description", f"{nested_name}列表"),
                    "items": {
                        "type": "object",
                        "properties": nested_properties
                    }
                }
            else:
                return {
                    "type": "object",
                    "description": field.get("description", nested_name),
                    "properties": nested_properties
                }
        else:
            # 未找到嵌套结构体定义
            if field.get("is_nested_array"):
                return {"type": "array", "description": field.get("description", ""), "items": {"type": "object"}}
            else:
                return {"type": "object", "description": field.get("description", "")}
    
    # 普通字段
    prop = {
        "type": field["type"],
        "description": field["description"]
    }
    
    # 添加 example
    if field.get("example"):
        if field["type"] == "integer":
            try:
                prop["example"] = int(field["example"])
            except:
                prop["example"] = field["example"]
        elif field["type"] == "number":
            try:
                prop["example"] = float(field["example"])
            except:
                prop["example"] = field["example"]
        elif field["type"] == "boolean":
            prop["example"] = field["example"].lower() == "true"
        else:
            prop["example"] = field["example"]
    
    # 添加 enum
    if field.get("enum"):
        if field["type"] == "integer":
            prop["enum"] = [int(e) for e in field["enum"] if e.strip()]
        else:
            prop["enum"] = [e.strip() for e in field["enum"] if e.strip()]
    
    return prop

def generate_openapi(apis: list) -> dict:
    """生成 OpenAPI 3.0 文档"""
    doc = {
        "openapi": "3.0.0",
        "info": {
            "title": "API",
            "description": "APP后端API接口",
            "version": "1.0.0"
        },
        "servers": [
            {"url": "http://localhost:10020", "description": "开发环境"}
        ],
        "tags": [],  # 动态生成，在所有 API 解析后填充
        "paths": {},
        "components": {
            "schemas": {
                "Response": {
                    "type": "object",
                    "properties": {
                        "code": {"type": "integer", "description": "状态码"},
                        "msg": {"type": "string", "description": "消息"},
                        "data": {"type": "object", "description": "数据"}
                    }
                }
            },
            "securitySchemes": {
                "MemberAuth": {
                    "type": "apiKey",
                    "in": "header",
                    "name": "x-token",
                    "description": "会员JWT Token，通过登录接口获取，放在请求头 x-token 中"
                }
            }
        }
    }
    
    # 公开接口路径关键词
    public_keywords = ["sendSmsCode", "phoneLogin", "wechatLogin", "appleLogin", "login", "captcha", "Public", "webhook"]
    
    for api in apis:
        path = api["path"]
        method = api["method"]
        
        if path not in doc["paths"]:
            doc["paths"][path] = {}
        
        # 中文 tag 名（二级目录）
        raw_tag = api["tags"] if api["tags"] else ""
        tag_cn = _get_tag_cn(raw_tag)
        # (一级目录, 二级目录)
        l1, l2 = _get_project_folder(raw_tag, path, tag_cn)
        # API Fox 用 tag 表达层级：「一级目录/二级目录」
        tag_name = f"{l1}/{l2}" if l1 else l2
        
        # 构建响应 schema
        response_schema = {"$ref": "#/components/schemas/Response"}
        
        # 如果有具体的响应结构体，使用它
        if api.get("response_body"):
            struct_name = api["response_body"]
            is_array = api.get("response_is_array", False)
            fields = RESPONSE_SCHEMAS.get(struct_name, [])
            
            if fields:
                data_properties = {}
                for field in fields:
                    data_properties[field["name"]] = build_field_schema(field)
                
                # 构建 data 结构
                if is_array:
                    data_schema = {
                        "type": "array",
                        "description": f"{struct_name}列表",
                        "items": {
                            "type": "object",
                            "properties": data_properties
                        }
                    }
                else:
                    data_schema = {
                        "type": "object",
                        "description": struct_name,
                        "properties": data_properties
                    }
                
                response_schema = {
                    "type": "object",
                    "properties": {
                        "code": {"type": "integer", "description": "状态码", "example": 0},
                        "msg": {"type": "string", "description": "消息", "example": "成功"},
                        "data": data_schema
                    }
                }
        
        operation = {
            "tags": [tag_name] if tag_name else [],
            "x-apifox-folder": tag_name,
            "summary": _fix_summary(api["summary"]),
            "operationId": f"{path.replace('/', '_')}_{method}",
            "responses": {
                "200": {
                    "description": "成功",
                    "content": {
                        "application/json": {
                            "schema": response_schema
                        }
                    }
                }
            }
        }
        
        # 添加 query/path/header 参数
        if api.get("params"):
            operation["parameters"] = []
            for param in api["params"]:
                operation["parameters"].append({
                    "name": param["name"],
                    "in": param["in"],
                    "description": param["description"],
                    "required": param["required"],
                    "schema": {"type": param["type"]}
                })
        
        # 添加 formData 参数（文件上传）
        if api.get("formData_params"):
            properties = {}
            required_fields = []
            
            for param in api["formData_params"]:
                if param["type"] == "file":
                    properties[param["name"]] = {
                        "type": "string",
                        "format": "binary",
                        "description": param["description"]
                    }
                else:
                    properties[param["name"]] = {
                        "type": param["type"],
                        "description": param["description"]
                    }
                if param["required"]:
                    required_fields.append(param["name"])
            
            request_schema = {
                "type": "object",
                "properties": properties
            }
            if required_fields:
                request_schema["required"] = required_fields
            
            operation["requestBody"] = {
                "description": "文件上传",
                "required": True,
                "content": {
                    "multipart/form-data": {
                        "schema": request_schema
                    }
                }
            }
        # 添加 request body
        elif api.get("request_body"):
            struct_name = api["request_body"]["struct_name"]
            description = api["request_body"]["description"]
            
            # 从 REQUEST_SCHEMAS 获取字段信息
            fields = REQUEST_SCHEMAS.get(struct_name, [])
            
            properties = {}
            required_fields = []
            
            for field in fields:
                properties[field["name"]] = {
                    "type": field["type"],
                    "description": field["description"]
                }
                if field["required"]:
                    required_fields.append(field["name"])
            
            # 如果没有解析到字段，使用默认结构
            if not properties:
                properties = {"data": {"type": "object", "description": description}}
            
            request_schema = {
                "type": "object",
                "properties": properties
            }
            if required_fields:
                request_schema["required"] = required_fields
            
            operation["requestBody"] = {
                "description": description,
                "required": True,
                "content": {
                    "application/json": {
                        "schema": request_schema
                    }
                }
            }
        
        # 为所有接口添加 x-appid 请求头参数
        if "parameters" not in operation:
            operation["parameters"] = []
        operation["parameters"].append({
            "name": "x-appid",
            "in": "header",
            "description": "应用ID",
            "required": True,
            "schema": {"type": "string"},
            "example": "wp0006490b11f4fcc0"
        })
        
        # 判断是否需要认证
        need_auth = not any(kw in path for kw in public_keywords)
        if need_auth:
            operation["security"] = [{"MemberAuth": []}]
            # 添加 x-token 请求头参数
            operation["parameters"].append({
                "name": "x-token",
                "in": "header",
                "description": "会员JWT Token，通过登录接口获取",
                "required": True,
                "schema": {"type": "string"},
                "example": "{{token}}"
            })
        
        doc["paths"][path][method] = operation

    # 动态填充 tags 列表（去重，保持顺序）
    seen_tags = set()
    for api in apis:
        raw_tag = api["tags"] if api["tags"] else ""
        tag_cn = _get_tag_cn(raw_tag)
        l1, l2 = _get_project_folder(raw_tag, api["path"], tag_cn)
        tag_name = f"{l1}/{l2}" if l1 else l2
        if tag_name and tag_name not in seen_tags:
            seen_tags.add(tag_name)
            doc["tags"].append({"name": tag_name, "description": f"{tag_cn}相关接口"})

    return doc

def sync_to_apifox(doc: dict) -> None:
    """同步到 API Fox"""
    url = f"{APIFOX_BASE_URL}/projects/{APIFOX_PROJECT_ID}/import-openapi"
    
    headers = {
        "Authorization": f"Bearer {APIFOX_TOKEN}",
        "Content-Type": "application/json",
        "X-Apifox-Api-Version": "2024-03-28"
    }
    
    # API Fox 需要的请求格式
    # input 直接是 OpenAPI JSON 字符串
    # deleteUnmatchedResources: 删除项目中已存在但在导入数据中不存在的资源
    payload = {
        "input": json.dumps(doc, ensure_ascii=False),
        "options": {
            "endpointOverwriteBehavior": "OVERWRITE_EXISTING",
            "schemaOverwriteBehavior": "OVERWRITE_EXISTING",
            "updateFolderOfChangedEndpoint": True,
            "prependBasePath": False,
            "deleteUnmatchedResources": True
        }
    }
    
    print(f"请求 URL: {url}")
    print(f"Payload input 长度: {len(payload['input'])} 字符")
    
    try:
        resp = requests.post(url, json=payload, headers=headers, timeout=30)
        print(f"API Fox 响应 ({resp.status_code}): {resp.text}")
    except Exception as e:
        print(f"同步失败: {e}")

def main():
    # 获取项目根目录（脚本在 .windsurf/skills/sync-apifox/scripts/ 下，向上4层是 monorepo 根）
    script_dir = Path(__file__).parent
    repo_root = script_dir.parent.parent.parent.parent
    # Go server 代码在 apps/server/ 下
    root_dir = repo_root / "apps" / "server"
    
    print(f"扫描目录: {root_dir}")
    
    # 先扫描 request 模型
    scan_request_models(str(root_dir))
    print(f"发现 {len(REQUEST_SCHEMAS)} 个 request 结构体")
    
    # 扫描 response 模型
    scan_response_models(str(root_dir))
    print(f"发现 {len(RESPONSE_SCHEMAS)} 个 response 结构体")
    
    # 扫描 API
    apis = scan_go_files(str(root_dir))
    print(f"发现 {len(apis)} 个 API 接口")
    
    # 生成 OpenAPI 文档
    doc = generate_openapi(apis)
    
    # 保存到文件
    output_file = root_dir / "apifox.json"
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(doc, f, ensure_ascii=False, indent=2)
    print(f"已生成: {output_file}")
    
    # 同步到 API Fox
    sync_to_apifox(doc)

if __name__ == "__main__":
    main()
