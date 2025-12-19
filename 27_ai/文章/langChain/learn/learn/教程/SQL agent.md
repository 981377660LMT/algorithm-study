这份文档介绍了如何使用 LangChain 构建一个 **SQL Agent**，该 Agent 能够通过自然语言回答关于 SQL 数据库的问题。

它涵盖了从数据库连接、Schema 获取、SQL 生成、安全检查到最终执行的全过程。

以下是基于文档内容的完整 TypeScript 实现指南：

### 1. 安装依赖

首先，你需要安装 LangChain 核心库、SQLite 驱动以及 TypeORM：

```bash
npm install langchain @langchain/core @langchain/openai typeorm sqlite3 zod @langchain/community
```

### 2. 完整代码实现

这段代码整合了文档中的数据库下载、连接配置、SQL 清洗、工具定义以及 Agent 创建逻辑。

```typescript
import fs from "node:fs/promises";
import path from "node:path";
import { SqlDatabase } from "langchain/sql_db";
import { DataSource } from "typeorm";
import { SystemMessage } from "@langchain/core/messages";
import { createAgent } from "langchain/agents"; // 假设使用高层级 createAgent API
import { tool } from "@langchain/core/tools";
import { ChatOpenAI } from "@langchain/openai";
import * as z from "zod";

// --- 1. 数据库设置与下载 ---
const url = "https://storage.googleapis.com/benchmarks-artifacts/chinook/Chinook.db";
const localPath = path.resolve("Chinook.db");

async function resolveDbPath() {
  try {
    await fs.access(localPath);
    return localPath;
  } catch {
    console.log("Downloading database...");
    const resp = await fetch(url);
    if (!resp.ok) throw new Error(`Failed to download DB. Status code: ${resp.status}`);
    const buf = Buffer.from(await resp.arrayBuffer());
    await fs.writeFile(localPath, buf);
    return localPath;
  }
}

let db: SqlDatabase | undefined;
async function getDb() {
  if (!db) {
    const dbPath = await resolveDbPath();
    const datasource = new DataSource({ type: "sqlite", database: dbPath });
    await datasource.initialize(); // 确保初始化
    db = await SqlDatabase.fromDataSourceParams({ appDataSource: datasource });
  }
  return db;
}

async function getSchema() {
  const db = await getDb();
  return await db.getTableInfo();
}

// --- 2. SQL 安全清洗 ---
const DENY_RE = /\b(INSERT|UPDATE|DELETE|ALTER|DROP|CREATE|REPLACE|TRUNCATE)\b/i;
const HAS_LIMIT_TAIL_RE = /\blimit\b\s+\d+(\s*,\s*\d+)?\s*;?\s*$/i;

function sanitizeSqlQuery(q: string) {
  let query = String(q ?? "").trim();
  
  // 简单的防多语句注入检查
  const semis = [...query].filter((c) => c === ";").length;
  if (semis > 1 || (query.endsWith(";") && query.slice(0, -1).includes(";"))) {
    throw new Error("multiple statements are not allowed.");
  }
  query = query.replace(/;+\s*$/g, "").trim();

  if (!query.toLowerCase().startsWith("select")) {
    throw new Error("Only SELECT statements are allowed");
  }
  if (DENY_RE.test(query)) {
    throw new Error("DML/DDL detected. Only read-only queries are permitted.");
  }
  if (!HAS_LIMIT_TAIL_RE.test(query)) {
    query += " LIMIT 5";
  }
  return query;
}

// --- 3. 定义工具 ---
const executeSql = tool(
  async ({ query }) => {
    const database = await getDb();
    const q = sanitizeSqlQuery(query);
    try {
      const result = await database.run(q);
      return typeof result === "string" ? result : JSON.stringify(result, null, 2);
    } catch (e: any) {
      return `Error: ${e?.message ?? String(e)}`;
    }
  },
  {
    name: "execute_sql",
    description: "Execute a READ-ONLY SQLite SELECT query and return results.",
    schema: z.object({
      query: z.string().describe("SQLite SELECT query to execute (read-only)."),
    }),
  }
);

// --- 4. 创建 Agent ---
async function runAgent() {
  const model = new ChatOpenAI({ model: "gpt-4o", temperature: 0 });

  const getSystemPrompt = async () => new SystemMessage(`You are a careful SQLite analyst.

Authoritative schema (do not invent columns/tables):
${await getSchema()}

Rules:
- Think step-by-step.
- When you need data, call the tool \`execute_sql\` with ONE SELECT query.
- Read-only only; no INSERT/UPDATE/DELETE/ALTER/DROP/CREATE/REPLACE/TRUNCATE.
- Limit to 5 rows unless user explicitly asks otherwise.
- If the tool returns 'Error:', revise the SQL and try again.
- Limit the number of attempts to 5.
- If you are not successful after 5 attempts, return a note to the user.
- Prefer explicit column lists; avoid SELECT *.
`);

  // 注意：createAgent 是文档中假设的高级 API，实际使用中可能需要使用 createToolCallingAgent 或类似方法
  // 这里为了匹配文档逻辑保持一致
  const agent = createAgent({
    model: model,
    tools: [executeSql],
    systemPrompt: await getSystemPrompt(), // 动态获取 Schema
  });

  // --- 5. 运行测试 ---
  const question = "Which genre, on average, has the longest tracks?";
  console.log(`User Question: ${question}`);

  const stream = await agent.stream(
    { messages: [{ role: "user", content: question }] },
    { streamMode: "values" }
  );

  for await (const step of stream) {
    const message = step.messages[step.messages.length - 1];
    console.log(`[${message.role}]: ${JSON.stringify(message.content, null, 2)}`);
  }
}

// 执行
runAgent().catch(console.error);
```

### 关键点解析

1.  **SqlDatabase Wrapper**: 使用 `langchain/sql_db` 封装了底层的数据库操作，简化了与 TypeORM 的集成。
2.  **Schema 注入**: Agent 的 System Prompt 中动态注入了 `await getSchema()`，这让 LLM 知道数据库中有哪些表和字段，从而能写出正确的 SQL。
3.  **安全机制**:
    *   **Sanitization**: `sanitizeSqlQuery` 函数强制检查 SQL 语句，禁止 `DROP`, `DELETE` 等危险操作，并强制添加 `LIMIT`。
    *   **Read-only Tool**: 工具层再次确保只执行查询操作。
4.  **自我修正**: System Prompt 指示 Agent 如果收到 "Error" 返回值，需要修正 SQL 并重试。

### 注意事项
*   文档中提到的 `HumanInTheLoopMiddleware` 部分示例代码是 Python 的，在 TypeScript/JS 环境中，通常通过 LangGraph 的 `checkpointer` 和 `interrupt` 机制来实现类似的人工审核流程。
*   请确保设置了 `OPENAI_API_KEY` 环境变量。