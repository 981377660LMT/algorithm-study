// TextEditor() 用空文本初始化对象。
// void addText(string text) 将 text 添加到光标所在位置。添加完后光标在 text 的右边。
// int deleteText(int k) 删除光标左边 k 个字符。返回实际删除的字符数目。
// string cursorLeft(int k) 将光标向左移动 k 次。返回移动后光标左边 min(10, len) 个字符，其中 len 是光标左边的字符数目。
// string cursorRight(int k) 将光标向右移动 k 次。返回移动后光标左边 min(10, len) 个字符，其中 len 是光标左边的字符数目。

import { SqrtArray } from '../SqrtArray'

// https://leetcode.cn/problems/design-a-text-editor/
class TextEditor {
  readonly _sqrt: SqrtArray<string> = new SqrtArray(0, () => '', 1 + (Math.sqrt(1e5) | 0))
  _pos = 0

  addText(text: string): void {
    text.split('').forEach(char => this._sqrt.insert(this._pos++, char))
  }

  deleteText(k: number): number {
    const res = Math.min(k, this._pos)
    this._sqrt.erase(this._pos - res, this._pos)
    this._pos -= res
    return res
  }

  cursorLeft(k: number): string {
    this._pos = Math.max(0, this._pos - k)
    return this._sqrt.slice(this._pos - 10, this._pos).join('')
  }

  cursorRight(k: number): string {
    this._pos = Math.min(this._sqrt.length, this._pos + k)
    return this._sqrt.slice(this._pos - 10, this._pos).join('')
  }
}

/**
 * Your TextEditor object will be instantiated and called as such:
 * var obj = new TextEditor()
 * obj.addText(text)
 * var param_2 = obj.deleteText(k)
 * var param_3 = obj.cursorLeft(k)
 * var param_4 = obj.cursorRight(k)
 */
export {}

// ["TextEditor","addText","cursorLeft","deleteText","cursorLeft","addText","cursorRight"]
// [[],["bxyackuncqzcqo"],[12],[3],[5],["osdhyvqxf"],[10]]

if (require.main === module) {
  const e = new TextEditor()
  e.addText('bxyackuncqzcqo')
  console.log(e.cursorLeft(12))
  console.log(e.deleteText(3))
  console.log(e.cursorLeft(5))
  e.addText('osdhyvqxf')
  console.log(e._sqrt.toString(), e._pos)
  console.log(e.cursorRight(10))
}
