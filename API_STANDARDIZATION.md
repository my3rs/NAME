# API接口统一性改进方案

## 当前问题

### 1. 响应格式不统一
- 部分controller使用标准响应结构（model.BaseResponse等）
- 部分controller直接使用iris.Map
- JSON字段大小写不一致（Success vs success）

### 2. 需要修复的文件

#### controller/setting.go
- 问题：使用iris.Map，字段小写
- 解决：改用model.DetailResponse

#### controller/tag.go  
- 问题：使用iris.Map，字段大写
- 解决：改用model.EmptyResponse和model.ListResponse

#### controller/attachment.go
- 问题：使用iris.Map，字段大写
- 解决：创建AttachmentResponse结构或使用DetailResponse

### 3. 统一标准

#### 响应结构使用规范：
- 单个数据：model.DetailResponse
- 列表数据（分页）：model.PageResponse  
- 列表数据（不分页）：model.ListResponse
- 空响应：model.EmptyResponse
- 批量操作：model.BatchResponse

#### JSON字段规范：
- 统一使用小写：success, message, data
- 避免使用iris.Map，统一使用标准响应结构

#### HTTP状态码规范：
- 200：成功
- 400：请求参数错误
- 401：未授权
- 404：资源不存在
- 500：服务器内部错误

## 修复优先级
1. 高优先级：setting.go, tag.go, attachment.go
2. 中优先级：统一所有controller的错误处理
3. 低优先级：优化响应消息的国际化支持