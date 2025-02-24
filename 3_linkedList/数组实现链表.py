# https://kazun1998.github.io/library_for_python/Doubly_Linked_List.py


class DoublyLinkedList:
    """数组实现的双向链表,维护 0 ~ n-1 的节点."""

    __slots__ = ("_n", "_front", "_back")

    def __init__(self, n):
        self._n = n
        self._front = [-1] * n
        self._back = [-1] * n

    def previous(self, x, default=-1):
        return self._front[x] if self._front[x] != -1 else default

    def next(self, x, default=-1):
        return self._back[x] if self._back[x] != -1 else default

    def disconnect_front(self, x):
        """x から前に伸びるリンクを削除する."""
        front = self._front
        back = self._back
        y = front[x]
        if y >= 0:
            front[x] = -1
            back[y] = -1

    def disconnect_back(self, x):
        """x から後ろに伸びるリンクを削除する."""
        front = self._front
        back = self._back
        y = back[x]
        if y >= 0:
            back[x] = -1
            front[y] = -1

    def extract(self, x):
        """x に接続するリンクを削除し, x の前後が存在するならば, それらをつなぐ."""

        a = self._front[x]
        b = self._back[x]

        self.disconnect_front(x)
        self.disconnect_back(x)

        if a != -1 and b != -1:
            self.connect(a, b)

    def connect(self, x, y):
        """x から y へのリンクを生成する (すでにある x からのリンクと y へのリンクは削除される)."""

        self.disconnect_back(x)
        self.disconnect_front(y)
        self._back[x] = y
        self._front[y] = x

    def insert_front(self, x, y):
        """x の前に y を挿入する."""

        z = self._front[x]
        self.connect(y, x)
        if z != -1:
            self.connect(z, y)

    def insert_back(self, x, y):
        """x の後に y を挿入する."""

        z = self._back[x]
        self.connect(x, y)
        if z != -1:
            self.connect(y, z)

    def head(self, x):
        """x が属する弱連結成分の先頭を求める."""
        while self._front[x] != -1:
            x = self._front[x]
        return x

    def tail(self, x):
        """x が属する弱連結成分の末尾を求める."""
        while self._back[x] != -1:
            x = self._back[x]
        return x

    def enumerate(self, x):
        """x が属している弱連結成分を先頭から順に出力する."""
        x = self.head(x)
        res = [x]
        while self._back[x] >= 0:
            x = self._back[x]
            res.append(x)
        return res

    def depth(self, x):
        dep = 0
        while self._front[x] != -1:
            x = self._front[x]
            dep += 1
        return dep

    def __len__(self):
        return self._n

    def __str__(self):
        res = []
        used = [0] * self._n

        for x in range(self._n):
            if used[x]:
                continue

            a = self.enumerate(x)
            for y in a:
                used[y] = 1
            res.append(a)
        return str(res)

    def __repr__(self):
        return "[Doubly Linked List]: " + str(self)


if __name__ == "__main__":
    # 初始化一个包含 10 个节点的双向链表，所有的链接初始为 -1（表示未连接）
    dll = DoublyLinkedList(10)

    # 手动构造链表：连接 0 -> 1 -> 2 -> 3
    dll.connect(0, 1)
    dll.connect(1, 2)
    dll.connect(2, 3)
    print("初始链表 (0 -> 1 -> 2 -> 3):", dll.enumerate(0))

    # 在节点1前插入节点4，
    # 操作前链表: 0 -> 1 -> 2 -> 3
    # 操作后链表: 0 -> 4 -> 1 -> 2 -> 3
    dll.insert_front(1, 4)
    print("插入节点 4 到节点 1 之前:", dll.enumerate(0))

    # 在节点2后插入节点5，
    # 操作前链表: 0 -> 4 -> 1 -> 2 -> 3
    # 操作后链表: 0 -> 4 -> 1 -> 2 -> 5 -> 3
    dll.insert_back(2, 5)
    print("插入节点 5 到节点 2 之后:", dll.enumerate(0))

    # 从链表中提取（删除）节点1，提取后会将节点1的前后相连
    dll.extract(1)
    print("提取节点 1 后的链表:", dll.enumerate(0))

    # 测试 head 和 tail 方法：
    # 这里对于节点 5，我们获取其所属连通分组的头节点和尾节点
    head = dll.head(5)
    tail = dll.tail(5)
    print("包含节点 5 的连通分量的头节点:", head)
    print("包含节点 5 的连通分量的尾节点:", tail)

    # 最后，输出整个双向链表 (自动分组输出各连通分量)
    print("整个链表结构:", dll)
