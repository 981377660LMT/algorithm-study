# 脱离 Web：使用 Emscripten 生成独立 WASS 二进制文件

# Outside the web: standalone WebAssembly binaries using Emscripten

WebAssembly 不再仅仅是浏览器的附属品，v8.dev 探讨了其作为普适二进制格式的潜力。

### 1. 核心变革：WASI 整合

文章详细介绍了如何利用 WASI（WebAssembly System Interface）将 C++ 代码编译为不依赖 JS 环境的独立模块。WASI 提供了一套标准化、平台无关的系统调用接口（如读写文件、网络操作），使得 WASM 终于拥有了“操作系统”感知。

### 2. 移除 JS 胶水代码

以往使用 Emscripten 编译生成的 `.js` 胶水文件动辄数百 KB。通过新的独立编译选项，开发者可以生成仅包含业务逻辑和 WASI 调用的 `.wasm` 文件，极大地方便了在边缘计算（Serverless）和嵌入式设备上的分发。

### 3. `wasm-ld` 的作用

现代化的链接器 `wasm-ld` 取代了旧的转换工具。它能更好地处理死代码消除 (Tree Shaking) 和符号解析，生成的二进制文件不仅更小，而且启动速度更快，因为引擎不需要解析复杂的 JS 加载逻辑。

### 4. 一针见血的见解

WASM 正在经历其“Java 虚拟机”时刻——Write once, run anywhere。通过剥离 JS 依赖，WASM 正试图从一种 Web 技术转型为一种跨平台、高性能、安全隔离的底层通用计算引擎，这种转变对云原生计算具有深远影响。
