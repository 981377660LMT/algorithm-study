# V8 release v3.3 发布：垃圾回收平滑化与实验特性

# V8 release v3.3: GC Smoothing and Experimental Features

发布于 2011 年中旬，此版本重点解决了 Web 应用中的“卡顿”问题。

### 1. 垃圾回收暂停时间的微调 (GC Smoothing)

在 Crankshaft 普及后，V8 的峰值性能已经很高，但垃圾回收（GC）产生的停顿（Stutters）成为了体验杀手。v3.3 开始引入更精细的增量标记（Incremental Marking）实验，试图将漫长的 Full GC 拆分为多个微小的片段。

### 2. 实验性的 ES6 提案支持

在这个版本中，V8 团队开始尝试实现一些当时尚处在草案阶段的 Harmony (ES6) 特性，如 `let` 和 `const` 的早期实验。这开启了 V8 作为标准先行者的角色。

### 3. 一针见血的见解

性能不仅是平均速度（Throughput），更是平滑度（Latency）。v3.3 在通过 Crankshaft 获取速度后，开始回头修补 GC 带来的延迟问题。这种“先求快，再求稳”的迭代策略，是 V8 长期领先的核心逻辑。

---

- **归档链接**: [V8 Blog (Archived)](https://v8.dev/blog)
- **核心关键词**: `GC Latency`, `Harmony`, `ES6 Previews`, `Incremental Marking`
