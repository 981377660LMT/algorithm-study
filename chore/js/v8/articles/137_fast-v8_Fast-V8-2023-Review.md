# Fast V8: V8 2023 性能综述 | Fast V8: A Review of V8 Performance in 2023

### 背景：2023 的核心议题 —— “极致速度与体系化协同”

2023 年，V8 团队的核心议题是**在新的行业基准（Speedometer 3）下追求极致的实际表现**。不同于以往追求单一编译器的突破，这一年的重点在于通过引入全新的中间层编译器、重构顶层架构，以及跨组件（与 Blink DOM 引擎）的协同，消除性能木桶的短板，将“速度”从单一指标转化为全链路的丝滑体验。

### JavaScript 优化亮点 | JavaScript Optimization Highlights

#### 1. Maglev 集成：填补性能鸿沟

- **定位**：Maglev 是一款全新的 SSA（静态单赋值）基础的中间层优化编译器，成功填补了 **Sparkplug**（极速生成的基准代码）与 **TurboFan**（耗时较长、追求最高性能的优化代码）之间的空白。
- **性能提升**：它生成代码的速度比 TurboFan 快 **10-100 倍**，而代码运行效率显著优于 Sparkplug（使 JetStream 提升了 **8.2%**，Speedometer 提升了 **6%**）。
- **洞察**：Maglev 的引入使得 V8 可以在无需等待 TurboFan 漫长编译的情况下，更早地切换到优化代码运行，极大地缓解了在密集任务下的 CPU 争抢。

#### 2. 脚本解析与流式解析优化

- **HTML 解析加速**：V8 团队将其性能特长应用到了 Blink 的 HTML 解析器中，使其在 Speedometer 分数上获得了 **3.4%** 的显著提升。
- **启动速度提升**：通过改进脚本流式解析（Streaming Parsing）和并发编译策略，降低了主线程的占用。在关键路径上，V8 现在能更早地在后台完成解析和基准代码生成，使部分网页的**启动响应速度（Startup Latency）提升了约 20%-40%**。

### WebAssembly 优化亮点 | WebAssembly Optimization Highlights

#### 1. WasmGC 和 JSPI 的落地

- **WasmGC**：正式支持 WebAssembly 垃圾回收机制。这意味着 Java、Kotlin、Dart 等带有内存管理特性的语言可以更直接、高效地运行在 V8 上。实测显示，通过 WasmGC 移植的应用运行速度可比直接编译为 JS 快 **2 倍**，且二进制体积更小。
- **JSPI (JS Promise Integration)**：解决了 Wasm 传统的同步执行模型与 JS 异步 API（如 Fetch）之间的矛盾，使得开发者能在不重构 C++ 代码的情况下，以同步风格调用异步操作。

#### 2. Turboshaft：编译器流水线的革命

- **Turboshaft 转型**：TurboFan 的后端架构正逐步重构为全新的 **Turboshaft**。
- **效能飞跃**：在 Chrome 120 中，通过 Turboshaft 处理的编译阶段比旧版 TurboFan 快了 **2 倍**。这种架构提供了更灵活的优化通路，为未来更复杂的机器码优化奠定了基础。

### 行业标准：针对 Speedometer 3 的毫秒级打磨

V8 并非盲目追求跑分，而是通过分析 Speedometer 3 中的实际工作负载进行“微诊查”：

- **TDZ Check Elision**：通过静态分析消除了冗余的 `let`/`const` 暂时性死区检查。
- **DOM 分配优化**：重构了 Oilpan（DOM 对象分配器），使其分配性能在 DOM 密集型任务中提升了 **3 倍**。
- **针对性 API 优化**：对 `Array.prototype.groupBy` 等新标准 API 进行了深度硬件加速。

---

### 一针见血的洞察：从“单兵作战”到“体系化性能协同”

2023 年标志着 V8 架构迈向了**“精细化多级分层”**的新纪元。

不再存在所谓“最强的编译器”，而是形成了一个**协同作战的编译器梯队**：

- **Ignition** 负责极速启动；
- **Sparkplug** 负责瞬间基准性能；
- **Maglev** 在几毫秒内接管并提供极高性价比的优化；
- **TurboFan/Turboshaft** 则在后台默默打磨长期运行的顶尖性能代码。

这种**体系化性能协同**意味着 V8 把性能战场从“单纯的代码转换”扩展到了“跨主线程、多编译器、GC 调度与内存分配”的全局博弈。通过这种毫秒级的打磨，V8 巩固了其作为现代 Web 基础设施“最快心脏”的统治地位。
