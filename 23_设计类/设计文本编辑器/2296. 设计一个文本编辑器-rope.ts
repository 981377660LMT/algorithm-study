/* eslint-disable no-console */
// https://leetcode.cn/problems/design-a-text-editor/solution/vim-by-981377660lmt-3y3i/
import { Rope } from './rope'

// rope非常适合在很长的字符串上插入和删除的场景
class TextEditor {
  pos = 0
  private readonly rope = new Rope('')

  // 将 text 添加到光标所在位置。添加完后光标在 text 的右边。
  addText(text: string): void {
    this.rope.insert(this.pos, text)
    this.pos += text.length
  }

  // 删除光标左边 k 个字符。返回实际删除的字符数目。
  deleteText(k: number): number {
    const res = Math.min(k, this.pos)
    this.rope.remove(this.pos - res, this.pos)
    this.pos -= res
    return res
  }

  // 将光标向左移动 k 次。返回移动后光标左边 min(10, len) 个字符，其中 len 是光标左边的字符数目。
  cursorLeft(k: number): string {
    this.pos = Math.max(0, this.pos - k)
    return this.rope.slice(Math.max(0, this.pos - 10), this.pos)
  }

  //  将光标向右移动 k 次。返回移动后光标左边 min(10, len) 个字符
  cursorRight(k: number): string {
    this.pos = Math.min(this.pos + k, this.rope.length)
    return this.rope.slice(Math.max(0, this.pos - 10), this.pos)
  }
}

if (require.main === module) {
  const textEditor = new TextEditor()
  for (let _ = 0; _ < 100000; _++) {
    textEditor.addText('a'.repeat(4))
  }
}

export {}
