from typing import List

#  !python的字符串操作非常快(注意是字符串切片而不是数组切片)
class TextEditor:
    def __init__(self):
        self.pos = 0
        self.text = ''

    def addText(self, text: str) -> None:
        self.text = self.text[: self.pos] + text + self.text[self.pos :]
        self.pos += len(text)

    def deleteText(self, k: int) -> int:
        res = min(k, self.pos)
        self.text = self.text[: self.pos - res] + self.text[self.pos :]
        self.pos -= res
        return res

    def cursorLeft(self, k: int) -> str:
        self.pos = max(0, self.pos - k)
        return self.text[max(0, self.pos - 10) : self.pos]

    def cursorRight(self, k: int) -> str:
        self.pos = min(len(self.text), self.pos + k)
        return self.text[max(0, self.pos - 10) : self.pos]

