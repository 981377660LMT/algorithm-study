from typing import List, Optional

MOD = int(1e9 + 7)
INF = int(1e20)

# Definition for singly-linked list.


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# right/down/left/up 右/下/左/上
DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


class Solution:
    def spiralMatrix(self, m: int, n: int, head: Optional[ListNode]) -> List[List[int]]:
        ROW, COL = m, n
        res = [[-1] * COL for _ in range(ROW)]

        r, c = 0, 0
        di = 0
        while head:
            res[r][c] = head.val
            head = head.next
            nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            if nr < 0 or nr >= ROW or nc < 0 or nc >= COL or res[nr][nc] != -1:
                di = (di + 1) % 4
                nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            r, c = nr, nc

        return res
