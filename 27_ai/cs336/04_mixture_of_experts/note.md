# Lecture 4: Mixture of Experts (MoE) | 混合专家模型

> 讲师：Percy Liang (推测)
> 主题：MoE 架构原理、路由机制、负载均衡、训练稳定性、DeepSeek V1→V3 架构演进

**核心洞见**：

1. **MoE 的本质** — 不是"不同领域的专家"，而是**稀疏激活的 FFN 副本**，同 FLOPs 换更多参数
2. **路由已收敛** — 业界几乎全部收敛到 **Token Choice Top-K** 路由，RL/最优传输/哈希等方案均被淘汰
3. **负载均衡是核心难题** — 不做均衡 → 专家坍缩（只有 1-2 个专家存活）；辅助损失 + bias 调整是标配
4. **Fine-grained Experts** — DeepSeek 提出的将专家切细的方案，是目前无争议的最佳实践
5. **架构惊人稳定** — DeepSeek V1→V3，MoE 架构图完全不变，只调路由/损失的细节

**技术解构**：

- Top-K 路由公式完整推导 + 梯度分析
- 辅助损失（F·P 内积结构）如何抑制专家坍缩
- DeepSeek V3 的 auxiliary-loss-free balancing 机制
- MLA (Multi-head Latent Attention) 压缩 KV-cache 原理
- MTP (Multi-Token Prediction) 训练策略
- Upcycling：从 Dense 模型初始化 MoE

---

## 一、为什么要 MoE？| Why MoE?

### 1.1 核心动机：同 FLOPs，更多参数

MoE 的核心 idea 极其简单：

```
Dense model:  1 个大 FFN，每次都完整计算
MoE model:    N 个小 FFN（专家），每次只激活 K 个
```

如果 K=1 且每个专家的大小 = Dense FFN 的大小，那么：

$$
\text{FLOPs}_{\text{MoE}} = \text{FLOPs}_{\text{Dense}}
$$

但参数量是 Dense 的 N 倍。**更多参数 → 更大的记忆容量（如世界知识的存储）→ 更好的性能。**

> **关键澄清**：MoE 的"experts"并**不是**领域专家（比如"代码专家""英语专家"）。这个名字极具误导性。它们只是多个被稀疏激活的 FFN 副本，没有预设的语义分工。

### 1.2 MoE 只改 FFN，不改 Attention

Transformer 的核心组件：Self-Attention + FFN。

MoE 架构与 Dense 架构的**唯一区别**：

| 组件           | Dense        | MoE                          |
| -------------- | ------------ | ---------------------------- |
| Self-Attention | 不变         | 不变                         |
| FFN            | 1 个大 FFN   | Router + N 个小 FFN（专家）  |

虽然理论上也可以对 Attention 做稀疏路由（MoE Attention），但实践中：
- 极少有人做（训练更不稳定）
- 主流 release 全部只对 FFN 做 MoE

### 1.3 实验证据：同 FLOPs 下 MoE 远胜 Dense

**证据 1：Fedus et al. 2022 (Switch Transformer)**
- 固定训练 FLOPs，增加专家数（1 → 256），训练 loss 持续下降
- 128 专家的模型比 Dense baseline 快 **7×** 达到同等 perplexity

**证据 2：OLMo (AI2, 2024/2025)**
- 精心控制的 Dense vs MoE 对比实验
- MoE（粉色线）的训练 loss 下降速度远快于 Dense（青色线）

> **结论**：在 2025 年，MoE 优于 Dense 的优势已经非常明确。几乎在所有计算规模下，如果做得好，MoE 都会给你带来超越 Dense 的收益。

### 1.4 谁在用 MoE？

| 模型          | 类型         |
| ------------- | ------------ |
| GPT-4         | 泄露为 MoE   |
| Grok          | MoE          |
| DeepSeek V2/V3| MoE          |
| LLaMA-4       | MoE          |
| Qwen 1.5 MoE | MoE          |
| Mixtral       | MoE          |
| DBRX          | MoE          |

→ 东西方都在做 MoE，已成为 2025 年构建最强模型的标准架构。

---

## 二、为什么 MoE 没有更早流行？| Why Not Earlier?

MoE 虽好，但有两个根本难题导致它长期不是默认选择：

### 2.1 基础设施复杂度

MoE 的**最大优势出现在多节点训练**。当模型大到必须分布在多个节点时，将不同专家放在不同设备上是自然的并行策略。但在单节点场景下，MoE 的优势没那么大。

### 2.2 路由决策不可微

深度学习喜欢 **可微的、平滑的** 目标函数。但路由决策是 **离散的**——必须选择并提交到某个专家。这导致：
- 优化问题很棘手
- 训练目标要么是启发式的，要么不稳定
- 需要精心的工程来让它 work

---

## 三、MoE 的三大设计问题 | Three Design Questions

设计一个 MoE，需要回答三个问题：

1. **如何路由？**（Routing function）
2. **多少个专家，每个多大？**（Expert count & size）
3. **如何训练路由器？**（Training the router）

---

## 四、路由机制详解 | Routing Mechanisms

### 4.1 三种路由范式

| 范式              | 描述                                             | 特点                     |
| ----------------- | ------------------------------------------------ | ------------------------ |
| **Token Choice**  | 每个 token 对专家打分，选 Top-K 个专家           | ✅ 几乎所有人都用这个    |
| **Expert Choice** | 每个专家对 token 打分，选 Top-K 个 token         | 天然负载均衡，但效果较差 |
| **Global Assignment** | 求解全局最优分配（线性规划/最优传输）         | 优雅但计算成本太高       |

**OLMo 的消融实验**：Token Choice 的验证损失远优于 Expert Choice。

> **结论**：业界已收敛到 **Token Choice Top-K**。

### 4.2 Token Choice Top-K 路由公式（核心）

这是 DeepSeek V1/V2 以及大多数 MoE 使用的路由公式：

**Step 1：计算亲和度分数**

$$
s_i(t) = \text{softmax}\left( \mathbf{u}^{(l)} \cdot \mathbf{e}_i \right)
$$

- $\mathbf{u}^{(l)}$：第 $l$ 层的残差流输入（token 的隐藏状态）
- $\mathbf{e}_i$：第 $i$ 个专家的**学习向量**（路由器参数，与 FFN 权重无关）
- 这就是一个**向量内积 + softmax**，类似 Attention 的打分机制

**Step 2：Top-K 选择**

$$
g_i(t) = \begin{cases} s_i(t) & \text{if } i \in \text{TopK}(s(t), K) \\ 0 & \text{otherwise} \end{cases}
$$

零掉非 Top-K 的权重。

**Step 3：加权求和 + 残差连接**

$$
\mathbf{h}^{(l)} = \sum_{i=1}^{N} g_i(t) \cdot \text{FFN}_i(\mathbf{u}^{(l)}) + \mathbf{u}^{(l)}
$$

最终输出 = 被选中专家的加权输出之和 + 原始输入（残差连接）。

### 4.3 关于路由公式的关键 Q&A

**Q：softmax 之后再 Top-K，权重还能 sum to 1 吗？**
A：不能。softmax 在 Top-K 之前做，所以选完之后不再 sum to 1。但这**不是问题**——后续有 LayerNorm，模型可以自行调整 scale。有些架构会在 Top-K 之后重新 normalize（如 DeepSeek V3），有些不会，效果差别不大。

**Q：$\mathbf{e}_i$ 向量和 FFN 权重有关系吗？**
A：没有。$\mathbf{e}_i$ 是路由器的独立参数，纯粹为了计算亲和度。

**Q：为什么路由器这么简单（就一个线性层）？**
A：
1. 更复杂的路由器（如 MLP router）消耗 FLOPs，但收益有限
2. 路由学习本身就很难——梯度信号非常间接（只能通过 Top-K 选中的专家回传）
3. 即使用更复杂的路由器，也无法保证学到最优路由

**Q：K 一般选多少？**
A：
- 最初的论文认为 K≥2，因为 K=1 没有探索（总是 exploit 最好的专家）
- K=2 最经典，让第二个专家提供一些探索信号
- K=2 意味着**翻倍 FLOPs**（激活两个专家）
- 模型报告中的 "activated parameters" 会包含这个倍数

**Q：多个专家的输出怎么合并？**
A：就是**加权求和**（weighted sum），权重就是 $g_i(t)$。

### 4.4 一个令人震惊的事实：哈希路由也有效

即使用**完全无语义信息的哈希函数**（对 token 做哈希映射到专家），MoE 仍然比 Dense 好。

**为什么？**
- 即使哈希，同类 token 总是去同一个专家 → 仍有某种"特化"
- 如果 token 频率分布是 Zipfian 的（"the"这类高频词占据某个专家），专家仍会获得某种语义特化
- 但如果是**纯随机路由**（每次随机分配，不依赖输入），效果会非常差

> **启示**：MoE 的增益很大程度来自"多参数"本身，而不仅仅是"智能路由"。

### 4.5 其他被尝试但已放弃的路由方案

| 方案                   | 思路                            | 现状                        |
| ---------------------- | ------------------------------- | --------------------------- |
| **RL 学习路由**        | 离散决策 → 用 RL 优化          | 计算成本太高，效果不比哈希好多少 |
| **线性分配/最优传输**  | 求解全局最优 token-expert 映射  | 优雅但成本远超收益          |
| **随机扰动路由**       | 给 logits 加噪声促进探索        | 被启发式损失方法替代        |

#### 随机扰动路由（Shazeer 2017）

$$
H(x)_i = (x \cdot W_g)_i + \text{StandardNormal()} \cdot \text{Softplus}((x \cdot W_{\text{noise}})_i)
$$

- 给路由 logits 加了可学习尺度的正态噪声
- 类似 $\epsilon$-greedy 探索：随机拉一些非最优的专家 bandit arm
- 缺点：测试时不加噪声 → train/test mismatch
- **已被辅助损失方法取代**

---

## 五、专家数量与大小 | Expert Count & Size

### 5.1 Fine-grained Experts（细粒度专家）—— DeepSeek 的关键创新

**核心思想**：与其有 N 个"标准大小"的专家，不如有 4N 个"1/4 大小"的专家。

```
标准 FFN: hidden_dim → 4 × hidden_dim（或 2.67× for gated）

Fine-grained: hidden_dim → 1 × hidden_dim（切成 1/4）
              但有 4× 数量的专家
```

**好处**：
- 更多专家 = 更细粒度的路由 = 更好的性能
- 如果 K 也成比例增加，FLOPs 不变（每个专家更小，激活更多个）

**DeepSeek MoE 的经典消融**：
- GShard（基础 MoE） → +共享专家 → +细粒度专家 → 逐步大幅提升

**OLMo 的验证**：
- Fine-grained experts 8 → 32 → 64，loss 和下游指标持续改善
- **无争议的最佳实践**

### 5.2 Shared Experts（共享专家）

**思路**：也许有一些处理是不管什么 token 都需要的（"共享结构"）。那就设 1-2 个专家始终激活，不参与路由。

**DeepSeek 的消融**：共享专家在某些任务上有明显提升。
**OLMo 的验证**：共享专家**几乎没有增益**。

> **结论**：Fine-grained experts 是确定性洞见；Shared experts **效果存疑**，但很多模型还是用了。

### 5.3 主流模型配置一览

| 模型             | 总专家数 | 激活数 | 共享专家 | 大致 ratio |
| ---------------- | -------- | ------ | -------- | ---------- |
| **早期 Google**  | 64-2048  | 1-2    | 无       | -          |
| Mixtral          | 8        | 2      | 无       | 1          |
| DBRX             | 16       | 4      | 无       | ~1/4       |
| Grok             | 8        | 2      | 无       | 1          |
| **DeepSeek MoE (V1)** | 64 | 6      | 2        | ~1/4       |
| DeepSeek V2      | 160      | 6      | 2        | ~1/4       |
| **DeepSeek V3**  | 256      | 8      | 1        | ~1/14      |
| Qwen 1.5 MoE    | 60       | 4      | 4        | ~1/4       |
| OLMo MoE        | 64       | 8      | 0        | ~1/8       |
| LLaMA-4          | 128      | 8      | 1        | -          |

**ratio** = 每个专家的 FFN 中间维度 / 标准 Dense FFN 中间维度。

> **模式**：中国系 MoE（DeepSeek/Qwen/MiniMax）全面采用 fine-grained + shared。最新模型趋向更多、更小的专家。

---

## 六、训练路由器：负载均衡 | Training the Router: Load Balancing

### 6.1 核心问题：专家坍缩 (Expert Collapse)

如果不加任何约束，会发生什么？

**OLMo 的实验**：不做负载均衡时：
- 2 个专家接管了 ~50% 的 token
- 其余 6 个专家几乎完全"死亡"（不接收任何 token）
- 白白浪费了大量参数 → loss 变差

**原因**：正反馈循环——某个专家碰巧表现好 → 更多 token 被路由到它 → 它得到更多训练 → 表现更好 → ... 其他专家得不到训练 → 越来越差。

### 6.2 辅助负载均衡损失（Switch Transformer, Fedus et al. 2022）

**几乎所有人都用的核心公式**：

$$
\mathcal{L}_{\text{balance}} = \alpha \cdot N \cdot \sum_{i=1}^{N} f_i \cdot P_i
$$

其中：
- $N$：专家总数
- $f_i$：**实际被路由到专家 $i$ 的 token 比例**（Top-K 之后的实际分配）
- $P_i$：**路由器给专家 $i$ 的平均概率**（softmax 之后、Top-K 之前的意图概率）
- $\alpha$：平衡系数

**为什么 $f_i \cdot P_i$ 能 work？**

对 $P_i$ 求梯度：

$$
\frac{\partial \mathcal{L}_{\text{balance}}}{\partial P_i} \propto f_i
$$

→ 梯度与 $f_i$ 成正比：获得越多 token 的专家，$P_i$ 被**压低**得越狠。这自然把 token 推向其他专家。

### 6.3 多层次的均衡需求

不仅要在专家层面均衡，还要在**设备层面**均衡（因为专家被分片到不同 GPU 上）：

| 均衡层次         | 目标                                | 公式结构     |
| ---------------- | ----------------------------------- | ------------ |
| Per-expert/batch | 每个 batch 内各专家得到均匀 token   | $\sum f_i \cdot P_i$（按专家） |
| Per-device        | 每个 GPU 得到均匀 token             | $\sum f_d \cdot P_d$（按设备） |
| Per-sequence     | 每个序列内各专家得到均匀 token      | 同结构，按序列统计 |

### 6.4 DeepSeek V3 的创新：Auxiliary-Loss-Free Balancing

DeepSeek V3 提出了一种不需要辅助损失的均衡方法：

**机制**：给每个专家加一个 bias $B_i$：

$$
\text{routing\_score}_i = s_i(t) + B_i
$$

- $B_i$ 仅用于**路由决策**，不参与门控权重计算
- 通过在线学习更新 $B_i$：

$$
B_i \leftarrow \begin{cases}
B_i + \gamma & \text{if expert } i \text{ 的 token 不够} \\
B_i - \gamma & \text{if expert } i \text{ 的 token 太多}
\end{cases}
$$

**优势**：避免辅助损失干扰主损失的梯度 → 训练更稳定。

**但实际上**：DeepSeek V3 论文继续读下去就会发现——他们**仍然加了**一个 per-sequence 辅助损失（"complementary sequence-wise auxiliary loss"），因为 $B_i$ 方法无法保证**单个序列**内的均衡。

> **所以"auxiliary-loss-free"并不完全属实**——只是去掉了 per-batch 的辅助损失，换成了 bias 方法；per-sequence 的辅助损失还在。

---

## 七、系统相关 | Systems Concerns

### 7.1 Expert Parallelism（专家并行）

MoE 天然适合一种新的并行策略——**专家并行**：

```
Step 1: Router 决定每个 token 去哪个专家
Step 2: All-to-All 通信：把 token 发送到对应专家所在的设备
Step 3: 各设备上的专家计算 FFN 输出
Step 4: All-to-All 通信：把结果返回原设备
Step 5: 加权求和 + 残差连接
```

这是你工具箱里的**又一种并行策略**，可以和 Data Parallel、Tensor Parallel、Pipeline Parallel 组合使用。

### 7.2 稀疏矩阵乘法优化

一个设备上有多个专家时，不同 token 激活不同专家 → 本质是**稀疏矩阵乘法**。

现代库（如 MegaBlocks）可以将多个小矩阵乘法融合成一个高效的稀疏矩阵运算，避免浪费 FLOPs。

### 7.3 Token Dropping（令牌丢弃）

当某个专家收到了超出其容量（设备内存限制）的 token 时：

- 超出的 token 会被**直接丢弃**（MLP 输出为 0，仅靠残差连接传递）
- 这在训练和推理时都会发生
- 推理时：由于 batch 内的其他请求会影响路由结果，**同一个 prompt 在不同 batch 中可能得到不同的输出**

> **有趣的副作用**：这解释了为什么 GPT-4 在 temperature=0 时仍然可能给出不同回复——如果它是 MoE，batch 内的其他请求会导致 token dropping 的不确定性。

### 7.4 DeepSeek V2 的 Top-M Devices 优化

Fine-grained experts 的问题：专家太多 → token 可能需要发送到太多设备 → 通信爆炸。

**解决方案**：先选 Top-M 个设备，再在这些设备内选 Top-K 个专家。

$$
\text{路由} = \text{TopK}(\text{experts on TopM devices})
$$

→ 严格控制了通信范围，使大规模训练（236B 参数）成为可能。

---

## 八、训练稳定性 | Training Stability

### 8.1 MoE 为什么不稳定？

- Softmax 是稳定性的"危险区"——路由器里就有 softmax
- 离散路由决策 → 梯度信号不连续
- 微调时尤其容易过拟合（参数量大，但大部分被冻结）

### 8.2 稳定性 tricks

**Trick 1：路由计算用 float32**

即使模型主体用 bf16/fp16，路由器的 softmax 计算**必须用 float32**——防止数值溢出导致的 loss spike。

**Trick 2：Router Z-Loss**

$$
\mathcal{L}_z = \frac{1}{B} \sum_{b=1}^{B} \left( \log \sum_{i=1}^{N} e^{x_i^{(b)}} \right)^2
$$

- 保持 softmax normalizer 接近 1
- 这个技巧**最早就是在 MoE 论文中提出的**，后来被推广到一般的 Transformer 训练
- 去掉 Z-Loss → 出现明显的 validation loss spike

### 8.3 微调/RLHF 的问题

MoE 微调容易过拟合（巨大的 train-val gap）。应对策略：

| 策略                           | 描述                                          |
| ------------------------------ | --------------------------------------------- |
| 交替 Dense/MoE 层，只微调 Dense | 隔层设 MoE，微调时冻结 MoE 层                |
| 用大量 SFT 数据                | DeepSeek 用了 1.4M 条训练样本                 |

---

## 九、Upcycling：从 Dense 初始化 MoE

### 9.1 原理

```
1. 取一个训练好的 Dense 模型
2. 复制其 FFN 为 N 份（可加微小扰动）
3. 随机初始化路由器
4. 继续训练 → 得到 MoE
```

### 9.2 效果

- **MiniCPM**：Dense → MoE upcycling，下游任务显著提升
- **Qwen**：Dense 模型 → upcycled MoE，2.7B 激活参数达到 7B Dense 的水平

> **价值**：不用从头训练 MoE（省去大量预训练成本），推理时又能享受 MoE 的稀疏激活优势。

---

## 十、DeepSeek V1 → V3 架构全景 | DeepSeek Architecture Evolution

### 10.1 DeepSeek MoE (V1)

- **规模**：16B 总参数，2.8B 激活
- **MoE 配置**：2 shared + 64 fine-grained experts，6 active routed
- **路由**：标准 Top-K routing（softmax 在 Top-K 之前）
- **均衡**：Expert-level + Device-level 辅助损失

### 10.2 DeepSeek V2

- **规模**：236B 总参数，21B 激活
- **MoE 架构**：**完全不变**（与 V1 相同的 shared + fine-grained 结构）
- **路由**：Top-K selector 公式不变
- **新增**：
  - **Top-M Devices**：先选设备再选专家，控制通信
  - **Communication Balancing Loss**：均衡输出端通信成本

### 10.3 DeepSeek V3

- **规模**：671B 总参数，37B 激活
- **MoE 架构**：**仍然不变**
- **路由变化**：
  - Sigmoid 替代 Softmax（更温和的归一化）
  - 门控权重做 normalize to 1
  - Auxiliary-loss-free balancing（$B_i$ bias 方法）
  - 但仍有 per-sequence auxiliary loss
- **去掉了**：V2 的 communication balancing loss
- **保留了**：V2 的 Top-M devices

### 10.4 DeepSeek V3 的非 MoE 创新

#### MLA (Multi-head Latent Attention)

目标：压缩 KV-cache（与 GQA 是替代方案）。

$$
\mathbf{c}_t = W^{DKV} \mathbf{h}_t \quad (\text{下投影到低维 latent vector})
$$
$$
\mathbf{K}_t = W^{UK} \mathbf{c}_t, \quad \mathbf{V}_t = W^{UV} \mathbf{c}_t \quad (\text{上投影恢复 K, V})
$$

**关键 trick**：$W^{UK}$ 可以和 Q 的投影矩阵**合并**（矩阵乘法结合律）→ 不增加额外 FLOPs。

**KV-cache 只需缓存低维的 $\mathbf{c}_t$**，而非完整的 K 和 V。

**与 RoPE 的兼容性问题**：RoPE 旋转矩阵夹在 Q 投影和 K 上投影之间，破坏了矩阵合并。**解决方案**：只对非压缩维度做 RoPE。

#### MTP (Multi-Token Prediction)

- 在标准 next-token prediction 之外，增加一个轻量的单层 Transformer
- 该模块预测**第 $t+2$ 个 token**（即往后看 2 步）
- 虽然论文画了可以预测多步的架构图，但**实际只做了 1 步额外预测**

---

## 十一、总结与要点 | Summary

```
┌────────────────────────────────────────────────────────────┐
│                    MoE 核心知识图谱                          │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  架构：只改 FFN → Router + N 个小 FFN（专家）               │
│                                                            │
│  路由：Token Choice Top-K（已成事实标准）                    │
│    - 向量内积 + softmax + top-k + 加权求和                  │
│    - K=2 最经典，K=1 也有人用                               │
│                                                            │
│  专家设计：                                                 │
│    - Fine-grained experts（切细，更多）= 无争议的最佳实践    │
│    - Shared experts = 效果存疑但很多人用                     │
│                                                            │
│  训练关键：负载均衡                                         │
│    - 不均衡 → 专家坍缩（只有 1-2 个存活）                   │
│    - F·P 辅助损失（Switch Transformer）                     │
│    - Bias 调整（DeepSeek V3）                               │
│    - Router Z-Loss（稳定性）                                │
│                                                            │
│  系统：                                                     │
│    - Expert Parallelism = 新的并行维度                      │
│    - All-to-All 通信 = 主要瓶颈                             │
│    - Top-M devices = 控制通信范围                           │
│    - Token Dropping = MoE 特有的不确定性                    │
│                                                            │
│  实践 tricks：                                              │
│    - 路由器用 float32                                       │
│    - Upcycling = Dense → MoE 低成本转换                    │
│    - 微调需要大量数据防止过拟合                              │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

**一句话总结**：MoE 是用"稀疏激活"换"更多参数"的架构，核心挑战是离散路由的训练（用启发式损失解决），2025 年已是构建最强模型的标准选择。