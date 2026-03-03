# API 注解格式参考

## Swagger 注解规范

Go 文件中的 API 接口需要包含以下注解：

```go
// @Tags     模块名
// @Summary  接口描述
// @Produce  application/json
// @Param    data  body      request.XXX  true  "参数说明"
// @Success  200   {object}  response.Response{data=xxx}
// @Router   /path [method]
```

## 注解说明

| 注解 | 说明 | 示例 |
|------|------|------|
| @Tags | 接口分组 | Member |
| @Summary | 接口描述 | 发送短信验证码 |
| @Produce | 响应格式 | application/json |
| @Param | 请求参数 | data body request.SendSmsCode true "手机号" |
| @Success | 成功响应 | 200 {object} response.Response |
| @Router | 路由路径 | /member/sendSmsCode [post] |

## API Fox 导入格式

生成的 `apifox.json` 遵循 OpenAPI 3.0 规范，可直接导入 API Fox。

### 导入步骤

1. 打开 API Fox 项目
2. 点击「导入」->「OpenAPI/Swagger」
3. 选择 `apifox.json` 文件
4. 选择导入策略（覆盖/合并）
5. 完成导入
