# https://leetcode.cn/problems/print-immutable-linked-list-in-reverse/solutions/1075083/fen-er-zhi-zhi-zhen-zheng-de-onshi-jian-5ljbe/
# https://cs.stackexchange.com/questions/68769/print-the-nodes-of-an-immutable-single-linked-list-in-reverse-order
#
# 在使用递归或者栈方式时，由于需要一次性将链表全部压栈，所以需要n的栈空间。
# 为了减少栈空间，这里使用的是分区的方法，也就是将链表划分为前后两部分分别进行递归，这样栈空间就缩小了一半。


class ImmutableListNode:
    def printValue(self) -> None: ...  # print the value of this node.
    def getNext(self) -> "ImmutableListNode": ...  # return the next node.


class Solution:
    def printLinkedListInReverse(self, head: "ImmutableListNode") -> None:
        if not head:
            return

        # ---------- 1) 首遍扫描，记录每 √n 本书的块头 ----------
        import math

        n, cur = 0, head
        while cur:
            n += 1
            cur = cur.getNext()

        B = int(math.isqrt(n)) + 1  # 块大小（≈√n）
        block_heads = []  # 至多 √n 个块头
        idx, cur = 0, head
        while cur:
            if idx % B == 0:
                block_heads.append(cur)  # 记录块头
            cur = cur.getNext()
            idx += 1

        # ---------- 2) 逆序遍历块，再块内反向打印 ----------
        for bh_idx in range(len(block_heads) - 1, -1, -1):
            start = block_heads[bh_idx]
            end = block_heads[bh_idx + 1] if bh_idx + 1 < len(block_heads) else None

            tmp_stack = []
            cur = start
            while cur is not end:
                tmp_stack.append(cur)
                cur = cur.getNext()

            while tmp_stack:
                tmp_stack.pop().printValue()
