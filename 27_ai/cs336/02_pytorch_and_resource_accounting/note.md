# Lecture 2: PyTorch Primitives & Resource Accounting | PyTorch 原语与资源审计

> 讲师：Percy Liang
> 主题：Tensor 原语、浮点表示、内存审计、FLOPs 审计、MFU、模型构建与训练循环

**核心洞见**：

1. **资源审计是心态，不是技术** — 深度学习的效率来自"知道每一个 byte 和 flop 去了哪里"
2. **6× 法则** — 训练总 FLOPs ≈ 6 × 参数量 × Token 数（前向 2×，反向 4×）
3. **MFU 是效率的核尺** — 实际 flops/s ÷ 硬件承诺 flops/s，>0.5 算合格
4. **数据类型决定一切** — BF16 是当前甜蜜点：float32 的动态范围 + float16 的内存

**技术解构**：

- Tensor 存储模型：storage + stride 机制，view vs copy 的内存代价
- 浮点表示谱系：FP32 → FP16 → BF16 → FP8 的权衡图
- 矩阵乘法 FLOPs = 2 × 三维度之积，支配全部计算
- 内存四大户：Parameters + Activations + Gradients + Optimizer States
- einops / einsum：命名维度替代索引地狱

**关键反思**：

- 效率不是"跑完再优化"，而是"训练前先算清账"
- 所有 napkin math 的基础都是矩阵乘法
- 硬件 spec sheet 上的数字只是上界，benchmark 才是真相

---

## 一、开场：Napkin Math 的力量

Percy 一开场就抛出两个"餐巾纸计算"问题，目的是：**在写第一行代码之前，你就应该知道训练要花多久、需要多少内存。**

### 问题 1：70B 模型在 1024 张 H100 上训 15T tokens 要多久？

$$
\text{总 FLOPs} = 6 \times N_{\text{params}} \times N_{\text{tokens}} = 6 \times 70 \times 10^9 \times 15 \times 10^{12} = 6.3 \times 10^{24}
$$

$$
\text{每天可用 FLOPs} = 1024 \times \text{H100 flops/s} \times 86400 \times \text{MFU}
$$

取 MFU = 0.5，H100 BF16 无稀疏峰值约 990 TFLOPS：

$$
\text{每天 FLOPs} = 1024 \times 990 \times 10^{12} \times 86400 \times 0.5 \approx 4.38 \times 10^{22}
$$

$$
\text{天数} = \frac{6.3 \times 10^{24}}{4.38 \times 10^{22}} \approx 144 \text{ 天}
$$

> **洞见**：6× 中的 "6" 来自前向 2× + 反向 4×，本讲会推导这个数字。

### 问题 2：8 张 H100 最大能训多大的模型？（AdamW，不做花活）

$$
\text{总内存} = 8 \times 80\text{GB} = 640\text{GB}
$$

AdamW 每个参数需要 **16 bytes**（后面推导），不算 activations：

$$
N_{\text{params}} = \frac{640 \times 10^9}{16} = 40 \times 10^9 \approx 40\text{B}
$$

> **注意**：这里忽略了 activations（取决于 batch size 和 sequence length），实际可训练参数量更少。

---

## 二、Tensor：深度学习的原子

### 2.1 创建与基础

Tensor 是存储一切的容器：参数、梯度、优化器状态、数据、激活值。

```python
import torch

# 默认 float32，在 CPU 上
x = torch.randn(4, 8)
x.dtype   # torch.float32
x.shape   # torch.Size([4, 8])
x.numel() # 32
```

### 2.2 内存计算公式

$$
\text{内存 (bytes)} = \text{元素数量} \times \text{每个元素大小 (bytes)}
$$

```python
x = torch.randn(4, 8)  # float32
memory = x.numel() * x.element_size()  # 32 × 4 = 128 bytes
```

**真实规模感受**：GPT-3 FFN 层中的一个权重矩阵 $(12288 \times 49152)$ → 约 **2.3 GB**（float32）。

---

## 三、浮点表示：从 FP32 到 FP8

### 3.1 浮点格式全景

| 格式 | 总位数 | Sign | Exponent | Fraction | 每元素 Bytes | 动态范围 | 精度 |
|------|--------|------|----------|----------|-------------|---------|------|
| **FP32** | 32 | 1 | 8 | 23 | 4 | ★★★★★ | ★★★★★ |
| **FP16** | 16 | 1 | 5 | 10 | 2 | ★★☆☆☆ | ★★★☆☆ |
| **BF16** | 16 | 1 | 8 | 7 | 2 | ★★★★★ | ★★☆☆☆ |
| **FP8 (E5M2)** | 8 | 1 | 5 | 2 | 1 | ★★★☆☆ | ★☆☆☆☆ |
| **FP8 (E4M3)** | 8 | 1 | 4 | 3 | 1 | ★★☆☆☆ | ★★☆☆☆ |

- **Exponent** → 动态范围（能表示多大/多小的数）
- **Fraction** → 精度/分辨率（相邻可表示数字之间的间距）

### 3.2 为什么 BF16 是当前最佳选择？

**FP16 的致命问题**：动态范围太小。

```python
torch.tensor(1e-8, dtype=torch.float16)
# → tensor(0., dtype=torch.float16)  ← 下溢为 0！
```

**BF16 的设计哲学**：深度学习更需要动态范围，而非极高精度。

```python
torch.tensor(1e-8, dtype=torch.bfloat16)
# → tensor(1.0134e-08, dtype=torch.bfloat16)  ← 不为零，够用
```

BF16 = FP32 的动态范围 + FP16 的内存开销。

> **产业共识**：BF16 由 Google Brain 在 2018 年设计（"Brain Float"），现已成为 LLM 训练的事实标准。

### 3.3 FP8：激进的前沿

- 2022 年由 NVIDIA 推出，H100 原生支持
- 前代 GPU（A100）**不支持** FP8
- 两种变体：E5M2（更大动态范围）和 E4M3（更高精度）
- 目前主要用于推理量化，训练中还需要大量 trick

### 3.4 混合精度训练 (Mixed Precision Training)

**核心原则**：

| 阶段 | 推荐精度 | 原因 |
|------|---------|------|
| 参数主拷贝 (master weights) | FP32 | 长期累积更新，需要高精度 |
| 优化器状态 (m, v) | FP32 | 同上 |
| 前向传播 | BF16 | 计算密集，低精度 = 快 |
| 反向传播 | BF16 | 同上 |
| 梯度累积 | FP32 | 避免小梯度下溢 |
| Attention 计算 | 可能需要 FP32 | 数值敏感 |

> **Percy 的总结**：BF16 可以看作"临时的"——你把 FP32 参数 cast 到 BF16 去算前向/反向，但长期积累的东西（参数、优化器状态）必须保持 FP32。

**历史追溯**：混合精度的思想可追溯到 2017 年的论文 *"Mixed Precision Training"* (Micikevicius et al.)。

---

## 四、GPU 与数据传输

### 4.1 CPU vs GPU：Tensor 在哪？

```python
# 默认在 CPU
x = torch.zeros(32, 32)   # x.device = cpu

# 显式移到 GPU
x_gpu = x.to('cuda')      # 数据从 CPU RAM → GPU HBM

# 直接在 GPU 创建（避免传输开销）
y_gpu = torch.zeros(32, 32, device='cuda')
```

**心智模型**：

```
┌────────────┐     PCIe / NVLink     ┌────────────┐
│    CPU     │  ←─────────────────→  │    GPU     │
│   System   │    数据传输有代价！     │   H100     │
│    RAM     │                       │  80GB HBM  │
└────────────┘                       └────────────┘
```

> **原则**：永远清楚每个 tensor 在哪。如有必要，加 `assert x.device.type == 'cuda'` 来保证。

### 4.2 验证内存分配

```python
before = torch.cuda.memory_allocated()
x = torch.zeros(32, 32, device='cuda')
y = torch.zeros(32, 32, device='cuda')
after = torch.cuda.memory_allocated()

assert after - before == 2 * 32 * 32 * 4  # 2 个 FP32 矩阵 = 8192 bytes
```

---

## 五、Tensor 内部机制：Storage + Stride

### 5.1 Tensor 不是矩阵，是 metadata + 指针

PyTorch 的 Tensor 由两部分组成：

1. **Storage**：一段连续的一维内存
2. **Metadata**：shape, stride, offset — 告诉你如何索引 storage

```
逻辑视图 (4×4 矩阵):          物理存储 (一维数组):
┌───┬───┬───┬───┐
│ 0 │ 1 │ 2 │ 3 │             [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
├───┼───┼───┼───┤              ↑
│ 4 │ 5 │ 6 │ 7 │           stride = (4, 1)
├───┼───┼───┼───┤
│ 8 │ 9 │10 │11 │           element[i,j] = storage[i * stride[0] + j * stride[1]]
├───┼───┼───┼───┤           element[1,2] = storage[1*4 + 2*1] = storage[6] = 6 ✓
│12 │13 │14 │15 │
└───┴───┴───┴───┘
```

### 5.2 View：零拷贝操作

多个 Tensor 可以共享同一个 Storage，只是 metadata 不同：

```python
x = torch.tensor([[1, 2, 3], [4, 5, 6]])  # shape (2, 3)

y = x[0]             # 取第 0 行 → view，不拷贝
z = x[:, 1]          # 取第 1 列 → view，不拷贝
w = x.view(3, 2)     # reshape → view，不拷贝
t = x.T              # 转置 → view，不拷贝

# 验证共享 storage
assert x.storage().data_ptr() == y.storage().data_ptr()  # True！
```

**重要副作用**：修改 view 会修改原始 tensor！

```python
y = x[0]
y[0] = 999
print(x[0, 0])  # 999 ← x 也被改了！
```

### 5.3 Contiguous 与非 Contiguous

**Contiguous** = 按 storage 顺序遍历 tensor 时，地址连续递增。

```python
x = torch.tensor([[1, 2, 3], [4, 5, 6]])
x.is_contiguous()            # True, stride = (3, 1)

t = x.T                      # 转置
t.is_contiguous()            # False! stride 变成了 (1, 3)
# 遍历 t 时跳跃访问 storage：1, 4, 2, 5, 3, 6
```

**为什么重要**：非 contiguous tensor 不能直接 `view()`。

```python
t.view(6)  # ❌ RuntimeError!

# 解决方案 1：先 contiguous()（会拷贝！）
t.contiguous().view(6)  # ✅ 但 allocate 了新内存

# 解决方案 2：用 reshape()（必要时自动拷贝）
t.reshape(6)            # ✅ 等价于 contiguous().view()
```

### 5.4 总结：什么操作是 Free 的？

| 操作 | 是否拷贝 | 说明 |
|------|---------|------|
| `x[0]`，`x[:, 1]` | ❌ | 切片 = view |
| `x.view(...)` | ❌ | 仅改 metadata |
| `x.T` / `x.transpose(...)` | ❌ | 仅改 stride |
| `x.contiguous()` | ⚠️ 可能 | 已 contiguous 则不拷贝 |
| `x.reshape(...)` | ⚠️ 可能 | = contiguous() + view() |
| Element-wise ops (`+`, `*`, `triu`) | ✅ | 必须分配新 tensor |

> **实践建议**：View 操作是免费的，尽管用。但注意 `contiguous()` 和 `reshape()` 可能隐式拷贝。

---

## 六、Element-wise 操作与 MatMul

### 6.1 Element-wise 操作

创建新 tensor，每个元素独立计算：

```python
z = x + y        # 逐元素加法
z = x * y        # 逐元素乘法（Hadamard product）
z = torch.triu(x)  # 上三角（用于 causal attention mask）
```

### 6.2 矩阵乘法 (MatMul)：深度学习的心脏

```python
# 基本 matmul
x = torch.randn(16, 32)  # (M, K)
w = torch.randn(32, 2)   # (K, N)
y = x @ w                # (M, N) = (16, 2)
```

**FLOPs 计算**（记住这个公式）：

$$
\boxed{\text{FLOPs}_{\text{matmul}} = 2 \times M \times K \times N}
$$

- 因子 2 = 每个输出元素需要 $K$ 次乘法 + $K$ 次加法

### 6.3 Batched MatMul

实际训练中，MatMul 总是 batched 的（batch × sequence × hidden）：

```python
# Batched matmul
x = torch.randn(8, 64, 16, 32)  # (batch, seq, M, K)
w = torch.randn(32, 2)          # (K, N)
y = x @ w                       # (8, 64, 16, 2)
# 对每个 (batch, seq) 执行独立的 matmul
```

FLOPs = $2 \times 8 \times 64 \times 16 \times 32 \times 2$（所有维度的乘积 × 2）

### 6.4 关键直觉：MatMul 支配一切

> **深度学习中，没有任何其他操作在 FLOPs 上能与矩阵乘法相比。** 对于足够大的矩阵，其他操作（加法、激活函数、归一化）的 FLOPs 是线性的，而 MatMul 是立方的。这就是为什么 napkin math 只需要数 MatMul。

---

## 七、Einops 与 Einsum：告别索引地狱

### 7.1 问题：裸索引不可读

```python
# 💀 传统写法：三个月后自己都看不懂
scores = x @ y.transpose(-2, -1)  # -2 是啥？-1 是啥？
```

### 7.2 Einsum：带名字的矩阵乘法

```python
from torch import einsum

x: Float[Tensor, "batch seq1 hidden"]
y: Float[Tensor, "batch seq2 hidden"]

# ✅ Einsum：维度语义一目了然
scores = einsum(x, y, "batch seq1 hidden, batch seq2 hidden -> batch seq1 seq2")
# hidden 没出现在输出 → 被 sum 掉
# batch 出现在输出 → 被 iterate 掉
```

**规则**：
- 输出中 **存在** 的维度 → iterate（保留）
- 输出中 **不存在** 的维度 → sum（缩约）

### 7.3 Reduce：命名维度上的聚合

```python
from einops import reduce

x: Float[Tensor, "batch seq hidden"]

# 传统写法
mean = x.mean(dim=-1)

# Einops 写法
mean = reduce(x, "batch seq hidden -> batch seq", "sum")
# hidden 消失 → 在该维度上求和
```

### 7.4 Rearrange：拆分/合并维度

多头注意力中的经典场景：

```python
from einops import rearrange

# hidden = num_heads × head_dim，需要拆开
x: Float[Tensor, "batch seq hidden"]  # hidden = 8

x_heads = rearrange(x, "batch seq (heads hidden1) -> batch seq heads hidden1", heads=2)
# shape: (batch, seq, 2, 4)

# ... 对每个 head 独立计算 ...

# 合并回去
x_flat = rearrange(x_heads, "batch seq heads hidden1 -> batch seq (heads hidden1)")
```

### 7.5 Broadcasting 与省略号

```python
# 用 ... 代替任意数量的 batch 维度
scores = einsum(x, y, "... seq1 hidden, ... seq2 hidden -> ... seq1 seq2")
# 适用于 (batch, seq1, hidden) 也适用于 (batch1, batch2, seq1, hidden)
```

> **Percy 建议**：一开始可能觉得 einops 繁琐，但当模型变复杂时，它比 `-1, -2` 索引安全得多。`torch.compile` 也能正确优化 einsum。

---

## 八、FLOPs 审计：Compute Accounting

### 8.1 FLOP 的定义与易混淆点

| 写法 | 含义 | 用途 |
|------|------|------|
| FLOPs (小写 s) | Floating point **operations** (数量) | 衡量计算量 |
| FLOPS (大写 S) / flops/s | Floating point operations **per second** (速率) | 衡量硬件速度 |

> **Percy 的约定**：本课程用 flops 表示数量，flops/s 表示速率，避免歧义。

### 8.2 直觉：FLOPs 的量级

| 模型/事件 | FLOPs |
|-----------|-------|
| GPT-3 训练 | ~$3 \times 10^{23}$ |
| GPT-4 训练 (推测) | ~$2 \times 10^{25}$ |
| 美国行政令报告门槛 (已撤销) | $1 \times 10^{26}$ |
| 欧盟 AI 法案门槛 (仍有效) | $1 \times 10^{25}$ |

### 8.3 硬件 FLOPs/s 速率

| 硬件 | FP32 | BF16/FP16 | FP8 (无稀疏) | FP8 (2:4 稀疏) |
|------|------|-----------|-------------|----------------|
| A100 | 19.5 TFLOPS | 312 TFLOPS | N/A | N/A |
| H100 | 67 TFLOPS | ~990 TFLOPS | ~990 TFLOPS | ~1979 TFLOPS |

> ⚠️ **NVIDIA spec sheet 上的数字带 \*（稀疏）**。大多数 LLM 训练用的是 **dense** 矩阵，所以实际峰值是标注值的 **一半**。这是一个常见陷阱。

### 8.4 粗算：8×H100 一周能做多少计算？

$$
8 \times 990 \times 10^{12} \times 86400 \times 7 \approx 4.78 \times 10^{21} \text{ FLOPs (BF16)}
$$

对比：GPT-3 训练需要 $3 \times 10^{23}$ FLOPs → 8 卡大约需要 **63 周**（不现实，所以要上千卡）。

### 8.5 线性模型的 FLOPs 示例

```python
B, D, K = 1024, 512, 128
X = torch.randn(B, D, device='cuda')
W = torch.randn(D, K, device='cuda')
Y = X @ W  # (B, K)

flops = 2 * B * D * K  # = 2 × 1024 × 512 × 128 = 134,217,728
```

**推广到模型**：对于线性模型，$D \times K$ 就是参数量：

$$
\text{前向 FLOPs} = 2 \times N_{\text{tokens}} \times N_{\text{params}}
$$

> 这个公式对 Transformer 也近似成立（当 sequence length 不太长时），因为 Transformer 本质上是一堆 MatMul。

---

## 九、MFU：Model FLOPs Utilization

### 9.1 定义

$$
\boxed{\text{MFU} = \frac{\text{模型的理论 FLOPs} / \text{实际用时}}{\text{硬件承诺的 peak flops/s}}}
$$

或等价地：

$$
\text{MFU} = \frac{\text{实际 flops/s（模型视角）}}{\text{硬件 peak flops/s}}
$$

### 9.2 实测示例

**FP32 MatMul on H100**：

```
实际时间：0.16s
理论 FLOPs：2 × B × D × K
实际 flops/s：~5.4 × 10¹³
H100 FP32 peak：6.7 × 10¹³
MFU ≈ 0.8
```

**BF16 MatMul on H100**：

```
实际时间：0.03s（快了 5×!）
理论 FLOPs：相同
实际 flops/s：更高
H100 BF16 peak：~9.9 × 10¹⁴
MFU ≈ 偏低（因为 peak 很激进）
```

### 9.3 MFU 经验法则

| MFU | 评价 |
|-----|------|
| > 0.5 | 合格，大多数训练框架的目标 |
| 0.3 – 0.5 | 一般，有优化空间 |
| < 0.1 | 严重浪费硬件，需要排查 |
| > 0.8 | 很好，MatMul 高度主导 |

### 9.4 MFU vs HFU

- **MFU (Model FLOPs Utilization)**：以 **模型理论计算量** 为分子。不惩罚聪明的优化（如 activation checkpointing 导致的重复计算不算在内）。
- **HFU (Hardware FLOPs Utilization)**：以 **实际硬件执行的全部 FLOPs** 为分子。

> **Percy 的立场**：MFU 更有意义，因为它衡量的是"你的模型每秒训练了多少有效计算"，而不是硬件在忙什么。

---

## 十、反向传播的 FLOPs：为什么是 4×

### 10.1 两层线性模型

$$
X \xrightarrow{W_1} H_1 \xrightarrow{W_2} H_2 \rightarrow \mathcal{L}
$$

- $X$: $(B, D)$, $W_1$: $(D, D)$, $W_2$: $(D, K)$

**前向 FLOPs**：

$$
\text{Forward} = \underbrace{2 \times B \times D \times D}_{X \times W_1} + \underbrace{2 \times B \times D \times K}_{H_1 \times W_2}
$$

### 10.2 反向传播分析（以 $W_2$ 为例）

链式法则给出：

$$
\frac{\partial \mathcal{L}}{\partial W_2} = H_1^T \cdot \frac{\partial \mathcal{L}}{\partial H_2}
$$

这是一个 $(D, B) \times (B, K)$ 的 MatMul → FLOPs = $2 \times B \times D \times K$

但还需要继续反向传播到 $H_1$：

$$
\frac{\partial \mathcal{L}}{\partial H_1} = \frac{\partial \mathcal{L}}{\partial H_2} \cdot W_2^T
$$

这是一个 $(B, K) \times (K, D)$ 的 MatMul → FLOPs = $2 \times B \times D \times K$

所以对 $W_2$ 层：

$$
\text{Backward}_{W_2} = \underbrace{2BDK}_{\nabla W_2} + \underbrace{2BDK}_{\nabla H_1} = 4BDK
$$

### 10.3 总结：6× 法则

对每个线性层：
- **前向**：$2 \times B \times (\text{layer params})$
- **反向**：$4 \times B \times (\text{layer params})$（两个 MatMul：一个算参数梯度，一个传播到上层）

汇总：

$$
\boxed{\text{总 FLOPs} = 6 \times N_{\text{tokens}} \times N_{\text{params}}}
$$

| 阶段 | FLOPs | 占比 |
|------|-------|------|
| Forward | $2 \times B \times N_{\text{params}}$ | 1/3 |
| Backward (参数梯度) | $2 \times B \times N_{\text{params}}$ | 1/3 |
| Backward (激活梯度传播) | $2 \times B \times N_{\text{params}}$ | 1/3 |
| **总计** | $6 \times B \times N_{\text{params}}$ | 3/3 |

> **适用范围**：这个近似对大多数 Transformer 成立（当 attention 的 $O(n^2)$ 部分不占主导时）。不适用于参数共享严重的模型。

---

## 十一、模型构建

### 11.1 参数初始化：为什么不能用纯 randn

```python
D = 4096
x = torch.randn(1, D)
W = torch.randn(D, D)  # ← 标准正态初始化
y = x @ W
# y 的每个元素 ~ N(0, D)，标准差 ≈ √D ≈ 64
# 值会爆炸！多层堆叠后更严重
```

**Xavier Initialization（缩放修正）**：

$$
W \sim \mathcal{N}\left(0, \frac{1}{D_{\text{in}}}\right) \quad \text{即} \quad W = \frac{\text{randn}(D_{\text{in}}, D_{\text{out}})}{\sqrt{D_{\text{in}}}}
$$

```python
W = torch.randn(D, D) / (D ** 0.5)
y = x @ W
# y 的每个元素 ~ N(0, 1)，稳定！
```

**更保守的做法**：截断正态分布到 $[-3\sigma, 3\sigma]$，防止极端值。

### 11.2 nn.Parameter 与简单模型

```python
import torch.nn as nn

class LinearLayer(nn.Module):
    def __init__(self, d_in, d_out):
        super().__init__()
        self.W = nn.Parameter(torch.randn(d_in, d_out) / (d_in ** 0.5))

    def forward(self, x):
        return x @ self.W

class Cruncher(nn.Module):
    """深度线性网络 (教学用)"""
    def __init__(self, d, num_layers):
        super().__init__()
        self.layers = nn.ModuleList([LinearLayer(d, d) for _ in range(num_layers)])
        self.head = nn.Parameter(torch.randn(d) / (d ** 0.5))

    def forward(self, x):
        for layer in self.layers:
            x = layer(x)
        return x @ self.head  # → scalar per example

# 参数量
model = Cruncher(d=128, num_layers=2)
n_params = sum(p.numel() for p in model.parameters())
# = 128² + 128² + 128 = 32,896
```

### 11.3 随机性管理

随机性出现在：初始化、Dropout、数据 Loader 顺序。

**最佳实践**：

```python
# 为每个随机源设置独立 seed
torch.manual_seed(42)
torch.cuda.manual_seed_all(42)

# 调试时追求可重现
# 不同 seed 控制不同来源 → 可以固定初始化但变化数据顺序
```

> **Percy 的建议**：Determinism is your friend when debugging.

---

## 十二、数据加载

### 12.1 语言模型的数据格式

Tokenizer 输出 = 整数序列 → 序列化为 NumPy 数组存磁盘。

### 12.2 Memory-Mapped Files：不加载全部数据

```python
import numpy as np

# LLaMA 数据约 2.8TB，不可能全部加载到 RAM
data = np.memmap("tokens.bin", dtype=np.uint16, mode='r')
# data 看起来像普通数组，但实际按需从磁盘加载
```

---

## 十三、优化器

### 13.1 优化器简史

```
SGD           → 最朴素：θ -= lr * ∇L
Momentum      → 维护梯度的指数移动平均 (EMA)
AdaGrad       → 按参数缩放：除以历史梯度平方和的根号
RMSprop       → AdaGrad 的改进：用 EMA 替代全量平均
Adam (2014)   → RMSprop + Momentum = 一阶矩 + 二阶矩
AdamW         → Adam + 解耦权重衰减
```

### 13.2 手写 AdaGrad（理解优化器结构）

```python
class SimpleAdaGrad(torch.optim.Optimizer):
    def __init__(self, params, lr=0.01):
        super().__init__(params, defaults=dict(lr=lr))

    def step(self):
        for group in self.param_groups:
            for p in group['params']:
                if p.grad is None:
                    continue
                grad = p.grad.data
                state = self.state[p]

                # 初始化累积梯度平方
                if 'grad_squared_sum' not in state:
                    state['grad_squared_sum'] = torch.zeros_like(p.data)

                # 累积 g²
                state['grad_squared_sum'].add_(grad ** 2)

                # 更新参数：θ -= lr * g / √(Σg²)
                p.data.add_(
                    -group['lr'] * grad / (state['grad_squared_sum'].sqrt() + 1e-10)
                )

                # 释放梯度内存
                p.grad = None
```

### 13.3 AdamW 每个参数的存储需求

| 存储项 | 内容 | 精度 | Bytes/param |
|--------|------|------|-------------|
| 参数 (θ) | 模型权重 | FP32 | 4 |
| 梯度 (∇θ) | 当前 batch 的梯度 | FP32 | 4 |
| 一阶矩 (m) | 梯度的 EMA | FP32 | 4 |
| 二阶矩 (v) | 梯度² 的 EMA | FP32 | 4 |
| **总计** | | | **16** |

> **这就是开头问题 2 中 "16 bytes per parameter" 的来源！**

---

## 十四、内存审计：全景

### 14.1 内存的四大户

$$
\text{总内存} = \underbrace{N_p \times b_p}_{\text{Parameters}} + \underbrace{N_a \times b_a}_{\text{Activations}} + \underbrace{N_p \times b_g}_{\text{Gradients}} + \underbrace{N_o \times b_o}_{\text{Optimizer States}}
$$

对于两层线性模型（$D=128$, $B=8$, $K=1$, AdaGrad, FP32）：

| 类别 | 计算 | 数量 |
|------|------|------|
| Parameters | $D^2 \times 2 + D = 32896$ | 32,896 |
| Activations | $B \times D \times 2 = 2048$ | 2,048 |
| Gradients | = Parameters = 32,896 | 32,896 |
| Optimizer (AdaGrad, 1 copy) | = Parameters = 32,896 | 32,896 |

$$
\text{总内存} = 4 \times (32896 + 2048 + 32896 + 32896) = 4 \times 100736 = 402944 \text{ bytes}
$$

### 14.2 对 Transformer 的推广

同样的框架，但维度更复杂：

- **Parameters**：Attention 的 $W_Q, W_K, W_V, W_O$ + FFN 的 $W_1, W_2$ × 层数
- **Activations**：取决于 batch size、sequence length、hidden dim
- **Gradients**：= Parameters
- **Optimizer (AdamW)**：= 2 × Parameters（m 和 v 各一份）

$$
\text{最低内存 (AdamW, FP32)} \approx 4 \times (1 + 1 + 2) \times N_{\text{params}} = 16 \times N_{\text{params}} \text{ bytes}
$$

> Activations 是 **变量**（依赖 batch size），可以通过 **Activation Checkpointing** 来用计算换内存（后续课程讲）。

---

## 十五、训练循环与检查点

### 15.1 标准训练循环

```python
model = Cruncher(d=128, num_layers=2).cuda()
optimizer = SimpleAdaGrad(model.parameters(), lr=0.01)

for step in range(num_steps):
    x, y = get_batch()
    pred = model(x)
    loss = mse_loss(pred, y)
    loss.backward()
    optimizer.step()
```

### 15.2 检查点：不要失去进度

> 训练大模型一定会崩溃。**一定。**

保存内容 = 模型状态 + 优化器状态 + 当前 step：

```python
# 保存
torch.save({
    'model': model.state_dict(),
    'optimizer': optimizer.state_dict(),
    'step': step,
}, 'checkpoint.pt')

# 恢复
ckpt = torch.load('checkpoint.pt')
model.load_state_dict(ckpt['model'])
optimizer.load_state_dict(ckpt['optimizer'])
step = ckpt['step']
```

> ⚠️ 必须同时保存 optimizer state！否则恢复后 Adam 的 m, v 会从零开始，训练会出问题。

---

## 十六、总结：本讲的知识框架

```
Lecture 2 知识树
│
├── Mechanics (机制)
│   ├── Tensor: storage + stride + view vs copy
│   ├── 浮点类型: FP32 / FP16 / BF16 / FP8
│   ├── GPU 内存管理: .to('cuda'), memory_allocated
│   ├── Einops / Einsum: 命名维度
│   ├── 优化器实现: param_groups, state dict
│   └── 训练循环 + 检查点
│
├── Mindset (心态)
│   ├── ⭐ 写代码前先算内存和 FLOPs
│   ├── MatMul FLOPs = 2 × M × K × N
│   ├── 总训练 FLOPs ≈ 6 × params × tokens
│   ├── 每参数 16 bytes (AdamW, FP32)
│   ├── MFU > 0.5 才算合格
│   └── Benchmark，不要信 spec sheet
│
└── Intuitions (直觉)
    ├── BF16 是当前训练的甜蜜点
    ├── 混合精度：临时计算用低精度，积累用高精度
    ├── Xavier Init: 缩放防爆炸
    └── 硬件在推动模型设计（量化 → 架构协同设计）
```

### 从 Lecture 1 到 Lecture 2 的连接

| Lecture 1 | Lecture 2 |
|-----------|-----------|
| "效率是核心约束" | 具体化：内存审计 + FLOPs 审计 |
| "BPE 压缩比影响效率" | Token 数直接进入 6× 公式 |
| "三种知识：Mechanics / Mindset / Intuitions" | 本讲：Mechanics = PyTorch, Mindset = 资源审计 |
| Assignment 1 = BPE + Transformer + AdamW | 本讲给出了所有需要的 napkin math 工具 |

---

## 延伸阅读

- **Transformer FLOPs 详细推导**：课件中提到的参考文章（用于 Assignment 1 的 Transformer 资源审计）
- **混合精度训练**：Micikevicius et al., 2017, *"Mixed Precision Training"*
- **einops 教程**：[einops tutorial](https://einops.rocks/)
- **JAXTyping**：给 Tensor shape 加类型标注
