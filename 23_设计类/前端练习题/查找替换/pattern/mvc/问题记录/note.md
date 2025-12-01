- 幽灵高亮 (Ghost Decorations): 在 \_researchModule 中，你正确地添加了 this.\_renderModule(module, [])。
  效果: 只要开始重搜，编辑器上的旧框立即消失。无论搜索多慢，用户都不会看到错位的框。
- 替换异常 (Replace Safety): 在 replace 和 replaceAll 中都包裹了 try-catch。
  效果: 即使文件只读或插件报错，UI 线程不会崩溃。
- 状态同步 (Sync): setActiveResult 中调用了 \_renderModule。
  效果: 点击列表时，编辑器不仅滚动，还能正确切换“选中/非选中”的颜色。
- 竞态处理 (Race Conditions): SearchTask 中的 finally 块检查 if (this.\_moduleCancelSources.get(module.id) === localSource)。
  评价: 这是一行非常高级的代码。它防止了旧任务结束时误删了新任务的 Token。非常棒。

---

- 只有用户手动 prev/next/自动修复，才能改变 currentIndex
- reSearch 时，页面不滚动
- 替换完后需要 gotoIndex
- Controller 控制任务超时，任务本身不关心超时，只关心任务执行
- updateView 需要在 setTimeout 中
