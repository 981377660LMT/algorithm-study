# Visual Studio Code 在 Build 2021

链接：https://code.visualstudio.com/blogs/2021/05/06/code-build-2021

## 深入分析

Build 2021的VS Code演讲，反映了远程开发和协作的深化。

### 主要宣布
1. **Remote Tunnels** - 无需SSH密钥，仅凭VS Code账户即可连接远程
2. **GitHub Copilot的早期版本** - AI辅助代码编写（当时还叫"GitHub Copilot for Technical Preview"）
3. **Workspace Trust** - 安全特性，防止恶意代码执行

### GitHub Copilot的意义
- 这是LLM（大语言模型）在IDE中的首次大规模应用
- 基于OpenAI的Codex模型，经过GitHub代码库微调
- 提供上下文感知的代码补全和生成

### Workspace Trust的必要性
- VS Code可以运行任意脚本（通过扩展），存在安全隐患
- Workspace Trust要求用户明确信任某个工作区，才能执行脚本
- 这对在不信任的代码库（如开源项目）中工作特别重要

### 时代转折点
- Build 2021标志着VS Code从"编辑工具"向"AI驱动开发环境"的转变
- GitHub Copilot虽然在2021年还不完美，但预示了未来的方向
