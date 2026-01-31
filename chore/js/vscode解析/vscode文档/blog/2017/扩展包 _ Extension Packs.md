# 扩展包

链接：https://code.visualstudio.com/blogs/2017/03/07/extension-pack-roundup

## 一针见血的分析

1. **解决依赖地狱**：随着 VS Code 插件的细粒度化，配置一个完整的开发环境（如 C++ 或 Azure 开发）可能需要安装 5-10 个插件。「扩展包」的推出是官方层面对生态系统的梳理，将这种配置成本从用户转移到了包维护者身上。
2. **场景化导流**：扩展包不仅仅是插件的集合，它们是 VS Code 渗透到不同领域的「登陆舰」。例如通过一个「Azure Tools」包，微软可以打包销售其云服务的开发套件。
3. **极简的开发者体验**：基于 `package.json` 中的 `extensionPack` 字段，实现「一键安装/一键禁用」。这种极简设计鼓励了社区通过「策展」（Curating）来建立影响力。

## 摘要

VS Code 介绍了**扩展包 (Extension Packs)** 这一特性，旨在帮助用户一次性安装多个相关的插件，简化复杂开发环境的搭建过程。

- **工作原理**：通过在扩展的 `package.json` 中定义 `extensionPack` 列表来实现依赖捆绑。
- **制作成本**：使用 Yeoman 的 `yo code` 生成器可以非常快速地创建新扩展包，甚至可以基于当前已打开的扩展自动生成。
- **推荐列表**：
  - **Azure Tools Extension Pack**：集合了管理各种 Azure 资源的工具。
  - **Node.js Extension Pack**：汇集了 Node 开发所需的常用调试与性能工具。
  - **React Native iOS Pack**：为跨平台开发提供一站式支持。
- **使用建议**：用户可以为自己的特定项目或技术栈创建扩展包，方便与团队共享或跨机器维护配置。
