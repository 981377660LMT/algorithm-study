## Pi Coding Agent 清晰总结

**是什么**：一个极简的终端编码代理（harness），哲学是"让 Pi 适应你，而非你适应它"。

---

### 安装

```bash
npm install -g @mariozechner/pi-coding-agent
```

---

### 核心特性

| 特性           | 说明                                                                               |
| -------------- | ---------------------------------------------------------------------------------- |
| **模型支持**   | 15+ 提供商（Anthropic、OpenAI、Google、Ollama 等），会话中途可用 `/model` 切换     |
| **会话管理**   | 树形结构存储历史，`/tree` 导航分支，`/export` 导出 HTML，`/share` 上传 Gist        |
| **上下文工程** | `AGENTS.md`（项目指令）、`SYSTEM.md`（自定义系统提示）、自动压缩、按需加载 Skills  |
| **消息队列**   | `Enter` = 转向消息（打断当前工具）；`Alt+Enter` = 后续跟进（等代理完成）           |
| **扩展系统**   | TypeScript 模块，可访问工具、命令、快捷键、TUI，可实现子代理/沙箱/MCP/SSH 等       |
| **包管理**     | `pi install npm:@foo/pi-tools` 或 `pi install git:...`，用 `pi-package` 关键字分享 |

---

### 四种运行模式

1. **Interactive** — 完整 TUI 终端交互
2. **Print/JSON** — `pi -p "query"` 脚本模式 / `--mode json` 事件流
3. **RPC** — stdin/stdout JSON 协议，供非 Node 集成使用
4. **SDK** — 嵌入到自己的应用中

---

### 设计哲学（刻意不做的事）

- ❌ 无内置 MCP（让你自己用扩展加）
- ❌ 无子代理（用 tmux 或自建）
- ❌ 无权限弹窗（在容器里跑或自建）
- ❌ 无计划模式（写 `TODO.md` 或自建）
- ❌ 无后台 Bash（用 tmux，保持完全可观测）

**核心理念**：保持内核极简，所有"特性"都可以通过扩展或第三方包按需添加。

---

## tmux 快速上手

**是什么**：Terminal Multiplexer，终端复用器。一个终端窗口内运行多个"虚拟终端"，关闭终端后进程仍在后台保活。

---

### 核心概念三层结构

```
Session（会话）
  └── Window（窗口，类似标签页）
        └── Pane（面板，窗口的分割区域）
```

---

### 必备快捷键

所有快捷键需先按 **前缀键** `Ctrl+b`，然后再按目标键。

#### Session 管理

```bash
tmux                    # 新建会话
tmux new -s mywork      # 新建命名会话
tmux ls                 # 列出所有会话
tmux attach -t mywork   # 重新连接会话
tmux kill-session -t mywork
```

| 快捷键     | 作用                               |
| ---------- | ---------------------------------- |
| `Ctrl+b d` | **detach**，断开会话（进程继续跑） |
| `Ctrl+b s` | 列出并切换会话                     |
| `Ctrl+b $` | 重命名当前会话                     |

#### Window（标签页）

| 快捷键           | 作用              |
| ---------------- | ----------------- |
| `Ctrl+b c`       | 新建 window       |
| `Ctrl+b n` / `p` | 下/上一个 window  |
| `Ctrl+b 0-9`     | 跳转到指定 window |
| `Ctrl+b ,`       | 重命名 window     |
| `Ctrl+b &`       | 关闭当前 window   |

#### Pane（分屏）

| 快捷键           | 作用                          |
| ---------------- | ----------------------------- |
| `Ctrl+b %`       | 垂直分屏（左右）              |
| `Ctrl+b "`       | 水平分屏（上下）              |
| `Ctrl+b 方向键`  | 切换 pane                     |
| `Ctrl+b z`       | **zoom**，当前 pane 全屏/还原 |
| `Ctrl+b x`       | 关闭当前 pane                 |
| `Ctrl+b {` / `}` | 移动 pane 位置                |

---

### 最常用工作流

```bash
# 1. 开始工作
tmux new -s work

# 2. 分屏：左边编辑器，右边运行
Ctrl+b %

# 3. 跑个长任务，然后安全离开
Ctrl+b d      # detach，任务继续跑

# 4. 第二天回来继续
tmux attach -t work
```

---

### 与 Pi 的关系

Pi 刻意**不内置后台 Bash**，推荐用 tmux 替代：

- 每个 tmux pane = 完全可见、可直接交互的独立终端
- 可以在 pane 里直接运行 `pi`，另一个 pane 看结果
- 比隐藏的后台进程更透明、更安全

---

### 滚动查看历史输出

```
Ctrl+b [      # 进入滚动模式（copy mode）
方向键/PgUp   # 滚动
q             # 退出
```
