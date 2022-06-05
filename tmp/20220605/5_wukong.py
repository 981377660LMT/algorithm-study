class Block:
    def __init__(self, text):
        self.text = text
        self.prev = None
        self.next = None

    def add(self, b):
        b.next = self.next
        b.prev = self
        self.next.prev = b
        self.next = b

    def deled(self):
        self.prev.next = self.next
        self.next.prev = self.prev


class TextEditor:
    def __init__(self):
        self.head = Block("")
        self.tail = Block("")

        self.head.next = self.tail
        self.tail.prev = self.head

        self.cur = self.head
        self.index = 0  # 光标，下一个字符输入的位置，以及左侧字符的长度

    def addText(self, text: str) -> None:
        if self.index == len(self.cur.text):
            self.cur.add(Block(text))
            self.cur = self.cur.next
            self.index = len(text)
        elif self.index == 0:
            if self.cur != self.head:
                self.cur = self.cur.prev
            self.cur.add(Block(text))
            self.cur = self.cur.next
            self.index = len(text)
        else:
            self.cur.add(Block(self.cur.text[self.index :]))
            self.cur.text = self.cur.text[: self.index]
            self.cur.add(Block(text))
            self.cur = self.cur.next
            self.index = len(text)

    def deleteText(self, k: int) -> int:
        if self.cur is self.head:
            return 0
        ans = 0

        while self.cur != self.head and k:
            # 如果当前block可以删完剩下的
            if self.index >= k:
                self.cur.text = self.cur.text[: self.index - k] + self.cur.text[self.index :]
                self.index -= k
                ans += k
                k = 0
            else:
                # 不能直接删完，就先把当前光标左侧block的删了
                k -= self.index
                ans += self.index
                nxt = self.cur.prev

                if self.index == len(self.cur.text):
                    self.cur.deled()
                else:
                    self.cur.text = self.cur.text[self.index :]

                self.cur = nxt
                self.index = len(self.cur.text)
        return ans

    def cursorLeft(self, k: int) -> str:
        while self.cur != self.head and k:
            if self.index >= k:
                self.index -= k
                k = 0
            else:
                k -= self.index
                self.cur = self.cur.prev
                self.index = len(self.cur.text)

        record = (self.cur, self.index)

        # 取字符
        ans = []
        r = 10
        while self.cur != self.head and r:
            if self.index >= r:
                ans.append(self.cur.text[self.index - r : self.index])
                self.index -= r
                r = 0
            else:
                r -= self.index
                ans.append(self.cur.text[: self.index])
                self.cur = self.cur.prev
                self.index = len(self.cur.text)

        self.cur, self.index = record
        return ''.join(ans[::-1])

    def cursorRight(self, k: int) -> str:
        while self.cur != self.tail and k:
            if len(self.cur.text) - self.index >= k:
                self.index += k
                k = 0
            else:
                k -= len(self.cur.text) - self.index
                self.cur = self.cur.next
                self.index = 0

        if self.cur == self.tail:
            self.cur = self.cur.prev
            self.index = len(self.cur.text)

        record = (self.cur, self.index)

        # 取字符
        ans = []
        r = 10
        while self.cur != self.head and r:
            if self.index >= r:
                ans.append(self.cur.text[self.index - r : self.index])
                self.index -= r
                r = 0
            else:
                r -= self.index
                ans.append(self.cur.text[: self.index])
                self.cur = self.cur.prev
                self.index = len(self.cur.text)

        self.cur, self.index = record
        return ''.join(ans[::-1])

    # Your TextEditor object will be instantiated and called as such:
