# V8 v9.4 发布：类静态块与控制流完整性

# V8 release v9.4: Class Static Blocks and Control-flow Integrity

V8 v9.4 聚焦于现代 JS 语法支持和底层安全架构的加固。

### 1. 技术解析：类静态初始化块 (Class Static Initialization Blocks)

本版本引入了 ECMAScript 标准中的 Class Static Initialization Blocks。

- **功能突破**：静态块允许在类定义作用域内部执行复杂的静态初始化逻辑。
- **深度洞察**：最关键的改进是静态块可以访问类的**私有静态字段**，这在以前是无法通过外部逻辑实现的。这提供了更强大的封装能力和初始化灵活性。

### 2. 安全增强：控制流完整性 (Control-flow Integrity, CFI)

在 x64 架构上增强了控制流完整性防御。

- **硬件协作**：通过硬件指令（如 Intel CET）防御 ROP（返回导向编程）攻击。
- **安全边界**：这进一步收紧了 JIT 代码生成的安全边界，减少了因内存破坏漏洞被利用的可能性。

### 3. 一针见血的见解

`super` 的私密性与安全性的硬件化，共同构成了 v9.4 的主旋律。V8 正在从软件层面和底层指令集两端同步发力，确保现代 JS 应用在变得更复杂的同时，依然能运行在坚不可摧的基础之上。
