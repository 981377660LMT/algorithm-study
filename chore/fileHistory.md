在 VS Code 中查看文件改动历史有以下几种方式：

**1. Timeline 面板（最简单）**

- 在资源管理器侧边栏底部有 **Timeline** 面板
- 点击任意文件后，Timeline 会显示该文件的 Git 提交历史

**2. 右键菜单**

- 在编辑器标签页或资源管理器中右键文件
- 选择 **Open Timeline** 或 **Git: View File History**

**3. Source Control 面板**

- 点击左侧 Source Control 图标（或 `Ctrl/Cmd+Shift+G`）
- 可以查看当前工作区的改动对比

**4. Git Lens 扩展（推荐）**

- 安装 **GitLens** 扩展后功能更强大
- 可以看到每行代码的最后修改者（inline blame）
- 文件历史、分支对比、提交详情等

**5. 命令面板**

- `Ctrl/Cmd+Shift+P` → 输入 `Git: View File History`

**6. 终端命令**

```bash
git log --follow -p -- <文件路径>   # 查看带 diff 的历史
git log --oneline -- <文件路径>     # 简洁历史
```
