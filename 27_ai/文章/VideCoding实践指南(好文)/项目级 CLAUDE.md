<!-- 一个简约的模板可以如下：

- 根目录使用下面的文件
- 子目录也放置一个 CLAUDE.md 文件，可以描述类似的内容 -->

# 项目简介

- 本项目是一个 ...（描述项目目的和技术栈）
- 编程语言：python 3.12
- 核心技术栈：
  - uv/react/

# Project Meta

Codebase Repo: ...
PSM: ...
SCM: ...
Regions: China-BOE, China-North
RDS: ...
Redis: ...

# 常用命令

- `pip install -r requirements.txt`：安装依赖（请在虚拟环境下执行）。
- `python manage.py migrate`：应用数据库迁移。
- `pytest`：运行测试，需确保测试通过。

# 目录结构

- `/frontend/`: ...
- `/backend/`: ... （描述主要目录和模块）

# 开发约定（简述）

- 代码风格：xxx。
- 架构模式：使用Django框架的MTV模式；禁止直接在视图中访问数据库（应通过Model层）。
- 错误处理：尽量捕获预期异常，记录日志，避免裸 `except`。

# Workflow（按需）

- 开发新功能请新建feature分支。提交前请运行 `black` 格式化代码和 `pytest` 确认测试通过。合并需至少一人代码审查。

# 注意事项

- **重要**：切勿将密钥上传至仓库，使用环境变量注入配置
- 首次启动时如出现数据库连接错误，请检查本地Postgres服务是否启动。

# 提交与合并请求

- 参考现有提交前缀风格（`feat:`、`fix:` 等）。
- PR 需包含：变更目的与行为总结、新增配置/环境变量说明、测试命令与结果；可关联 issue/需求编号。
- 如新增接口或行为，请附关键请求/响应示例或截图（可选）。
- 引入新参数或运行模式时同步更新文档（`README.md`、`QUICKSTART.md`、`.env.example`）。

//大项目

# 更多说明

- 服务架构规约：/path/to/architecture.md
- CodeStyle：/path/to/codestyle.xml

---

# Claude Context for [Project Name]

## 1. Project Overview (The WHY)

- **Goal**: 构建一个基于 Next.js 的高性能电商仪表盘。
- **Core Principles**:
  - 优先使用服务器组件 (RSC)。
  - UI 必须支持移动端响应式。
  - **Zero-trust**: 所有用户输入必须在 Zod Schema 中验证。

## 2. Tech Stack & Map (The WHAT)

- **Stack**: Next.js 14 (App Router), TypeScript, TailwindCSS, Prisma, PostgreSQL.
- **Key Directories**:
  - `/src/app`: 页面路由 (Business Logic).
  - `/src/components/ui`: 通用组件 (Radix UI wrappers). **禁止在此添加业务逻辑**。
  - `/src/lib/actions`: Server Actions (Data Mutation).

## 3. Workflow & Commands (The HOW)

- **Start**: `pnpm dev`
- **Test**: `pnpm test` (Run this before suggesting any commit).
- **Database**: Run `pnpm prisma generate` after changing `schema.prisma`.

## 4. Documentation Index (Progressive Disclosure)

- For Database Schema rules: `read docs/database-rules.md`
- For API Response formats: `read docs/api-standards.md`

## 5. Style & Constraints

- **Naming**: Use `camelCase` for variables, `PascalCase` for components.
- **Error Handling**: Never define empty catch blocks. Log all errors to console.
- **Testing**: Use React Testing Library. Do not test implementation details, test user behavior.
