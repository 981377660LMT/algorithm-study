# 72. V8 release v6.4 | V8 v6.4 发布：Unicode 属性转义与具名捕获组

V8 v6.4 是正则表达式进化的一个里程碑，正式将多个 TC39 正则提案落地。

## 1. 正则新语法落地

- **Unicode 属性转义 (`\p{...}`)**：开发者可以轻松匹配 Emoji、特定语种字符等，而不需要庞大的硬编码范围表。
- **具名捕获组**：允许通过 `(?<name>...)` 为匹配组命名，提高代码可读性。

## 2. 内存优化：Orinoco 并行清除

此版本的 Minor GC (Scavenger) 现在是**全并行**的。

- 之前标记工作是并行的，但移动对象仍有部分串行。
- 现在利用 `TaskRunner` 将整个清理阶段分布在多核上。

## 3. 性能：Instanceof 优化

V8 对 `instanceof` 操作符进行了编译器层面的内联优化。

- 如果原型链相对稳定，TurboFan 可以直接预测结果，将昂贵的原型链遍历简化为一次 Shape（Map）检查。

## 4. 其它改进

- **String.prototype.padStart / padEnd**：内建实现性能优化。
- **SharedArrayBuffer**：此版本出于安全考虑对 SAB 进行了临时的调整，预示着随后对 Meltdown/Spectre 漏洞的全面修复。
