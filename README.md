# Cinderella Engine

一个基于 Go 语言实现的轻量级 BPMN 2.0 工作流引擎。

## 功能特性

### 1. 任务类型支持
- **用户任务 (UserTask)**
  - 支持普通用户任务和自动用户任务
  - 支持灵活的任务分配
    - 指定处理人 (assignee)
    - 候选用户组 (candidateUsers)
    - 候选组 (candidateGroups)
  - 内置任务监听器机制

- **服务任务 (ServiceTask)**
  - HTTP 服务任务：支持 RESTful 服务调用
  - Pipeline 服务任务：支持管道处理

### 2. 网关支持
- **排他网关 (ExclusiveGateway)**
  - 支持条件判断
  - 支持默认路径配置

- **包容网关 (InclusiveGateway)**
  - 支持多路径并行执行
  - 支持条件判断

- **并行网关 (ParallelGateway)**
  - 支持并行分支处理
  - 支持同步合并

### 3. 事件支持
- **中间捕获事件 (IntermediateCatchEvent)**
  - 支持条件事件定义
  - 支持事件触发机制

### 4. 多实例支持
- **顺序多实例 (SequentialMultiInstance)**
  - 支持循环基数设置
  - 支持完成条件配置
  - 支持集合遍历处理

- **并行多实例 (ParallelMultiInstance)**
  - 支持基础框架
  - 并行执行能力（开发中）

### 5. 流程控制
- **顺序流 (SequenceFlow)**
  - 支持条件表达式
  - 支持默认路径配置

### 6. 扩展特性

#### 变量处理
- 支持本地变量和全局变量
- 支持表达式解析和计算

#### 事件监听
- 支持任务创建事件
- 支持任务完成事件

#### 多实例特性
- 支持循环计数
- 支持完成条件判断
- 支持集合元素处理

## 快速开始

### 安装
```bash
go get github.com/go-cinderella/cinderella-engine
```

### 基本使用
1. 引入包
```go
import "github.com/go-cinderella/cinderella-engine/engine"
```

2. 配置引擎
```go
// TODO: 添加配置示例
```

3. 部署流程
```go
// TODO: 添加部署示例
```

4. 启动流程
```go
// TODO: 添加启动示例
```

## 与flowable的区别

- 状态流转：sequenceFlow的执行实例是在continueProcessOperation中删除的，而flowable的sequenceFlow的执行实例是在takeOutgoingSequenceFlowsOperation中删除的; flowNode的执行实例的删除时机跟flowable一致，都是在takeOutgoingSequenceFlowsOperation中删除的

## 贡献指南

欢迎提交 Issue 和 Pull Request。

## 许可证

[License Name] // TODO: 添加许可证信息 