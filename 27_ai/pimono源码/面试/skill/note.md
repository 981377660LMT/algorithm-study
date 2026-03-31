# Pi-mono Skill 系统深入分析

## 一、整体架构概览

Skill 是 pi-mono 的**按需加载能力包**机制，实现了 [Agent Skills 标准](https://agentskills.io/specification)。核心设计理念是**渐进式披露 (Progressive Disclosure)**：

```
启动时只加载 name + description → 注入 system prompt（XML 格式）
    ↓ 任务匹配时
Agent 用 read 工具加载完整 SKILL.md → 执行技能指令
```

关键源码文件：

| 文件 | 职责 |
|------|------|
| `src/core/skills.ts` | 核心：类型定义、发现、加载、验证、prompt 格式化 |
| `src/core/resource-loader.ts` | ResourceLoader 统一管理 skill 生命周期 |
| `src/core/agent-session.ts` | `/skill:name` 命令展开、`parseSkillBlock` 解析 |
| `src/core/system-prompt.ts` | 将 skills 注入 system prompt |
| `src/modes/interactive/interactive-mode.ts` | TUI 中 skill 命令注册、渲染 |
| `src/modes/interactive/components/skill-invocation-message.ts` | Skill 调用消息折叠/展开 UI |

---

## 二、数据结构

### 2.1 Skill 接口

```typescript
interface Skill {
  name: string;                    // 小写字母+数字+连字符, 最长64字符, 必须匹配父目录名
  description: string;             // 最长1024字符, 决定模型何时调用
  filePath: string;                // SKILL.md 绝对路径
  baseDir: string;                 // 技能目录 (SKILL.md 的父目录)
  source: string;                  // 来源: "user" | "project" | "path"
  disableModelInvocation: boolean; // true = 从 system prompt 隐藏, 仅 /skill:name 可调用
}
```

### 2.2 Frontmatter 规范

```yaml
---
name: my-skill           # 必须匹配父目录名
description: >-          # 必填, 决定模型是否加载
  详细描述技能用途和触发条件
disable-model-invocation: false  # 可选, 默认 false
---
```

### 2.3 Skill 文件结构

```
my-skill/
├── SKILL.md              # 必须: frontmatter + 指令
├── scripts/              # 可选: 辅助脚本
├── references/           # 可选: 参考文档 (按需 read)
└── assets/               # 可选: 模板/资源
```

---

## 三、发现与加载机制 (skills.ts)

### 3.1 多层级发现路径 (优先级从高到低)

```
1. ~/.pi/agent/skills/       (用户全局, source="user")
2. .pi/skills/               (项目级, source="project")
3. settings.json:skills[]    (配置路径)
4. --skill <path>            (CLI, source="path", 即使 --no-skills 也生效)
```

### 3.2 发现规则

```typescript
// loadSkillsFromDirInternal() 核心逻辑：
// 根目录: 加载所有直接 .md 文件  (includeRootFiles=true)
// 子目录: 递归搜索 SKILL.md 文件  (includeRootFiles=false)
// 跳过: .开头文件、node_modules、被 .gitignore 忽略的路径
```

关键代码实现：

```typescript
function loadSkillsFromDirInternal(dir, source, includeRootFiles, ignoreMatcher?, rootDir?) {
  // 1. 构建 ignore matcher (读取 .gitignore/.ignore/.fdignore)
  addIgnoreRules(ig, dir, root);
  
  const entries = readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    // 跳过 .dotfiles 和 node_modules
    if (entry.name.startsWith(".") || entry.name === "node_modules") continue;
    
    // 符号链接解析
    if (entry.isSymbolicLink()) { /* statSync 跟踪 */ }
    
    // .gitignore 检查
    if (ig.ignores(ignorePath)) continue;
    
    // 目录 → 递归 (includeRootFiles=false)
    if (isDirectory) {
      loadSkillsFromDirInternal(fullPath, source, false, ig, root);
    }
    
    // 文件 → 根目录接受任意 .md, 子目录只接受 SKILL.md
    const isRootMd = includeRootFiles && entry.name.endsWith(".md");
    const isSkillMd = !includeRootFiles && entry.name === "SKILL.md";
  }
}
```

### 3.3 加载单个 Skill 文件

```typescript
function loadSkillFromFile(filePath, source) {
  const rawContent = readFileSync(filePath, "utf-8");
  const { frontmatter } = parseFrontmatter<SkillFrontmatter>(rawContent);
  
  // 验证 description (必填, ≤1024 字符)
  // 验证 name (匹配父目录, ≤64 字符, 仅 [a-z0-9-], 无连续--)
  // description 缺失 → 返回 null (不加载)
  // 其他验证失败 → warning diagnostics (仍然加载)
  
  return {
    skill: { name, description, filePath, baseDir, source, disableModelInvocation },
    diagnostics
  };
}
```

### 3.4 冲突处理

```typescript
// loadSkills() 中使用 Map 去重:
const skillMap = new Map<string, Skill>();   // name → Skill
const realPathSet = new Set<string>();        // 符号链接去重

// 同名 skill → 保留先发现的, 后来的产生 collision 诊断
// 同文件(符号链接) → 静默跳过
```

---

## 四、注入 System Prompt (system-prompt.ts)

### 4.1 XML 格式化

```typescript
function formatSkillsForPrompt(skills: Skill[]): string {
  // 过滤掉 disableModelInvocation=true 的
  const visibleSkills = skills.filter(s => !s.disableModelInvocation);
  
  // 生成 XML:
  // <available_skills>
  //   <skill>
  //     <name>...</name>
  //     <description>...</description>
  //     <location>...</location>          ← 文件绝对路径
  //   </skill>
  // </available_skills>
}
```

### 4.2 注入条件

```typescript
// buildSystemPrompt() 中：
// 1. 只有 read 工具可用时才注入 skills (因为模型需要 read 来加载完整内容)
// 2. skills 附加在 system prompt 末尾
if (hasRead && skills.length > 0) {
  prompt += formatSkillsForPrompt(skills);
}
```

引导语(写在 XML 之前)：
```
The following skills provide specialized instructions for specific tasks.
Use the read tool to load a skill's file when the task matches its description.
When a skill file references a relative path, resolve it against the skill directory.
```

---

## 五、Skill 命令执行流程 (agent-session.ts)

### 5.1 `/skill:name` 命令展开

```typescript
// 用户输入: /skill:brave-search query about AI
private _expandSkillCommand(text: string): string {
  if (!text.startsWith('/skill:')) return text;
  
  // 1. 解析: skillName="brave-search", args="query about AI"
  const skillName = text.slice(7, spaceIndex);
  const args = text.slice(spaceIndex + 1).trim();
  
  // 2. 查找 skill
  const skill = this.resourceLoader.getSkills().skills.find(s => s.name === skillName);
  if (!skill) return text; // 未知 skill, 原样传递
  
  // 3. 读取文件, 去掉 frontmatter
  const content = readFileSync(skill.filePath, 'utf-8');
  const body = stripFrontmatter(content).trim();
  
  // 4. 构造 skill block XML
  const skillBlock = `<skill name="${skill.name}" location="${skill.filePath}">
References are relative to ${skill.baseDir}.

${body}
</skill>`;
  
  // 5. 拼接用户参数
  return args ? `${skillBlock}\n\n${args}` : skillBlock;
}
```

### 5.2 展开的调用时机

```typescript
// agent-session.ts sendMessage() 中:
let expandedText = currentText;
if (expandPromptTemplates) {
  expandedText = this._expandSkillCommand(expandedText);    // 先展开 skill
  expandedText = expandPromptTemplate(expandedText, [...]);  // 再展开 prompt template
}

// steer() (中断当前流式输出后注入) 中同样展开:
let expandedText = this._expandSkillCommand(text);
expandedText = expandPromptTemplate(expandedText, [...]);
```

### 5.3 Skill Block 解析 (用于渲染)

```typescript
// TUI 渲染时解析 skill block:
function parseSkillBlock(text: string): ParsedSkillBlock | null {
  // 正则解析:
  // <skill name="..." location="...">content</skill>
  // 可选的后续用户消息
  const match = text.match(
    /^<skill name="([^"]+)" location="([^"]+)">\n([\s\S]*?)\n<\/skill>(?:\n\n([\s\S]+))?$/
  );
  return { name, location, content, userMessage };
}
```

---

## 六、ResourceLoader 集成 (resource-loader.ts)

### 6.1 统一资源管理

```typescript
interface ResourceLoader {
  getSkills(): { skills: Skill[]; diagnostics: ResourceDiagnostic[] };
  reload(): Promise<void>;           // 重新发现+加载所有资源
  extendResources(paths): void;      // 扩展运行时动态添加 skill 路径
}
```

### 6.2 DefaultResourceLoader 加载流程

```typescript
async reload() {
  // 1. 收集路径: settings.skills + CLI --skill + extension 提供的路径
  const skillPaths = this.noSkills
    ? this.mergePaths(cliEnabledSkills, this.additionalSkillPaths)
    : this.mergePaths([...enabledSkills, ...cliEnabledSkills], this.additionalSkillPaths);
  
  // 2. 加载
  this.updateSkillsFromPaths(skillPaths);
}

private updateSkillsFromPaths(skillPaths: string[]) {
  // loadSkills() → loadSkillsFromDir() → loadSkillFromFile()
  let skillsResult = loadSkills({ cwd, agentDir, skillPaths, includeDefaults: false });
  
  // 3. 应用 override (SDK 可自定义过滤/添加)
  const resolved = this.skillsOverride ? this.skillsOverride(skillsResult) : skillsResult;
  this.skills = resolved.skills;
}
```

### 6.3 SDK 定制点

```typescript
// 使用者可通过 skillsOverride 函数自定义:
const loader = new DefaultResourceLoader({
  skillsOverride: (current) => ({
    skills: [
      ...current.skills.filter(s => s.name.includes("browser")),  // 过滤
      customSkill,                                                   // 添加自定义
    ],
    diagnostics: current.diagnostics,
  }),
});
```

---

## 七、TUI 交互层

### 7.1 Slash 命令注册

```typescript
// interactive-mode.ts
this.skillCommands.clear();
if (this.settingsManager.getEnableSkillCommands()) {
  for (const skill of this.session.resourceLoader.getSkills().skills) {
    const commandName = `skill:${skill.name}`;
    this.skillCommands.set(commandName, skill.filePath);
    // 注册到自动补全
    skillCommandList.push({ name: commandName, description: skill.description });
  }
}
```

### 7.2 消息渲染 (折叠/展开)

```typescript
// SkillInvocationMessageComponent
class SkillInvocationMessageComponent extends Box {
  private expanded = false;
  
  // 折叠态: [skill] brave-search (ctrl+t to expand)
  // 展开态: [skill] + 完整 SKILL.md 内容 (Markdown 渲染)
}
```

### 7.3 Settings 面板

```typescript
// settings-selector.ts
// "Skill commands" 开关: 控制是否注册 /skill:name 命令
// enableSkillCommands: boolean, 默认 true
```

---

## 八、面试高频问题与关键设计点

### Q1: 为什么用渐进式披露而不是把所有 skill 内容直接塞进 system prompt？

**Token 效率**: description 通常几十字, 完整 SKILL.md 可能几千字。10 个 skill 全量注入会浪费大量 context window。渐进式只有在任务匹配时才 read 全文, 节省 token。

### Q2: Skill 和 System Prompt / Context File 的区别？

| 维度 | Skill | Context File (AGENTS.md) | System Prompt |
|------|-------|--------------------------|---------------|
| 加载 | 按需 (read) | 启动全量 | 启动全量 |
| 粒度 | description → 全文 | 全文 | 全文 |
| 来源 | 多层目录 | .pi/, ~/.pi | 配置 |
| 可扩展 | SDK override | override | override |
| 用途 | 专项能力包 | 项目规范 | 基础行为 |

### Q3: disableModelInvocation 的设计意图？

某些 skill 不适合模型自动调用(如破坏性操作), 仅允许用户**显式**通过 `/skill:name` 触发。这是安全/可控性设计。

### Q4: 冲突处理策略？

**先到先得** + 诊断警告。用 Map + realpath (解析符号链接) 双层去重。不会静默丢弃，而是产生 collision diagnostic 让用户感知。

### Q5: Skill 的安全模型？

Skills 可以指示模型执行任何操作 (包括 bash 命令)。文档明确提醒:
> "Skills can instruct the model to perform any action and may include executable code the model invokes. Review skill content before use."

安全边界靠**用户审查** + **工具层权限控制** (如 bash 确认)，而不是 skill 层本身。

### Q6: 如何和 Extension 系统协作？

Extension 可以通过 `extendResources()` 动态添加 skill 路径:
```typescript
extendResources({ skillPaths: [{ path: "path/to/SKILL.md", metadata: {...} }] })
```
这使得 Extension 能够安装新的 skill 包。

---

## 九、数据流全景图

```
                    ┌─────────────────┐
                    │   Skill Files   │
                    │  (~/.pi/skills, │
                    │  .pi/skills,    │
                    │  --skill paths) │
                    └────────┬────────┘
                             │
              loadSkillsFromDir() / loadSkills()
              ┌──────────────┴──────────────┐
              │  发现 → 解析 → 验证 → 去重   │
              │  (ignore规则, frontmatter,    │
              │   name/desc校验, collision)  │
              └──────────────┬──────────────┘
                             │
              ┌──────────────┴──────────────┐
              │     ResourceLoader          │
              │  (统一管理, override 支持)    │
              └──────────────┬──────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
  System Prompt         /skill:name           TUI 展示
  注入 XML              命令展开               命令补全
  (description)    (read → XML block)        折叠/展开
  formatSkillsFor   _expandSkillCommand()    Settings
  Prompt()                                    面板
        │                    │
        ▼                    ▼
     模型看到              用户消息中
   available_skills      <skill> block
        │                    │
        ▼                    ▼
   模型 read 加载        直接执行
   完整 SKILL.md         skill 指令
```
