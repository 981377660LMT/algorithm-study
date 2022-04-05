# 还是dp[index][diff]
# 0 ≤ n ≤ 100,000

# 树中的最长等差数列
from collections import defaultdict, deque
from typing import Optional


# dp的本质是在DAG上用拓扑序求最短(长)路，而树本身就是一个DAG，所以可以直接用bfs层序遍历来dp求最短(长)路


class Tree:
    def __init__(self, val: int, left: Optional['Tree'] = None, right: Optional['Tree'] = None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def solve(self, root: Tree) -> int:
        if not root:
            return 0

        dp = defaultdict(lambda: 1)
        queue = deque([(None, root)])  # parent,cur

        while queue:
            parent, cur = queue.popleft()
            if cur.left:
                queue.append((cur, cur.left))
            if cur.right:
                queue.append((cur, cur.right))
            if parent is not None:
                diff = cur.val - parent.val
                dp[cur, diff] = max(dp[cur, diff], dp[parent, diff] + 1)

        return max(dp.values(), default=1)

