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

        self.cnt = self.head
        self.c = 0  # 光标，下一个字符输入的位置，以及左侧字符的长度

    def addText(self, text: str) -> None:
        if self.c == len(self.cnt.text):
            self.cnt.add(Block(text))
            self.cnt = self.cnt.next
            self.c = len(text)
        elif self.c == 0:
            if self.cnt != self.head:
                self.cnt = self.cnt.prev
            self.cnt.add(Block(text))
            self.cnt = self.cnt.next
            self.c = len(text)
        else:
            self.cnt.add(Block(self.cnt.text[self.c :]))
            self.cnt.text = self.cnt.text[: self.c]
            self.cnt.add(Block(text))
            self.cnt = self.cnt.next
            self.c = len(text)

    def deleteText(self, k: int) -> int:
        if self.cnt is self.head:
            return 0
        ans = 0

        while self.cnt != self.head and k:
            # 如果当前block可以删完剩下的
            if self.c >= k:
                self.cnt.text = self.cnt.text[: self.c - k] + self.cnt.text[self.c :]
                self.c -= k
                ans += k
                k = 0
            else:
                # 不能直接删完，就先把当前光标左侧block的删了
                k -= self.c
                ans += self.c
                nxt = self.cnt.prev

                if self.c == len(self.cnt.text):
                    self.cnt.deled()
                else:
                    self.cnt.text = self.cnt.text[self.c :]

                self.cnt = nxt
                self.c = len(self.cnt.text)
        return ans

    def cursorLeft(self, k: int) -> str:
        while self.cnt != self.head and k:
            if self.c >= k:
                self.c -= k
                k = 0
            else:
                k -= self.c
                self.cnt = self.cnt.prev
                self.c = len(self.cnt.text)

        record = (self.cnt, self.c)

        # 取字符
        ans = []
        r = 10
        while self.cnt != self.head and r:
            if self.c >= r:
                ans.append(self.cnt.text[self.c - r : self.c])
                self.c -= r
                r = 0
            else:
                r -= self.c
                ans.append(self.cnt.text[: self.c])
                self.cnt = self.cnt.prev
                self.c = len(self.cnt.text)

        self.cnt, self.c = record
        return ''.join(ans[::-1])

    def cursorRight(self, k: int) -> str:
        while self.cnt != self.tail and k:
            if len(self.cnt.text) - self.c >= k:
                self.c += k
                k = 0
            else:
                k -= len(self.cnt.text) - self.c
                self.cnt = self.cnt.next
                self.c = 0

        if self.cnt == self.tail:
            self.cnt = self.cnt.prev
            self.c = len(self.cnt.text)

        record = (self.cnt, self.c)

        # 取字符
        ans = []
        r = 10
        while self.cnt != self.head and r:
            if self.c >= r:
                ans.append(self.cnt.text[self.c - r : self.c])
                self.c -= r
                r = 0
            else:
                r -= self.c
                ans.append(self.cnt.text[: self.c])
                self.cnt = self.cnt.prev
                self.c = len(self.cnt.text)

        self.cnt, self.c = record
        return ''.join(ans[::-1])

    # Your TextEditor object will be instantiated and called as such:
