## CS336 第六讲：Kernels & Triton — 为 GPU 编写高性能代码

> "If you want to write high performance code, you should remember to **benchmark and profile** your code."

### 本讲核心要点

**一句话总结**：通过 Benchmarking、Profiling 和 Kernel Fusion（CUDA / Triton / torch.compile），将逐元素操作从 8ms 优化到 ~1ms。

**核心技能清单：**

| 技能                        | 关键点                                             |
| --------------------------- | -------------------------------------------------- |
| Benchmarking                | warmup + `torch.cuda.synchronize()` + 多次取均值   |
| Profiling（PyTorch）        | `torch.profiler` 看 CPU/CUDA 各操作耗时            |
| Profiling（Nsight Systems） | 可视化 CPU/GPU 时间线、kernel 调度、异步执行       |
| CUDA C++ Kernel             | 线程级编程，手动管理坐标与边界                     |
| Triton Kernel               | Block 级编程，Python 语法，编译器管理线程/合并访存 |
| torch.compile               | JIT 编译器自动融合算子，零代码改动                 |

**性能对比（GeLU，大向量）：**

| 实现方式               | 耗时    | 相对加速 |
| ---------------------- | ------- | -------- |
| 手写 PyTorch（逐操作） | 8.1 ms  | 1×       |
| PyTorch `F.gelu`       | 1.1 ms  | 7.4×     |
| CUDA C++ kernel        | 1.84 ms | 4.4×     |
| Triton kernel          | 1.85 ms | 4.4×     |
| `torch.compile`        | 1.47 ms | 5.5×     |

---

### 一、GPU 架构快速回顾

> 详细内容见第五讲笔记，这里只回顾本讲需要用到的关键概念。

#### 1.1 硬件结构

```
GPU（如 A100 / H100）
├── SM × 108+  (Streaming Multiprocessor)
│   ├── INT32 单元、FP32 单元、Tensor Core
│   ├── 寄存器文件（Register File）  ← 极快，每个线程私有
│   ├── 共享内存 / L1 Cache         ← 很快，Block 内线程共享
│   └── Warp 调度器
└── 全局内存(HBM/DRAM)              ← 大但慢
```

#### 1.2 执行模型

```
Grid（网格）
  └── Thread Block（线程块）→ 被调度到 一个 SM
        └── Warp（束）= 32 个线程 → 同时执行同一条指令（SIMT）
              └── Thread（线程）→ 最小执行单元
```

**关键约束：**

- **Block 内**线程可通过共享内存通信 → 快（L1 级别）
- **Block 间**无法直接同步或共享数据 → 慢（必须走全局内存）
- 理想情况：Thread Block 数远多于 SM 数，使每个 **wave**（一轮调度）负载均衡

#### 1.3 Warp 的存在意义

> 课堂问答：What's the function of Warp?

Warp 是 **32 个线程同时执行** 的最小调度单位。它存在的原因是**减少控制逻辑开销**：

- GPU 的计算单元 >> Warp 调度器数量
- 不需要为每个线程维护独立的控制逻辑
- 同一 Warp 内的线程执行相同指令 → 极少的控制硬件即可驱动大量并行计算
- **这也是 CPU 和 GPU 的核心设计取舍**：CPU 将大量硅片面积用于分支预测、乱序执行等控制逻辑；GPU 将面积用于计算单元，控制逻辑极简

#### 1.4 算术强度（Arithmetic Intensity）

$$\text{Arithmetic Intensity} = \frac{\text{FLOPs}}{\text{Bytes moved}}$$

- **目标**：算术强度尽量高（更多计算、更少搬数据）
- **原因**：计算能力的增长速度远超内存带宽的增长速度 → 大部分操作是 **memory bound**
- **经验法则**：
  - 矩阵乘法 → **compute bound**（做得好的话）
  - 其他所有操作 → **memory bound** → 优化重点应在减少内存搬运

---

### 二、Benchmarking — 测量端到端性能

#### 2.1 为什么要 Benchmark

> "I've seen students spend three hours optimizing something that turns out wasn't the bottleneck."

- 比较不同实现（手写 CUDA vs Triton vs torch.compile）的性能
- 理解操作随规模增长的性价比变化
- **理论有限，最终必须做端到端实测**（库版本、硬件差异、微架构行为都无法完全预测）

#### 2.2 Benchmark 函数的关键细节

```python
def benchmark(fn, num_warmup=5, num_trials=20):
    # 1️⃣ Warmup：排除首次编译/初始化开销
    for _ in range(num_warmup):
        fn()

    # 2️⃣ 同步 CPU 和 GPU 状态
    torch.cuda.synchronize()

    times = []
    for _ in range(num_trials):
        start = time.time()
        fn()
        # 3️⃣ 每次运行后同步，确保测到 GPU 真实执行时间
        torch.cuda.synchronize()
        times.append(time.time() - start)

    return sum(times) / len(times)
```

**两个关键陷阱：**

| 陷阱                                | 问题                                                                   | 解决                   |
| ----------------------------------- | ---------------------------------------------------------------------- | ---------------------- |
| **不做 Warmup**                     | 首次运行含 JIT 编译、内存分配等一次性开销                              | 跑几轮 warmup 后再计时 |
| **不调 `torch.cuda.synchronize()`** | CPU 异步发射 kernel 后立即返回，测到的是 CPU 发射时间而非 GPU 执行时间 | 每次计时前后都同步     |

> 忘了 synchronize 会怎样？你会发现 "大矩阵乘法瞬间完成" — 这显然不对。

#### 2.3 Benchmark 结果示例

**矩阵乘法随维度增长：**

- 小矩阵（1024, 2048）：时间几乎不变 → 常数开销主导（kernel launch、数据传输）
- 大矩阵：超线性增长（$O(N^3)$ 的矩阵乘法）

**MLP 随 steps/layers 增长：**

- Steps 增加 → 线性增长（每步是独立的前向+反向）
- Layers 增加 → 线性增长（每层是相同维度的 Linear + GeLU）

---

### 三、Profiling — 定位瓶颈在哪里

#### 3.1 PyTorch 内置 Profiler

```python
with torch.profiler.profile(
    activities=[
        torch.profiler.ProfilerActivity.CPU,
        torch.profiler.ProfilerActivity.CUDA,
    ]
) as prof:
    torch.cuda.synchronize()
    fn()
    torch.cuda.synchronize()

print(prof.key_averages().table(sort_by="self_cpu_time_total"))
```

**输出解读：**

Profiler 输出一张表，展示每个操作的 CPU 时间和 CUDA 时间占比。

#### 3.2 逐操作分析

**向量加法 (`A + B`, 2048×2048)：**

```
Python:  A + B
  └── aten::add               ← C++ 接口层，CPU 开销大（~98% CPU time）
       └── vectorized_elementwise_kernel<AddFunctor>  ← 实际 GPU kernel
  └── cudaLaunchKernel        ← CPU 发射 kernel 到 GPU
  └── cudaDeviceSynchronize   ← 等待 GPU 完成
```

- CPU 时间：~1.4 ms（主要是 C++ dispatch 层的开销）
- CUDA 时间：~17 μs（GPU 执行极快）

**矩阵乘法 (`A @ B`, 2048×2048)：**

```
Python:  A @ B
  └── aten::mm
       └── CUTLASS kernel（高性能矩阵乘法库）  ← 78%+ CUDA 时间
  └── cudaLaunchKernel
  └── cudaDeviceSynchronize
```

**矩阵乘法 (`A @ B`, 128×128, 小矩阵)：**

```
Python:  A @ B
  └── aten::mm
       └── xmma_gemm（不同于大矩阵的 kernel！）
```

> **关键洞察**：PyTorch 会根据矩阵维度和硬件自动选择不同的底层 kernel。
> `torch.compile` 还可以 micro-benchmark 不同 kernel 选择最优的 → 免费 ~10% 加速。

#### 3.3 复合操作分析

**`torch.cdist`（欧氏距离矩阵）：**

```
cdist 分解为:
├── aten::mm       → gemm kernel      (78% CUDA time)
├── aten::cat      → copy kernel      (6%)
├── aten::pow      → elementwise pow  (5%)
└── aten::sum      → reduction sum    (3%)
```

→ 通过 profiler 看到 78% 时间在矩阵乘法 → 优化重点清晰。

**`F.gelu` 和 `F.softmax`：**

- 这些核心操作在 PyTorch 中已有 **专门的融合 kernel**
- 不是由基本操作组合而成，而是一个 kernel 搞定全部
- 所以 `F.gelu` 比手写 `0.5 * x * (1 + tanh(...))` 快很多

#### 3.4 NVIDIA Nsight Systems — 更强大的 Profiler

PyTorch profiler 适合快速查看，但对复杂模型不够用（如 MLP 的 forward+backward 只能显示前 10 个操作）。

Nsight Systems 提供 **时间线视图**，可以同时看到：

```
┌─ CUDA HW ────────────────────────────────────────────────────┐
│  [kernel][kernel][kernel][kernel][kernel]...                  │  ← GPU 实际执行的 kernel
├─ GPU Memory ─────────────────────────────────────────────────┤
│  ▁▂▃▄▅▆▇████████████████████████████████                     │  ← 显存使用曲线
├─ NVTX Annotations ───────────────────────────────────────────┤
│  |define_model|  |step_0 |step_1|step_2|step_3|step_4|step_5|│  ← 代码标注
├─ CPU Threads ────────────────────────────────────────────────┤
│  [dispatch][dispatch][dispatch]...[wait]...[dispatch]...      │  ← CPU 发射 kernel
└──────────────────────────────────────────────────────────────┘
     时间轴 →
```

**代码中添加 NVTX 标注：**

```python
import torch.cuda.nvtx as nvtx

nvtx.range_push("define_model")
model = MLP(...)
nvtx.range_pop()

for i in range(num_steps):
    nvtx.range_push(f"step_{i}")
    # forward + backward
    nvtx.range_pop()
```

→ Nsight Systems 会在时间线上标出每个代码块对应的区间。

#### 3.5 CPU 与 GPU 的异步执行模型

> 这是本讲最深刻的洞察之一。

```
CPU:  [dispatch layer0][dispatch layer1]...[dispatch layer9][wait...][dispatch step1]...
       ↓               ↓                    ↓
GPU:               [exec layer0      ][exec layer1      ]...[exec layer9      ]

CPU 发射 kernel 远快于 GPU 执行 → CPU 跑在 GPU 前面好几个 step！
```

**异步执行的核心原理：**

1. CPU 执行 Python 代码，遇到 GPU 操作时将 kernel **排入队列**
2. CPU 不等 GPU 完成就继续执行下一行 Python 代码
3. GPU 从队列中取 kernel 逐个执行
4. CPU 可以领先 GPU **一整个 step** 甚至更多

**这也解释了为什么 Python 慢不重要：**

> CPU 不是瓶颈，它只是一个 "kernel 发射器"。GPU 执行才是耗时大头。Python 的速度完全足够应付 kernel dispatch。

#### 3.6 `print(loss)` 引发的性能悬崖

```python
for step in range(num_steps):
    loss = model(x).mean()
    loss.backward()
    print(loss.item())  # ← 这行看似无害！
```

**问题**：`loss.item()` 需要将数据从 GPU 拷回 CPU → 触发 `cudaStreamSynchronize` → **CPU 被迫等待 GPU 完成**

```
无 print:
  CPU: [d0][d1][d2][d3][d4][d5]...        ← CPU 一路狂奔
  GPU:      [e0   ][e1   ][e2   ][e3   ]  ← GPU 连续执行

有 print:
  CPU: [d0][wait~~~~~~~~][print][d1][wait~~~~~~~~][print][d2]...
  GPU:      [e0         ]            [e1         ]            ← GPU 间歇停顿
```

- 无 print：CPU 领先 GPU 一个完整 step，GPU 保持满载
- 有 print：每个 step 结束 CPU 必须等 GPU → GPU 利用率下降
- **极端情况**（频繁 print）：CPU 完全成为瓶颈，GPU 大量空闲

> **教训**：训练循环中避免频繁的 CPU-GPU 同步。logging 可以每 N 步做一次。

---

### 四、Kernel Fusion — 为什么手写 PyTorch 慢？

#### 4.1 问题的直觉

```
【未融合的 GeLU（手写 PyTorch）】

每一步：HBM → SM → 计算 → SM → HBM → SM → 计算 → SM → HBM → ...
                                    ↑
                              反复搬运数据！

步骤：x³ → 乘以 0.044715 → 加上 x → 乘以 √(2/π) → tanh → 加 1 → 乘以 0.5 → 乘以 x
     ↑每一步都是一个独立的 CUDA kernel，每次都从 HBM 读一遍、写一遍

【融合的 GeLU（单 kernel）】

一次：HBM → SM → 全部计算 → SM → HBM

只读一次、写一次！
```

#### 4.2 Profiler 验证

**手写 GeLU：**

```
╔═══════════════════════════════════╗
║ vectorized_elementwise(mul) ×3    ║  ← 多次 kernel launch
║ vectorized_elementwise(add) ×2    ║
║ vectorized_elementwise(tanh) ×1   ║
╚═══════════════════════════════════╝
耗时：~8.1 ms
```

**PyTorch `F.gelu`：**

```
╔═══════════════════════════════════╗
║ gelu_cuda_kernel                  ║  ← 单次 kernel launch
╚═══════════════════════════════════╝
耗时：~1.1 ms
```

→ 8× 差距完全来自内存往返次数。

---

### 五、CUDA C++ Kernel — 线程级编程

#### 5.1 编程模型

```
Grid（所有工作）
├── Block 0 → SM 0
│   ├── Thread 0  (blockIdx.x * blockDim.x + threadIdx.x = 0)
│   ├── Thread 1  (blockIdx.x * blockDim.x + threadIdx.x = 1)
│   └── ...
├── Block 1 → SM 1
│   ├── Thread 0  (blockIdx.x * blockDim.x + threadIdx.x = 1024)
│   └── ...
└── Block N-1
```

**坐标计算：**

$$i = \texttt{blockIdx.x} \times \texttt{blockDim.x} + \texttt{threadIdx.x}$$

- `blockIdx.x`：当前 Block 在 Grid 中的索引
- `blockDim.x`：每个 Block 的线程数（= BLOCK_SIZE）
- `threadIdx.x`：当前线程在 Block 内的索引

#### 5.2 完整 GeLU CUDA Kernel

**Kernel 函数（运行在 GPU 上）：**

```cpp
// .cu 文件
__global__ void geluKernel(const float* in, float* out, int num_elements) {
    // 1. 计算全局坐标
    int i = blockIdx.x * blockDim.x + threadIdx.x;

    // 2. 边界检查（最后一个 Block 可能越界）
    if (i < num_elements) {
        float x = in[i];
        // 3. GeLU 公式
        float cdf = 0.5f * (1.0f + tanhf(
            sqrtf(2.0f / M_PI) * (x + 0.044715f * x * x * x)
        ));
        out[i] = x * cdf;
    }
}
```

**Wrapper 函数（运行在 CPU 上，负责编排和启动 kernel）：**

```cpp
torch::Tensor gelu(torch::Tensor x) {
    // 检查输入合法性
    TORCH_CHECK(x.device().is_cuda(), "x must be on CUDA");
    TORCH_CHECK(x.is_contiguous(), "x must be contiguous");

    auto y = torch::empty_like(x);  // 分配输出（不初始化为 0）

    int num_elements = x.numel();
    int BLOCK_SIZE = 1024;
    int num_blocks = (num_elements + BLOCK_SIZE - 1) / BLOCK_SIZE;  // ceil division

    // 启动 kernel: <<<num_blocks, BLOCK_SIZE>>>
    geluKernel<<<num_blocks, BLOCK_SIZE>>>(
        x.data_ptr<float>(), y.data_ptr<float>(), num_elements
    );

    return y;
}
```

#### 5.3 关键设计要点

| 要点                           | 说明                                         |
| ------------------------------ | -------------------------------------------- |
| `__global__`                   | 标记为 CUDA kernel 函数                      |
| `is_contiguous()` 检查         | 确保内存连续，否则下标偏移算术无效           |
| `empty_like` 而非 `zeros_like` | 避免不必要的初始化 → 省一次内存写            |
| ceil division 计算 Block 数    | 确保最后一批元素也能被处理                   |
| `if (i < num_elements)`        | 最后一个 Block 的部分线程超出数组 → 必须跳过 |
| `x.data_ptr<float>()`          | 传递原始指针而非 Tensor 对象                 |
| `.cu` 文件扩展名               | CUDA C++ 源文件的约定                        |

> 课堂问答：不连续内存何时出现？`transpose()`、`view()`、维度置换等操作可能导致。
> 解决方法：在 wrapper 中调用 `.contiguous()` 强制连续化。

#### 5.4 调试技巧

```bash
CUDA_LAUNCH_BLOCKING=1 python my_script.py
```

→ 强制 CUDA kernel 同步执行，使错误信息能正确定位到出错的 kernel。

#### 5.5 Python 中加载 CUDA Kernel

```python
from torch.utils.cpp_extension import load_inline

cuda_gelu = load_inline(
    name='cuda_gelu',
    cpp_sources=wrapper_code,
    cuda_sources=kernel_code,
    functions=['gelu'],
)

# 像普通 Python 函数一样调用
y = cuda_gelu.gelu(x)
```

#### 5.6 性能

- **耗时**：1.84 ms（vs 手写 PyTorch 8.1 ms → **4.4× 加速**）
- **Profiler 确认**：单个 CUDA kernel，100% GPU 时间 → 融合成功
- 比 PyTorch 内置实现（1.1 ms）还慢一点 → PyTorch 做了更多针对性优化

---

### 六、Triton Kernel — Block 级编程

#### 6.1 Triton 是什么

| 特性     | 说明                                                   |
| -------- | ------------------------------------------------------ |
| 开发者   | OpenAI，2021 年发布                                    |
| 编程模型 | **Block 级**（不用想线程，只想 Block）                 |
| 语言     | Python                                                 |
| 自动管理 | 内存合并（coalescing）、共享内存管理、Block 内线程同步 |
| 手动管理 | SM 间调度、Block 间通信                                |
| 优势     | 接近 CUDA 性能，但代码量大幅减少，可调试               |

**CUDA vs Triton 对比：**

```
CUDA:  每个线程处理 1 个元素，你管理线程坐标
Triton: 每个 Block 处理一段向量，你管理 Block 坐标
         编译器自动把向量操作映射到线程上
```

#### 6.2 Triton GeLU Kernel

```python
import triton
import triton.language as tl

@triton.jit
def triton_gelu_kernel(
    x_ptr, y_ptr, block_size: tl.constexpr, num_elements
):
    # 1. 计算当前 Block 的起始位置
    block_start = tl.program_id(0) * block_size

    # 2. 生成 Block 内所有偏移量（向量化！）
    offsets = block_start + tl.arange(0, block_size)

    # 3. 越界 mask
    mask = offsets < num_elements

    # 4. 一次性加载整个 Block 的数据
    x = tl.load(x_ptr + offsets, mask=mask)

    # 5. 向量化计算 GeLU（看起来就像普通 Python！）
    # tanh 近似: tanh(x) ≈ 1 - 2/(1 + exp(2x))
    inner = 0.7978845608 * (x + 0.044715 * x * x * x)  # √(2/π) * (x + 0.044715x³)
    y = 0.5 * x * (1.0 + tl.libdevice.tanh(inner))

    # 6. 写回结果
    tl.store(y_ptr + offsets, y, mask=mask)
```

**Wrapper：**

```python
def triton_gelu(x: torch.Tensor):
    assert x.is_cuda and x.is_contiguous()
    y = torch.empty_like(x)
    num_elements = x.numel()
    BLOCK_SIZE = 1024
    num_blocks = (num_elements + BLOCK_SIZE - 1) // BLOCK_SIZE
    triton_gelu_kernel[(num_blocks,)](x, y, BLOCK_SIZE, num_elements)
    return y
```

#### 6.3 Triton vs CUDA 的关键差异

| 方面     | CUDA C++                          | Triton                                                     |
| -------- | --------------------------------- | ---------------------------------------------------------- |
| 编程视角 | 单个线程                          | 单个 Block                                                 |
| 坐标计算 | `blockIdx * blockDim + threadIdx` | `tl.program_id(0) * block_size + tl.arange(0, block_size)` |
| 数据加载 | `in[i]`（一个值）                 | `tl.load(ptr + offsets, mask=mask)`（一个向量）            |
| 计算     | 标量操作                          | 向量操作（编译器分配给线程）                               |
| 数据写回 | `out[i] = val`                    | `tl.store(ptr + offsets, val, mask=mask)`                  |
| 内存合并 | 手动确保                          | 编译器自动处理                                             |
| 共享内存 | 手动管理                          | 编译器自动管理                                             |

#### 6.4 PTX 机器码分析

Triton 编译后生成 PTX（接近 GPU 机器码的汇编）。分析 GeLU 的 PTX 可以看到：

```ptx
// 1. 寄存器声明
.reg .b32  %r<50>;    // 32-bit 通用寄存器
.reg .f32  %f<45>;    // 32-bit 浮点寄存器
.reg .b64  %rd<10>;   // 64-bit 寄存器（地址用）

// 2. 加载函数参数（X指针、Y指针等）
ld.param.u64 %rd1, [x_ptr];
ld.param.u64 %rd4, [y_ptr];

// 3. 计算 Block 内坐标偏移

// 4. 从全局内存加载数据（一次加载 4 个 float → 利用 burst mode！）
ld.global.v4.b32 {%r2, %r3, %r4, %r5}, [%rd1];

// 5. 大量浮点运算（mul, add, ex2 等）
mul.f32  %f10, %f1, %f1;      // x²
mul.f32  %f11, %f10, %f1;     // x³
mul.f32  %f12, %f11, 0x3D372713; // x³ × 0.044715

// ... tanh 用 exp 实现: tanh(x) = (e^2x - 1) / (e^2x + 1)
// GPU 没有 e^x 指令，用 2^x + 换底: e^x = 2^(x·log₂e)
ex2.approx.f32 %f25, %f24;    // 2^x 指令

// 6. 存储结果（一次写 4 个 float）
st.global.v4.b32 [%rd4], {%r38, %r39, %r40, %r41};
```

**关键观察：**

- **每个线程处理 4 个元素**：利用 burst mode，一次读写 4 个 float
- **所有中间结果在寄存器中**：极快，无需访问全局内存
- **GPU 没有 $e^x$ 指令**：用 $2^x$ 配合换底公式 $e^x = 2^{x \cdot \log_2 e}$ 实现
- **Triton 编译器自动做了内存合并**：`ld.global.v4` 表示 4-wide 向量化加载

#### 6.5 性能

- **耗时**：1.85 ms（基本与 CUDA C++ 持平）
- **代码量**：比 CUDA 少得多，全程 Python
- **Profiler**：单个 kernel launch，100% GPU 时间

---

### 七、torch.compile — 零成本优化

#### 7.1 使用方法

```python
# 就这一行
compiled_gelu = torch.compile(manual_gelu)

# 正常调用
y = compiled_gelu(x)
```

#### 7.2 它做了什么

1. **Tracing**：跟踪 Python 函数的计算图
2. **算子融合**：将多个小操作合并为一个 kernel
3. **代码生成**：底层生成 **Triton 代码**并编译
4. **矩阵乘法调优**：如果知道矩阵形状，micro-benchmark 不同 kernel 选最优的

```
torch.compile 下的 GeLU:
  → 生成 fused_add_mul_tanh Triton kernel
  → 比我们手写的 Triton 稍微更优化（编译器做了额外优化）
```

#### 7.3 性能

- **耗时**：1.47 ms（比我们手写的 Triton 1.85ms 还快！）
- **原因**：编译器做了额外的微优化（寄存器分配、指令调度等）
- **Profiler**：单个融合 kernel，100% GPU 时间

#### 7.4 何时用 torch.compile，何时手写 Kernel？

| 场景                               | 建议                              |
| ---------------------------------- | --------------------------------- |
| 简单算子融合（elementwise ops）    | `torch.compile` 完全够用          |
| 矩阵乘法调优                       | `torch.compile` 自动选最优 kernel |
| 标准操作组合                       | `torch.compile` 搞定              |
| Flash Attention 级别的复杂优化     | 手写 Triton                       |
| 利用特定硬件特性（如 H100 新指令） | 手写 CUDA/Triton                  |
| 全新架构中的复杂非标准操作         | 手写 Triton                       |

> "You shouldn't go home and say I'm gonna write CUDA kernels for every single part of my language model. That's probably not a good use of your time."

---

### 八、Triton Softmax — 带 Reduction 的 Kernel

#### 8.1 与 GeLU 的区别

- GeLU：**逐元素操作**，完全独立，最简单
- Softmax：有 **reduction**（求 max、求 sum），需要行内所有元素协作

#### 8.2 设计策略

```
关键决策：每个 Block 负责一整行

Matrix (M × N):
  Block 0 → Row 0: [x₀₀, x₀₁, ..., x₀ₙ]  → 在 SM 内完成 max, sum, normaliz
  Block 1 → Row 1: [x₁₀, x₁₁, ..., x₁ₙ]
  ...
  Block M → Row M: [xₘ₀, xₘ₁, ..., xₘₙ]

num_blocks = num_rows
block_size = next_power_of_2(num_cols)  ← padding 到 2 的幂
```

**为什么这样设计？**

- Softmax 需要对每一 **行** 做归一化
- 如果整行能装进一个 SM → 所有 reduction 都在共享内存/寄存器中完成 → 无需跨 Block 通信 → 极快

#### 8.3 Triton Softmax Kernel

```python
@triton.jit
def softmax_kernel(
    x_ptr, y_ptr, x_stride, y_stride, n_cols,
    BLOCK_SIZE: tl.constexpr
):
    # 哪一行
    row_idx = tl.program_id(0)

    # 该行的所有列偏移
    col_offsets = tl.arange(0, BLOCK_SIZE)
    mask = col_offsets < n_cols

    # 加载整行到 SM 本地内存
    row = tl.load(x_ptr + row_idx * x_stride + col_offsets, mask=mask, other=-float('inf'))

    # 标准 softmax: exp(x - max) / sum(exp(x - max))
    row_max = tl.max(row, axis=0)
    numerator = tl.exp(row - row_max)
    denominator = tl.sum(numerator, axis=0)
    result = numerator / denominator

    # 写回
    tl.store(y_ptr + row_idx * y_stride + col_offsets, result, mask=mask)
```

#### 8.4 Wrapper

```python
def triton_softmax(x: torch.Tensor):
    assert x.is_cuda and x.is_contiguous()
    M, N = x.shape
    y = torch.empty_like(x)
    BLOCK_SIZE = triton.next_power_of_2(N)  # padding
    num_blocks = M  # 每行一个 Block
    softmax_kernel[(num_blocks,)](
        x, y, x.stride(0), y.stride(0), N, BLOCK_SIZE=BLOCK_SIZE
    )
    return y
```

#### 8.5 Softmax 性能对比

| 实现方式                               | 耗时   |
| -------------------------------------- | ------ |
| 手写 PyTorch（逐操作 max/exp/sum/div） | 3.7 ms |
| PyTorch `F.softmax`                    | 1.5 ms |
| `torch.compile`                        | 1.3 ms |
| Triton kernel                          | 1.9 ms |

- `torch.compile` 居然**比 PyTorch 原生实现还快** → 因为它知道形状，能选更优的 kernel
- 手写 Triton 稍慢 → 还有优化空间（如分块策略、warp-level reduction 等）

---

### 九、总结与实践指南

#### 9.1 性能优化工作流

```
 1. 先跑起来 → 用 PyTorch 写正确的代码
                ↓
 2. Benchmark → 测端到端耗时，发现"慢"
                ↓
 3. Profile  → 用 PyTorch profiler 或 Nsight Systems 定位瓶颈
                ↓
 4. 先试 torch.compile → 零代码改动，看能否解决
                ↓
 5. 仍不够？→ 写 Triton kernel（Block 级，Python 语法）
                ↓
 6. 极端性能需求 → 写 CUDA C++ kernel（线程级，硬件控制最细）
```

#### 9.2 永远记住的原则

| 原则                     | 解释                                 |
| ------------------------ | ------------------------------------ |
| **先 Profile 再优化**    | 不要猜瓶颈在哪                       |
| **Warmup + Synchronize** | Benchmark 的两个铁律                 |
| **减少内存搬运**         | Kernel fusion 的核心目的             |
| **CPU 和 GPU 是异步的**  | 理解 dispatch 模型才能正确测时       |
| **print/logging 会同步** | 训练循环中避免每步都访问 GPU 数据    |
| **不要过度手写 kernel**  | torch.compile 能覆盖大部分场景       |
| **Python 慢不重要**      | CPU 只是 kernel 发射器，GPU 才是瓶颈 |
