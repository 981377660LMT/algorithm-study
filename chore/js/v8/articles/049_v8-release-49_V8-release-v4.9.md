# V8 v4.9 发布：ES6 覆盖率达 91% 与随机数革命

# V8 release v4.9: 91% ES6 coverage and Math.random improvements

V8 v4.9 让现代 JavaScript 的核心特性在浏览器和 Node.js 中几乎触手可及。

### 1. 现代语法的集体爆发

交付了解构赋值 (Destructuring)、默认参数 (Default Parameters) 以及极其复杂的 Proxies 和 Reflect API。这些特性的落地不仅需要语法支撑，更需要 V8 在隐藏类（Hidden Class）和 IC 系统中进行深层的架构适配。

### 2. 数学引擎改进：xorshift128+

重构了 `Math.random()`。老算法在统计学检测下存在显著的可预测性漏洞，v4.9 切换到了 **xorshift128+** 算法，这确保了 JS 环境能产生更高质量、更符合统计学期望的伪随机数序列。

### 3. 一针见血的见解

v4.9 标志着 V8 对待随机性（Randomness）态度的严谨化。同时，91% 的 ES6 覆盖率意味着“编译到 ES5”的必要性开始动摇，JS 开发者的工具链开始迎来从“降级”到“原生运行”的重大范式迁移。
