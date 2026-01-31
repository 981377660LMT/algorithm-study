1. 简短结论：**没有“会话累计阈值”这一说**，只有**单次请求的上限**。  
   判断是否达到上限，只看**当前这一次请求**的 usage（以及你提交的 prompt 总量），不是会话内累加。

   ### 计算机制（单次请求）
   - **输入上限**：`prompt_tokens <= maxPromptTokens`
   - **输出上限**：`completion_tokens <= maxResponseTokens`
   - **总上限**（如果系统也有总窗限制）：`prompt_tokens + completion_tokens <= context_window`

   ### 如何判断“接近上限”

   用本次 usage 看：
   - 输入占比 = $\frac{prompt\_tokens}{maxPromptTokens}$
   - 输出占比 = $\frac{completion\_tokens}{maxResponseTokens}$

   ### 关键点
   - **usage 是单次请求的统计**，不会在会话内累加。
   - 会话内“感觉变大”的原因是：**系统把更多历史/工具/文件树打包进了当前请求**，导致这次 `prompt_tokens` 大。
   - `cached_tokens` 只是**这次输入里命中缓存的部分**，仍然计入 `prompt_tokens`，但可能不计费或更便宜。

   如果你想做“会话级监控”，只能在客户端自己**累计每次 usage**，那是统计用途，不是模型的硬限制。

2. 判断“达到窗口限制”的依据是**单次请求的输入长度**是否接近或超过模型的上下文上限，而不是会话累计。

   `agent 无状态!`

   可用的判断方式：
   1. **直接硬指标（最佳）**
   - 发送前估算：`estimated_prompt_tokens >= maxPromptTokens`（或 `>= context_window - maxResponseTokens`）
   - 发送后：若服务端返回**截断/超限错误**，说明已超。
   2. **运行时指标（次优）**
   - 看本次 `usage.prompt_tokens`，若已接近 `maxPromptTokens`（比如 >90%），下一次就应压缩。
   - 如果你还需要较大输出，则用：  
     $\text{prompt\_tokens} \ge \text{context\_window} - \text{maxResponseTokens}$ 作为阈值。
   3. **实务做法**
   - 在客户端维护“本次输入的 token 估算器”。
   - 一旦估算值 > 80–90% 上限，触发**压缩/截断**。

   结论：**用本次请求的 prompt 大小**判断是否需要压缩，而不是看历史累计。
