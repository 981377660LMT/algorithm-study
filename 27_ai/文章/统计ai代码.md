- 目标50%代码贡献率
- 多种AI产品
- claudeCode、Cursor、OpenCode、Copilot、Codex 采集方式配置不同
- 数据hive表同步
- commit 提交时取前两周代码进行相似度匹配，三种方法算数均值>=0.5视为AI生成
  - 编辑距离、N-gram、文本token序列比较
  - 参考Git Diff Rename Threshold默认值50%的经验
- AI贡献率、提效总人天、人均提效市场、人均提效占比
- 数据洞察仪表盘

- 采集 -> 上报 -> Hive 同步 -> 数据洞察 -> 指标同步
- PostToolUse Hooks 配置 CLI

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "TEA_APP_ID=1220 TEA_CHANNEL=cn npx -y --prefix /tmp --registry xxx @dp/ab-agent-collect-event"
          }
        ]
      }
    ]
  }
}
```

```js
if (jsonData.hook_event_name === 'PostToolUse') {
    conversations.forEach((conversation) => {
        if (item.type === 'tool_use') {
            collectEvent('dev_agent_tool_call', {
                name: item.name,
                session_id: conversation.sessionId,
                uuid: item.id,
                conversation_uuid: conversation.uuid,
                input: formatToString(item.input),
                output: formatToString(item.content),
                is_error: item.isError,
                timestamp: item.timestamp,
                duration: item.duration,
                model: currentModel,
                patch: item.patch, // 上报代码，git diff 格式
                repo: repo,
                skill: conversation.skill,
            });
        }

    }
}
```

---

在统计 AI 代码贡献率的场景中，**编辑距离**、**N-gram 相似度**和**文本 Token 序列比较**是三种互补的文本相似度度量方法。通过计算这三者的算术均值（如你文档中提到的 $\ge 0.5$），可以有效地识别代码是否由 AI 生成或受到 AI 的显著影响。

以下是这三种技术的深入讲解：

### 1. 编辑距离 (Levenshtein Distance)

**核心思想**：衡量将一个字符串（AI 生成的代码）转换为另一个字符串（最终提交的代码）所需的最少单字符编辑操作次数。

- **操作类型**：
  - **插入** (Insertion)：增加一个字符。
  - **删除** (Deletion)：删除一个字符。
  - **替换** (Substitution)：将一个字符更改为另一个字符。
- **计算方式**：
  通常使用动态规划求解，时间复杂度为 $O(mn)$。
  $$Similarity = 1 - \frac{Levenshtein(s1, s2)}{\max(|s1|, |s2|)}$$
- **在代码统计中的特点**：
  - **优点**：对微小的字符变动（如改个变量名、加个分号）非常敏感。
  - **缺点**：对于代码块的搬移（Block Move）处理较差。如果 AI 生成了一段函数，用户只是将其在文件中移动了位置，编辑距离会认为变动巨大。
  - **适用场景**：检测用户是否对 AI 生成的代码进行了小修小改。

---

### 2. N-gram 相似度

**核心思想**：将文本分解为连续的 $n$ 个单元（字符或单词/Token）的滑窗集合，通过计算两个集合的重叠程度来衡量相似性。

- **定义**：
  - `n=1` (Unigram): "if", "(", "a", ">", "0"
  - `n=2` (Bigram): "if(", "(a", "a>", ">0"
- **计算方式（Jaccard 相似度）**：
  $$Similarity = \frac{|Ngram_{AI} \cap Ngram_{User}|}{|Ngram_{AI} \cup Ngram_{User}|}$$
- **在代码统计中的特点**：
  - **优点**：**位置无关性**。由于它关注的是局部片段的重复，即使代码块被重新排序或移动，N-gram 依然能捕获到大量重合的片段。
  - **缺点**：容易受到高频简单片段（如 `if (err != nil)`）的干扰，产生噪声。
  - **适用场景**：检测结构性相似，即使代码位置发生了变化。

---

### 3. 文本 Token 序列比较 (Token-based Comparison)

**核心思想**：不再从字符层面比较，而是先通过**词法分析器 (Lexer)** 将代码转换为 Token 流（如关键字、标识符、运算符、字面量），针对 Token 序列进行比较。

- **处理流程**：
  1.  **分词**：`const x = 1` $\rightarrow$ `[KEYWORD, IDENTIFIER, ASSIGN, NUMBER]`。
  2.  **归一化**（可选）：将所有变量名统一替换为 `VAR`，消除因重命名变量而导致的差异。
  3.  **序列比对**：使用 **LCS (最长公共子序列)** 算法计算两个序列的重合度。
- **计算方式**：
  $$Similarity = \frac{2 \times |LCS(Token_{AI}, Token_{User})|}{|Token_{AI}| + |Token_{User}|}$$
- **在代码统计中的特点**：
  - **优点**：**抗干扰能力强**。它能无视空格、注释、缩进的差异。如果进行了“重命名重构”，归一化后的 Token 序列依然能保持高度一致。
  - **缺点**：依赖于特定语言的 Lexer 编译器前端，实现成本比前两者高。
  - **适用场景**：识别深层的逻辑重用，防止用户通过修改格式或重命名来“洗掉”AI 贡献。

---

### 综合对比与建议

| 特性             | 编辑距离 (Levenshtein) | N-gram 相似度     | Token 序列比较 |
| :--------------- | :--------------------- | :---------------- | :------------- |
| **基础单元**     | 字符                   | 连续片段 (n-gram) | 语法 Token     |
| **对格式敏感度** | 极高（空格、注释都算） | 中等              | 低（忽略格式） |
| **对位置敏感度** | 敏感                   | **不敏感**        | 敏感 (LCS)     |
| **主要用途**     | 检测“逐字”修改         | 检测“片段”重用    | 检测“逻辑”复制 |

**为什么取三者均值？**
单一指标容易被绕过（例如：移动代码块绕过编辑距离，改变量名绕过简单的字符 N-gram）。将三者结合可以覆盖从**字符表面**、**局部片段**到**语法逻辑**的全方位追踪，从而得出更客观的 AI 贡献率数据。
