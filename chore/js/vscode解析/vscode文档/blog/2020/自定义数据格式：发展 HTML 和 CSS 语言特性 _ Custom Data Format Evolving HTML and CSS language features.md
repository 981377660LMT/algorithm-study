# 自定义数据格式：发展 HTML 和 CSS 语言特性

链接：https://code.visualstudio.com/blogs/2020/08/04/custom-data

## 深入分析

这是一个底层的架构功能，允许VS Code支持自定义的HTML/CSS扩展。

### 背景
- Web标准（HTML、CSS）在不断演进，新增新元素和属性
- VS Code需要及时支持这些新特性的自动补全、验证等
- 但硬编码每个标准可能跟不上变化

### Custom Data Format
- 一个JSON格式的定义文件，描述HTML元素、CSS属性等
- VS Code可以读取这个文件，动态生成补全和验证规则
- 扩展开发者可以提供自己的Custom Data，支持自定义HTML标签

### 使用场景
1. **Web Components** - 自定义标签如`<my-element>`
2. **Framework特定标签** - Vue的`<component>`、React的JSX等
3. **语言扩展** - 定义新的CSS属性支持

### 竞争优势
- WebStorm虽然对HTML/CSS的支持很全面，但不够灵活
- VS Code的Custom Data机制，允许社区扩展，形成了一个良好的生态
