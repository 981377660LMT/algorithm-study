### Kill Ring 精简讲解

**Kill Ring** 是 Emacs 及其衍生编辑器（如部分 TUI 框架）中的核心机制，用于管理剪切和粘贴。它与普通操作系统的“剪贴板”主要有两点不同：

1.  **历史记录（Ring Buffer）**：它不是只存一个值，而是一个环形缓冲区。你可以通过 `yank` (粘贴最新) 和 `yank-pop` (循环切换历史) 来找回之前的被删除内容。
2.  **自动合并（Accumulation）**：连续的“删除”操作会将内容合并到同一个 Kill Ring 条目中，而不是产生一堆细碎的记录。

---

### TypeScript 实现

```typescript
/**
 * 用于类似 Emacs 的 kill/yank 操作的环形缓冲区。
 */
export class KillRing {
  private ring: string[] = []
  private readonly maxSize: number

  constructor(maxSize = 60) {
    this.maxSize = maxSize
  }

  /**
   * 向 kill ring 添加文本。
   * @param text 被删除的文本。
   * @param opts.accumulate 是否与最近一条合并（用于处理连续删除）。
   * @param opts.prepend 合并时是前缀（反向删除）还是后缀（正向删除）。
   */
  push(text: string, opts: { accumulate?: boolean; prepend?: boolean } = {}): void {
    if (!text) return

    if (opts.accumulate && this.ring.length > 0) {
      const last = this.ring.pop()!
      this.ring.push(opts.prepend ? text + last : last + text)
    } else {
      this.ring.push(text)
      if (this.ring.length > this.maxSize) {
        this.ring.shift() // 移除最旧的记录
      }
    }
  }

  /** 获取最近的一次内容 */
  peek(): string | undefined {
    return this.ring[this.ring.length - 1]
  }

  /**
   * 旋转环（用于 yank-pop）。
   * 将最近的内容移到末尾，暴露出次新的内容。
   */
  rotate(): void {
    if (this.ring.length < 2) return
    const last = this.ring.pop()!
    this.ring.unshift(last)
  }

  get length(): number {
    return this.ring.length
  }
}
```

### 核心方法说明：

- `push`: 处理新文本。如果设置了 `accumulate`，它会修改最近的一条记录，这在处理用户连续点击 `Backspace` 或 `Delete` 时非常有用。
- `rotate`: 实现“循环粘贴”的关键。当用户粘贴了一次发现不是想要的，触发 `yank-pop` 时，底层通过旋转数组来切换当前指向的历史记录。
