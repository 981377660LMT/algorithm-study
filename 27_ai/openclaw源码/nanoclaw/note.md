# NanoClaw 源码深度解读

> **按启动顺序 + 请求链路 逐文件讲解，适合从零到一理解整个项目。**

## 一句话概括

NanoClaw 是一个**轻量级个人 AI 助手框架**，通过消息渠道（WhatsApp/Telegram/Discord 等）接收用户消息，在隔离的 Linux 容器中运行 Claude Agent，实现安全的多群组 AI 对话和定时任务。

---

## 0. 全览：启动顺序与请求链路

```
启动顺序 (main函数):
  ① ensureContainerRuntimeRunning() ─ 检查 Docker 是否可用
  ② cleanupOrphans()                ─ 清理上次残留的容器
  ③ initDatabase()                  ─ 初始化 SQLite, 跑 migration
  ④ loadState()                     ─ 从 DB 恢复游标/会话/分组状态
  ⑤ startCredentialProxy()          ─ 启动 HTTP 凭据代理 (:3001)
  ⑥ 注册 SIGTERM/SIGINT 处理       ─ 优雅关闭
  ⑦ 遍历 channels → connect()      ─ 连接所有已安装且有凭据的渠道
  ⑧ startSchedulerLoop()            ─ 启动定时任务轮询 (60s)
  ⑨ startIpcWatcher()               ─ 启动 IPC 文件监视 (1s)
  ⑩ queue.setProcessMessagesFn()    ─ 注入消息处理函数
  ⑪ recoverPendingMessages()        ─ 恢复崩溃前未处理的消息
  ⑫ startMessageLoop()              ─ 进入主消息轮询 (2s) ← 永不退出

请求链路 (一条消息从进入到回复):
  Channel.onMessage → storeMessage(SQLite)
       ↓ (2s 后被轮询到)
  startMessageLoop → getNewMessages → 检查触发词
       ↓
  queue.sendMessage(已有容器) 或 queue.enqueueMessageCheck(新建容器)
       ↓
  processGroupMessages → formatMessages(XML) → runAgent
       ↓
  runContainerAgent → docker run → stdin 写入 JSON
       ↓ (容器内)
  Agent Runner → readStdin → query(Claude Agent SDK)
       ↓
  Claude API → 工具调用 → 产生 result
       ↓
  writeOutput(stdout 标记) → 宿主机流式解析
       ↓
  onOutput 回调 → channel.sendMessage → 用户收到回复
```

---

## 1. 基础设施层

### 1.1 env.ts — 安全的 .env 读取

```typescript
// 关键设计：不把值加入 process.env，防止泄漏到子进程
export function readEnvFile(keys: string[]): Record<string, string> {
  // 只读指定的 key，忽略其他
  // 支持引号去除: "value" → value
}
```

**设计要点**：`readEnvFile` 只提取指定的 key，**不会写入 `process.env`**。这样敏感信息（API key）不会通过 `spawn` 泄漏到子进程环境变量。config.ts 只读业务配置，credential-proxy.ts 才读密钥。

### 1.2 config.ts — 配置中心

```typescript
// 业务配置（非密钥）从 .env 读取
const envConfig = readEnvFile(['ASSISTANT_NAME', 'ASSISTANT_HAS_OWN_NUMBER'])

export const ASSISTANT_NAME = process.env.ASSISTANT_NAME || envConfig.ASSISTANT_NAME || 'Andy'
export const POLL_INTERVAL = 2000 // 消息轮询间隔
export const SCHEDULER_POLL_INTERVAL = 60000 // 定时任务检查间隔
export const IPC_POLL_INTERVAL = 1000 // IPC 文件监视间隔
export const IDLE_TIMEOUT = 1800000 // 容器空闲 30min 后关闭
export const MAX_CONCURRENT_CONTAINERS = 5 // 最多同时运行 5 个容器

// 路径全部用绝对路径（容器 volume mount 要求）
export const STORE_DIR = path.resolve(PROJECT_ROOT, 'store') // SQLite + auth
export const GROUPS_DIR = path.resolve(PROJECT_ROOT, 'groups') // 分组目录
export const DATA_DIR = path.resolve(PROJECT_ROOT, 'data') // sessions/ipc/env

// 触发词正则：^@Andy\b（大小写不敏感）
export const TRIGGER_PATTERN = new RegExp(`^@${escapeRegex(ASSISTANT_NAME)}\\b`, 'i')

// 安全配置存放在项目外，容器永远无法访问
export const MOUNT_ALLOWLIST_PATH = path.join(HOME, '.config/nanoclaw/mount-allowlist.json')
export const SENDER_ALLOWLIST_PATH = path.join(HOME, '.config/nanoclaw/sender-allowlist.json')
```

**设计要点**：安全配置文件放在 `~/.config/nanoclaw/`（项目目录之外），因为项目目录会被挂载到容器中，但 home 目录下的 `.config` 不会。

### 1.3 logger.ts — 日志

```typescript
export const logger = pino({
  level: process.env.LOG_LEVEL || 'info',
  transport: { target: 'pino-pretty', options: { colorize: true } }
})
// 全局捕获未处理异常/rejection，通过 pino 输出带时间戳的日志
```

---

## 2. 数据层 — db.ts

使用 `better-sqlite3`（同步 API，在单线程 Node 中比异步方案更简单高效）。

### 2.1 Schema（7 张表）

```sql
chats             -- 所有已知聊天的元数据（JID, name, 最后活动时间, 渠道类型）
messages          -- 消息内容（PK: id+chat_jid），只存注册群组的消息
scheduled_tasks   -- 定时任务（cron/interval/once）
task_run_logs     -- 任务执行日志
router_state      -- KV 存储：last_timestamp, last_agent_timestamp
sessions          -- 每个 group 的 Claude session ID
registered_groups -- 注册群组信息（JID → folder 映射, 容器配置）
```

### 2.2 Migration 策略

不用 migration 文件，而是 **try ALTER TABLE + catch**：

```typescript
// 每次启动尝试加新列，已存在则 catch 忽略
try {
  database.exec(`ALTER TABLE messages ADD COLUMN is_bot_message INTEGER DEFAULT 0`)
} catch {
  /* column already exists */
}
// 并做数据回填
database
  .prepare(`UPDATE messages SET is_bot_message = 1 WHERE content LIKE ?`)
  .run(`${ASSISTANT_NAME}:%`)
```

**设计要点**：极简 migration，适合小型单用户应用。

### 2.3 关键查询模式

**双游标消息查询**：系统维护两个时间戳游标：

- `lastTimestamp` — "已看到"游标，Message Loop 用来发现新消息
- `lastAgentTimestamp[chatJid]` — "已处理"游标，每个群组独立，用来重组上下文

```typescript
// getNewMessages: 获取所有注册群组的新消息（用于 Message Loop 发现）
// getMessagesSince: 获取某个群组的所有待处理消息（用于构建 prompt）
// 两者都过滤 bot 消息（is_bot_message=0 AND content NOT LIKE 'Andy:%'）
// 子查询取最新 N 条再外层按时间正序 — 保证最新消息不丢，且按时间呈现
```

---

## 3. 渠道系统 — channels/

### 3.1 registry.ts — 工厂注册中心

```typescript
export type ChannelFactory = (opts: ChannelOpts) => Channel | null
const registry = new Map<string, ChannelFactory>()

export function registerChannel(name: string, factory: ChannelFactory): void {
  registry.set(name, factory)
}
```

核心不到 15 行。**Channel 接口**定义在 types.ts：

```typescript
interface Channel {
  name: string
  connect(): Promise<void> // 建立连接
  sendMessage(jid: string, text: string): Promise<void> // 发消息
  isConnected(): boolean
  ownsJid(jid: string): boolean // 判断某个 JID 是否属于此渠道
  disconnect(): Promise<void>
  setTyping?(jid: string, isTyping: boolean): Promise<void> // 可选：typing 指示
  syncGroups?(force: boolean): Promise<void> // 可选：同步群组名
}
```

**`ownsJid()` 设计精妙**：多渠道并存时，WhatsApp JID 格式 `xxx@g.us`、Telegram 用 `tg:xxx`、Discord 用 `dc:xxx`。路由时遍历所有 channel 找到 `ownsJid` 返回 true 的那个。

### 3.2 自注册流程

```
1. Skill 添加 src/channels/telegram.ts，里面在模块顶层调用 registerChannel()
2. Skill 在 src/channels/index.ts 添加 import './telegram.js'
3. index.ts 顶部 import './channels/index.js' → 触发所有 channel 注册
4. main() 中遍历 registry，对每个 factory 传入 channelOpts 调用
5. factory 检查凭据 → 有则返回 Channel 实例 → 无则返回 null（跳过）
```

### 3.3 channelOpts — 渠道回调

```typescript
const channelOpts = {
  onMessage: (chatJid, msg) => {
    // ① 拦截 /remote-control 命令
    // ② sender-allowlist drop 模式过滤
    // ③ storeMessage → SQLite
  },
  onChatMetadata: (chatJid, timestamp, name?, channel?, isGroup?) => storeChatMetadata(...),
  registeredGroups: () => registeredGroups,  // 实时引用
};
```

---

## 4. 消息循环 — index.ts 中的 startMessageLoop

这是整个系统的**心跳**，每 2 秒执行一次：

```typescript
while (true) {
  // ① 从 SQLite 获取所有注册群组的新消息
  const { messages, newTimestamp } = getNewMessages(jids, lastTimestamp, ASSISTANT_NAME);

  // ② 推进"已看到"游标
  lastTimestamp = newTimestamp;

  // ③ 按群组分组
  const messagesByGroup = new Map<string, NewMessage[]>();

  for (const [chatJid, groupMessages] of messagesByGroup) {
    // ④ 非 Main 群组检查触发词 @Andy
    if (needsTrigger) {
      const hasTrigger = groupMessages.some(m => TRIGGER_PATTERN.test(m.content.trim()));
      if (!hasTrigger) continue;  // 无触发词 → 跳过，消息累积在 DB 等下次触发
    }

    // ⑤ 拉取自上次处理以来的所有消息（包括累积的非触发消息作为上下文）
    const allPending = getMessagesSince(chatJid, lastAgentTimestamp[chatJid], ...);

    // ⑥ 格式化为 XML
    const formatted = formatMessages(allPending, TIMEZONE);

    // ⑦ 尝试发送到活跃容器（已有容器在跑）
    if (queue.sendMessage(chatJid, formatted)) {
      // 成功通过 IPC 管道追加 → 推进游标
    } else {
      // 无活跃容器 → 排队等新容器
      queue.enqueueMessageCheck(chatJid);
    }
  }
  await sleep(POLL_INTERVAL);  // 2s
}
```

**双路径设计**：

- **热路径**（`queue.sendMessage`）：群组已有活跃容器 → 通过文件 IPC 追加消息，毫秒级
- **冷路径**（`queue.enqueueMessageCheck`）：需要新建容器 → 排队 → GroupQueue 调度

### 消息格式化 — router.ts

```typescript
// 输入给 Claude 的 prompt 是 XML 格式
<context timezone="Asia/Shanghai" />
<messages>
<message sender="张三" time="3:42 PM">@Andy 今天天气怎么样？</message>
<message sender="李四" time="3:43 PM">我也想知道</message>
</messages>
```

输出过滤：Agent 可以用 `<internal>...</internal>` 包裹内部推理，`stripInternalTags` 会在发送前移除。

---

## 5. 并发控制 — group-queue.ts

`GroupQueue` 是**每组一队列 + 全局并发限制**的调度器：

```typescript
class GroupQueue {
  private groups = new Map<string, GroupState>() // 每组状态
  private activeCount = 0 // 当前活跃容器数
  private waitingGroups: string[] = [] // 等待槽位的组

  // GroupState 维护：
  // active: 是否有容器在跑
  // idleWaiting: 容器空闲等待中
  // pendingMessages: 有新消息待处理
  // pendingTasks: 排队的定时任务
  // process: 容器子进程引用
  // containerName: Docker 容器名（用于 stop）
  // retryCount: 指数退避重试计数
}
```

### 核心调度逻辑

```
enqueueMessageCheck(groupJid):
  if 已有活跃容器 → 标记 pendingMessages
  elif 超并发上限 → 加入 waitingGroups
  else → runForGroup() 立即启动

enqueueTask(groupJid, taskId, fn):
  if 已有活跃容器 → 入 pendingTasks 队列
    if 容器空闲 → closeStdin() 触发容器退出 → 排空后执行任务
  elif 超并发 → 入等待队列
  else → runTask() 立即执行

容器结束 (finally):
  activeCount--
  drainGroup():  ← 检查此组是否有 pending work
    pendingTasks.shift() → runTask()   // 任务优先
    pendingMessages → runForGroup()    // 然后消息
    否则 → drainWaiting()             // 释放槽位给等待的组
```

### 容器复用 — sendMessage

```typescript
sendMessage(groupJid: string, text: string): boolean {
  // 通过文件 IPC 向活跃容器发送后续消息
  const inputDir = path.join(DATA_DIR, 'ipc', state.groupFolder, 'input');
  // 原子写入：先写 .tmp 再 rename
  fs.writeFileSync(tempPath, JSON.stringify({ type: 'message', text }));
  fs.renameSync(tempPath, filepath);
}
```

### 重试机制

```typescript
// 指数退避：5s → 10s → 20s → 40s → 80s，最多 5 次
const delayMs = BASE_RETRY_MS * Math.pow(2, retryCount - 1)
```

---

## 6. 容器生命周期

### 6.1 container-runtime.ts — 运行时抽象

```typescript
export const CONTAINER_RUNTIME_BIN = 'docker'
export const CONTAINER_HOST_GATEWAY = 'host.docker.internal'

// Credential Proxy 绑定地址：
//   macOS/WSL → 127.0.0.1（Docker Desktop VM 代理）
//   Linux → docker0 网桥 IP（容器直连宿主机）

// 启动时检查 docker info，失败则 fatal 退出
// 清理以 nanoclaw- 开头的孤儿容器
```

### 6.2 container-runner.ts — 容器构建与启动

**Volume Mount 构建** — `buildVolumeMounts()` 是安全模型的核心：

```
Main 群组:
  /workspace/project   ← 项目根目录 (只读)
  /workspace/project/.env ← /dev/null 遮蔽 (防止读密钥)
  /workspace/group     ← groups/{folder} (读写)

Non-main 群组:
  /workspace/group     ← groups/{folder} (读写)
  /workspace/global    ← groups/global/ (只读)

通用挂载:
  /home/node/.claude   ← data/sessions/{group}/.claude/ (会话+设置)
  /workspace/ipc       ← data/ipc/{group}/ (IPC 通信)
  /app/src             ← data/sessions/{group}/agent-runner-src/ (可定制的 Agent Runner)
  /workspace/extra/*   ← 额外挂载 (经 mount-security.ts 验证)
```

**容器启动** — `buildContainerArgs()`:

```bash
docker run -i --rm --name nanoclaw-{folder}-{timestamp} \
  -e TZ=Asia/Shanghai \
  -e ANTHROPIC_BASE_URL=http://host.docker.internal:3001 \  # 走代理
  -e ANTHROPIC_API_KEY=placeholder \  # 或 CLAUDE_CODE_OAUTH_TOKEN=placeholder
  --user {uid}:{gid} \  # 以宿主机用户运行（非 root）
  -v groups/{folder}:/workspace/group \
  -v ... (其他挂载) \
  nanoclaw-agent:latest
```

**IO 协议**：

- stdin: 写入 JSON(`ContainerInput`) 后关闭
- stdout: Agent Runner 输出用标记包裹 `---NANOCLAW_OUTPUT_START---\n{json}\n---NANOCLAW_OUTPUT_END---`
- stderr: 调试日志（不触发超时重置）

**流式解析**：

```typescript
container.stdout.on('data', chunk => {
  parseBuffer += chunk
  // 在 buffer 中搜索 START/END 标记对
  while ((startIdx = parseBuffer.indexOf(OUTPUT_START_MARKER)) !== -1) {
    const endIdx = parseBuffer.indexOf(OUTPUT_END_MARKER, startIdx)
    if (endIdx === -1) break // 不完整，等更多数据
    const json = parseBuffer.slice(startIdx + MARKER_LEN, endIdx).trim()
    const parsed = JSON.parse(json)
    onOutput(parsed) // 流式回调 → 实时发消息给用户
  }
})
```

**超时策略**：

- 硬超时：`max(configTimeout, IDLE_TIMEOUT + 30s)` — 每次有输出时重置
- 软超时：IDLE_TIMEOUT (30min) → 写 `_close` → 容器自行退出
- 超时后已有输出 → 视为正常（idle cleanup），不算错误

### 6.3 Group session 设置

每个 group 首次创建时，在 `data/sessions/{group}/.claude/settings.json` 写入：

```json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1", // 启用子代理编排
    "CLAUDE_CODE_ADDITIONAL_DIRECTORIES_CLAUDE_MD": "1", // 加载额外目录的 CLAUDE.md
    "CLAUDE_CODE_DISABLE_AUTO_MEMORY": "0" // 启用自动记忆
  }
}
```

同时将 `container/skills/` 下的 Skill 同步到 group 的 `.claude/skills/`。

---

## 7. 容器内核心 — agent-runner

### 7.1 入口流程

```typescript
async function main() {
  // ① 从 stdin 读取完整的 ContainerInput JSON
  const containerInput = JSON.parse(await readStdin());

  // ② 准备环境（凭据通过代理注入，容器内是 placeholder）
  const sdkEnv = { ...process.env };

  // ③ IPC 清理 + 初始 prompt 构建
  let prompt = containerInput.prompt;
  const pending = drainIpcInput();  // 可能有累积的 IPC 消息
  if (pending.length > 0) prompt += '\n' + pending.join('\n');

  // ④ 查询循环：query → wait → query → wait
  while (true) {
    const result = await runQuery(prompt, sessionId, ...);
    if (result.closedDuringQuery) break;  // 收到关闭信号

    writeOutput({ status: 'success', result: null, newSessionId }); // session 更新标记

    const nextMessage = await waitForIpcMessage();  // 轮询 IPC input 目录
    if (nextMessage === null) break;  // _close 哨兵
    prompt = nextMessage;
  }
}
```

### 7.2 MessageStream — 异步消息流

```typescript
class MessageStream {
  private queue: SDKUserMessage[] = [];
  private waiting: (() => void) | null = null;
  private done = false;

  push(text: string): void { ... }  // 追加消息到队列
  end(): void { ... }               // 关闭流

  async *[Symbol.asyncIterator]() {
    while (true) {
      while (this.queue.length > 0) yield this.queue.shift()!;
      if (this.done) return;
      await new Promise(r => { this.waiting = r; }); // 挂起等待
    }
  }
}
```

**为什么需要 MessageStream**：SDK 的 `query()` 接受 `prompt` 参数可以是 `AsyncIterable`。当是单个字符串时，SDK 设置 `isSingleUserTurn=true`，Agent Teams 的子代理会被提前终止。用 MessageStream 保持流开放，子代理才能完整运行。

### 7.3 runQuery — 核心 Agent 调用

```typescript
for await (const message of query({
  prompt: stream, // MessageStream (AsyncIterable)
  options: {
    cwd: '/workspace/group',
    resume: sessionId, // 恢复会话
    resumeSessionAt: resumeAt, // 从特定消息恢复
    systemPrompt: globalClaudeMd // 追加全局 CLAUDE.md 到系统提示
      ? { type: 'preset', preset: 'claude_code', append: globalClaudeMd }
      : undefined,
    allowedTools: [
      'Bash', // 安全！因为在容器内
      'Read',
      'Write',
      'Edit',
      'Glob',
      'Grep', // 文件操作
      'WebSearch',
      'WebFetch', // 网络访问
      'Task',
      'TaskOutput',
      'TaskStop', // 子任务
      'TeamCreate',
      'TeamDelete',
      'SendMessage', // Agent Teams
      'TodoWrite',
      'ToolSearch',
      'Skill',
      'NotebookEdit',
      'mcp__nanoclaw__*' // IPC MCP Server 的所有工具
    ],
    permissionMode: 'bypassPermissions', // 容器已隔离，跳过权限确认
    settingSources: ['project', 'user'], // 加载 CLAUDE.md
    mcpServers: {
      nanoclaw: {
        command: 'node',
        args: [mcpServerPath], // ipc-mcp-stdio.js
        env: { NANOCLAW_CHAT_JID, NANOCLAW_GROUP_FOLDER, NANOCLAW_IS_MAIN }
      }
    },
    hooks: {
      PreCompact: [createPreCompactHook()] // compact 前归档对话
    }
  }
})) {
  if (message.type === 'system' && message.subtype === 'init') {
    newSessionId = message.session_id // 首次获取 session ID
  }
  if (message.type === 'result') {
    writeOutput({ status: 'success', result: message.result, newSessionId })
  }
}
```

**同时轮询 IPC**：`pollIpcDuringQuery()` 每 500ms 检查 IPC input 目录，发现新消息就 `stream.push()` 注入到活跃查询中。发现 `_close` 哨兵就 `stream.end()` 结束查询。

### 7.4 Pre-Compact Hook — 对话归档

Claude SDK 会在上下文过长时触发 compaction。NanoClaw 注册了 `PreCompact` hook：

```
compact 前 → 解析 transcript (JSONL) → 转成 Markdown → 存到 conversations/ 目录
```

这保证了即使 session 被压缩，完整的对话历史也有 Markdown 归档。

---

## 8. IPC MCP Server — ipc-mcp-stdio.ts

容器内的 Claude Agent 通过 MCP 协议调用宿主机功能。这是一个 **stdio 传输的 MCP Server**：

### 工具列表

| 工具             | 描述             | 权限控制                               |
| ---------------- | ---------------- | -------------------------------------- |
| `send_message`   | 中途发消息给用户 | Main 可发任意 chat, 非 Main 只能发自己 |
| `schedule_task`  | 创建定时任务     | 非 Main 只能给自己调度                 |
| `list_tasks`     | 列出定时任务     | 非 Main 只能看自己的                   |
| `pause_task`     | 暂停任务         | 非 Main 只能操作自己的                 |
| `resume_task`    | 恢复任务         | 同上                                   |
| `cancel_task`    | 取消任务         | 同上                                   |
| `update_task`    | 修改任务         | 同上                                   |
| `list_groups`    | 列出可用群组     | 仅 Main                                |
| `register_group` | 注册新群组       | 仅 Main                                |

### 通信机制

```
Agent 调 MCP 工具 → MCP Server 写 JSON 到 /workspace/ipc/{messages,tasks}/
  ↓ (文件系统，被挂载到宿主机)
IPC Watcher (宿主机) 轮询读取 → 验证权限 → 执行操作
```

**原子写入**：`writeIpcFile()` 先写 `.tmp` 后 `rename`，防止读到半写的文件。

---

## 9. IPC 通信 — ipc.ts

宿主机侧的 IPC Watcher 每 1 秒轮询 `data/ipc/` 目录：

```
data/ipc/
  {group_folder}/
    messages/   ← Agent 写的待发消息 JSON
    tasks/      ← Agent 写的任务操作 JSON
    input/      ← 宿主机写的后续消息 + _close 信号
```

### 安全授权

通过**目录名确定身份**（而非文件内容中的声明）：

```typescript
for (const sourceGroup of groupFolders) {
  const isMain = folderIsMain.get(sourceGroup) === true

  // 处理消息：验证源组是否有权发给目标
  if (isMain || (targetGroup && targetGroup.folder === sourceGroup)) {
    await deps.sendMessage(data.chatJid, data.text) // 授权通过
  } else {
    logger.warn('Unauthorized IPC message attempt blocked') // 拒绝
  }

  // 处理任务：类似的权限检查
}
```

**设计要点**：身份由 IPC 目录路径决定（`data/ipc/{group_folder}/`），而不是消息 JSON 里的 `groupFolder` 字段。因为 agent 只能写入自己被挂载的 IPC 目录，无法伪造目录名。

---

## 10. 凭据代理 — credential-proxy.ts

**核心原则**：容器永远不知道真实的 API 密钥。

```
容器内: SDK 发 HTTP 请求到 http://host.docker.internal:3001
  ↓
Credential Proxy (宿主机):
  收到请求 → 删除容器发的假凭据 → 注入真实凭据 → 转发到 api.anthropic.com
  ↓
Anthropic API: 收到带真实密钥的请求
```

两种认证模式：

```typescript
if (authMode === 'api-key') {
  // API key 模式：每个请求删除假 key，注入真 x-api-key
  delete headers['x-api-key']
  headers['x-api-key'] = secrets.ANTHROPIC_API_KEY
} else {
  // OAuth 模式：替换 Authorization Bearer token
  // 只在有 Authorization header 时替换（token 交换请求）
  // 后续请求用临时 API key，无需替换
  if (headers['authorization']) {
    headers['authorization'] = `Bearer ${oauthToken}`
  }
}
```

**安全细节**：

- 剥离 hop-by-hop headers（connection, keep-alive, transfer-encoding）
- 支持自定义上游 URL（`ANTHROPIC_BASE_URL`），可接第三方模型
- `.env` 在容器中被 `/dev/null` 遮蔽

---

## 11. 定时任务 — task-scheduler.ts

### 调度循环

每 60 秒执行一次：

```typescript
export function startSchedulerLoop(deps: SchedulerDependencies) {
  const check = async () => {
    const dueTasks = getDueTasks() // WHERE status='active' AND next_run <= now
    for (const task of dueTasks) {
      // 先更新 next_run 防止并发重复触发
      const nextRun = computeNextRun(task)
      updateTask(task.id, { next_run: nextRun })

      // 通过 GroupQueue 执行，保证并发控制
      deps.queue.enqueueTask(task.chat_jid, task.id, async () => {
        await runTask(task, deps)
      })
    }
    setTimeout(check, SCHEDULER_POLL_INTERVAL)
  }
  check()
}
```

### 防漂移算法

```typescript
export function computeNextRun(task: ScheduledTask): string | null {
  if (task.schedule_type === 'interval') {
    // 锚定到上次调度时间而非当前时间
    let next = new Date(task.next_run!).getTime() + ms
    while (next <= now) next += ms // 跳过错过的间隔
    return new Date(next).toISOString()
  }
}
```

### 任务执行

```
task → runContainerAgent({ isScheduledTask: true }) → 容器内识别到 [SCHEDULED TASK] 前缀
  → Agent 执行 → 结果通过 sendMessage 发给用户
  → 10s 后关闭容器（任务是单轮的，不需要等 30min idle）
```

---

## 12. 安全模块

### 12.1 mount-security.ts — 挂载验证

```
白名单文件: ~/.config/nanoclaw/mount-allowlist.json（不在项目中，容器无法访问）

验证流程:
  ① 必须在 allowedRoots 中声明
  ② 解析符号链接后验证（防止遍历攻击）
  ③ 不含 '..' 或绝对路径的 containerPath
  ④ 不匹配黑名单模式（.ssh/.gnupg/.aws/credentials 等 17 种）
  ⑤ nonMainReadOnly → 非 Main 组强制只读
```

### 12.2 sender-allowlist.ts — 发送者过滤

```json
// ~/.config/nanoclaw/sender-allowlist.json
{
  "default": { "allow": "*", "mode": "trigger" },
  "chats": {
    "xxx@g.us": { "allow": ["user1@s.whatsapp.net"], "mode": "drop" }
  }
}
```

两种模式：

- `trigger`: 只阻止不在白名单的人触发 @Andy，消息仍存入 DB
- `drop`: 直接丢弃不在白名单的人的消息，不存入 DB

### 12.3 group-folder.ts — 路径安全

```typescript
const GROUP_FOLDER_PATTERN = /^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$/
const RESERVED_FOLDERS = new Set(['global'])

// 验证 folder 名合法 + 解析后在 base 目录内（防止 ../../../ 遍历）
export function resolveGroupFolderPath(folder: string): string {
  assertValidGroupFolder(folder)
  const groupPath = path.resolve(GROUPS_DIR, folder)
  ensureWithinBase(GROUPS_DIR, groupPath)
  return groupPath
}
```

---

## 13. 辅助功能

### 13.1 remote-control.ts — 远程控制

Main 群组发送 `/remote-control` → 启动 `claude remote-control` 进程 → 返回 URL 给用户 → 用户在浏览器直接操作 Claude Code。

```
进程 detach (unref) + stdout/stderr 写文件 → 轮询文件等 URL → 返回给用户
重启后从 data/remote-control.json 恢复会话
```

### 13.2 容器内 Dockerfile

```dockerfile
FROM node:22-slim
# 安装 Chromium + 字体（浏览器自动化）
RUN apt-get install -y chromium fonts-noto-cjk ...
# 全局安装 agent-browser 和 claude-code
RUN npm install -g agent-browser @anthropic-ai/claude-code
# 复制 agent-runner 源码，编译 TS
COPY agent-runner/ ./
RUN npm run build
```

容器启动后，`/app/src` 被挂载为 group 专属副本，每个 group 可定制自己的 Agent Runner。

---

## 14. 状态管理与崩溃恢复

### 双游标设计

```
lastTimestamp        ─ 全局"已看到"游标，由 Message Loop 推进
lastAgentTimestamp[] ─ 每组"已处理"游标，由 processGroupMessages 推进

启动恢复: recoverPendingMessages()
  for 每个注册群组:
    if getMessagesSince(lastAgentTimestamp) 有消息:
      queue.enqueueMessageCheck()  // 重新处理
```

### 错误回滚

```typescript
// processGroupMessages 中:
const previousCursor = lastAgentTimestamp[chatJid] // 保存旧游标
lastAgentTimestamp[chatJid] = 最新消息时间 // 乐观推进

if (error && !outputSentToUser) {
  lastAgentTimestamp[chatJid] = previousCursor // 回滚！
  return false // → 触发 GroupQueue 重试
}
// 如果已经发过回复给用户，即使有错误也不回滚（防止重复发送）
```

---

## 15. 技术决策总结

| 决策                         | 理由                                                |
| ---------------------------- | --------------------------------------------------- |
| 单进程 + SQLite              | 极简运维，避免消息队列/Redis 等额外依赖             |
| 同步 SQLite (better-sqlite3) | 单线程场景比异步 API 更简单且快                     |
| 文件系统 IPC                 | 跨容器通信最简单方案，无需网络协议                  |
| 原子 rename 写入             | 防止读到半写的 JSON 文件                            |
| 触发词机制                   | 群组消息噪声大，需要 @Andy 前缀来过滤               |
| session 续接 (resume)        | 保持对话上下文，不必每次重新注入                    |
| Pre-Compact archive          | 上下文压缩前归档完整对话，避免丢失历史              |
| Credential Proxy             | 容器完全不碰真密钥，即使 Agent 被注入攻击也无法获取 |
| .env 用 /dev/null 遮蔽       | 项目根目录只读挂载时 .env 也在其中，必须单独遮蔽    |
| try-ALTER-catch migration    | 适合小应用的极简数据库迁移策略                      |
| AsyncIterable prompt         | 保持 isSingleUserTurn=false 让 Agent Teams 正常工作 |
