class TextEditor:
    def __init__(self):
        self.stack1 = []
        self.stack2 = []

    def addText(self, text: str) -> None:
        """将 text 添加到光标所在位置。添加完后光标在 text 的右边。"""
        self.stack1.extend(list(text))

    def deleteText(self, k: int) -> int:
        """删除光标左边 k 个字符。返回实际删除的字符数目。"""
        remain = k
        while remain and self.stack1:
            self.stack1.pop()
            remain -= 1
        return k - remain

    def cursorLeft(self, k: int) -> str:
        """将光标向左移动 k 次。返回移动后光标左边 min(10, len) 个字符，其中 len 是光标左边的字符数目。"""
        while k and self.stack1:
            self.stack2.append(self.stack1.pop())
            k -= 1
        return ''.join(self.stack1[-10:])

    def cursorRight(self, k: int) -> str:
        """将光标向右移动 k 次。返回移动后光标左边 min(10, len) 个字符，其中 len 是光标左边的字符数目。"""
        while k and self.stack2:
            self.stack1.append(self.stack2.pop())
            k -= 1
        return ''.join(self.stack1[-10:])
