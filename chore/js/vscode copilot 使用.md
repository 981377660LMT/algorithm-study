# vscode copilot 使用

https://code.visualstudio.com/updates/v1_97
https://code.visualstudio.com/docs/copilot/overview

## copilot 使用

1. `#` 开头的命令
   - #changes
     引用在 Git 源代码控制中修改的文件
   - #codebase
     引用代码库中的相关文件区块、符号和其他信息。
   - #sym
     引用代码库中的符号
   - #terminalLastCommand
     引用最后一次在终端中运行的命令
   - #vscodeAPI
     使用 VS Code API 引用回答有关 VS Code 扩展开发的问题。
2. `@` 开头的命令
   - @workspace
     引用有关工作区的信息
   - @github
     获取基于 Web 搜索、代码搜索和企业知识库的答案

## copilot 周边插件

## vscode-mermAId

还有一种方法是直接 claude 3.7 生成 PlantUML/mermaid 图。

## Data Analysis for Copilot

---

1. git blame
   ![alt text](image-3.png)
2. Output panel filtering 输出面板过滤
   ![alt text](image-4.png)
3. Agent mode (Experimental)
   代理模式（实验性）
   We've been working on a new agent mode for Copilot Edits. When in agent mode, Copilot can automatically search your workspace for relevant context, edit files, check them for errors, and run terminal commands (with your permission) to complete a task end-to-end.
   我们一直在为 Copilot Edits 开发一种新的代理模式。在代理模式下，Copilot 可以自动搜索您的工作区以获取相关上下文，编辑文件，检查错误，并在您的许可下运行终端命令，以完成端到端的任务。
