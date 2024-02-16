# CF1620E Replace the Numbers-在线查询
# https://www.luogu.com.cn/problem/CF1620E
# 给出 q 个操作，操作分为两种：

# 1 x 在序列末尾插入数字 x。
# 2 x y 把序列中的所有 x 替换为 y。

# 求这个序列操作后的结果。


import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")


class Node:
    __slots__ = ("value", "next")

    def __init__(self, value: int):
        self.value = value
        self.next = None


if __name__ == "__main__":
    q = int(input())
    posFirst = dict()  # !记录每个数字的下标的链表的头节点
    posLast = dict()  # !记录每个数字的下标的链表的尾节点
    count = 0

    def add(index: int, value: int) -> None:
        newNode = Node(index)
        if value not in posFirst:
            posFirst[value] = newNode
            posLast[value] = newNode
        else:
            last = posLast[value]
            last.next = newNode
            posLast[value] = newNode

    def merge(from_: int, to: int) -> None:
        if from_ == to:
            return
        if from_ not in posFirst:
            return
        if to not in posFirst:
            posFirst[to] = posFirst[from_]
            posLast[to] = posLast[from_]
            posFirst.pop(from_)
            posLast.pop(from_)
            return
        posLast[to].next = posFirst[from_]
        posLast[to] = posLast[from_]
        posFirst.pop(from_)
        posLast.pop(from_)

    for _ in range(q):
        t, *args = map(int, input().split())
        if t == 1:
            x = args[0]
            add(count, x)
            count += 1
        else:
            x, y = args
            merge(x, y)

    res = [0] * count
    for k, v in posFirst.items():
        head = v
        while head:
            res[head.value] = k
            head = head.next
    print(*res)
