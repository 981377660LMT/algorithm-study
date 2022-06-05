from dataclasses import dataclass
from heapq import heapify, heappop, heappush
from typing import List, Optional

# !文本编辑器的游标表示左边那个结点


# !对顶栈 双向链表 分块链表 rope
# 好写的链表
@dataclass(slots=True)
class TextNode:

    value: str
    left: Optional['TextNode'] = None
    right: Optional['TextNode'] = None

    def __repr__(self) -> str:
        return f'{self.value}'


# 1 <= text.length, k <= 40
# text 只含有小写英文字母。
# 调用 addText ，deleteText ，cursorLeft 和 cursorRight 的 总 次数不超过 2 * 104 次。

# !数组和字符串暴力不一样
class TextEditor:
    def __init__(self):
        self.head = TextNode('#', None, None)  # type: ignore
        self.cursor = self.head

    def addText(self, text: str) -> None:
        for char in text:
            preNode, nextNode = self.cursor, self.cursor.right
            newNode = TextNode(char, None, None)
            preNode.right = newNode
            newNode.left = preNode
            newNode.right = nextNode
            if nextNode:
                nextNode.left = newNode
            self.cursor = newNode

    def deleteText(self, k: int) -> int:
        def remove(node: 'TextNode') -> None:
            if node.left:
                node.left.right = node.right
            if node.right:
                node.right.left = node.left

        res = 0
        for _ in range(k):
            if self.cursor.value == '#':
                break
            remove(self.cursor)
            if self.cursor.left:
                self.cursor = self.cursor.left
            res += 1
        return res

    def cursorLeft(self, k: int) -> str:
        for _ in range(k):
            if self.cursor.value == '#':
                break
            if self.cursor.left:
                self.cursor = self.cursor.left
        return self._getLeft()

    def cursorRight(self, k: int) -> str:
        for _ in range(k):
            if not self.cursor.right:
                break
            if self.cursor.right:
                self.cursor = self.cursor.right
        return self._getLeft()

    def _getLeft(self) -> str:
        res = []
        p = self.cursor
        for _ in range(10):
            if not p or p.value == '#':
                break
            res.append(p.value)
            p = p.left
        return ''.join(res[::-1])


textEditor = TextEditor()  # 当前 text 为 "|" 。（'|' 字符表示光标）
textEditor.addText("leetcode")  # 当前文本为 "leetcode|" 。
print(textEditor.deleteText(4))  # 返回 4
# 当前文本为 "leet|" 。
# 删除了 4 个字符。
textEditor.addText("practice")
# 当前文本为 "leetpractice|" 。
print(textEditor.cursorRight(3))
# 返回 "etpractice"
# 当前文本为 "leetpractice|".
# 光标无法移动到文本以外，所以无法移动。
# "etpractice" 是光标左边的 10 个字符。
print(textEditor.cursorLeft(8))
# 返回 "leet"
# 当前文本为 "leet|practice" 。
# "leet" 是光标左边的 min(10, 4) = 4 个字符。
print(textEditor.deleteText(10))
# 返回 4
# 当前文本为 "|practice" 。
# 只有 4 个字符被删除了。
print(textEditor.cursorLeft(2))
# 返回 ""
# 当前文本为 "|practice" 。
# 光标无法移动到文本以外，所以无法移动。
# "" 是光标左边的 min(10, 0) = 0 个字符。
print(textEditor.cursorRight(6))
# 返回 "practi"
# 当前文本为 "practi|ce" 。
# "practi" 是光标左边的 min(10, 6) = 6 个字符。


# ["TextEditor","addText","cursorLeft","deleteText","cursorLeft","addText","cursorRight"]
# [[],["bxyackuncqzcqo"],[12],[3],[5],["osdhyvqxf"],[10]]
# [null,null,"bx",2,"",null,"yackuncqzc"]


# [null, null, 4, null, "etpractice", "leet", 4, "", "practi"]
