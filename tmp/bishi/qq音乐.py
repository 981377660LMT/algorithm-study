from collections import Counter

from typing import List, Optional


class TreeNode:
    def __init__(
        self, val: int, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 给定一棵二叉树，二叉树的每个结点只有0或2个孩子。
# !你需要对每个结点赋值一个正整数，使得每个结点的左右子树权值和相等。
# !你需要返回所有结点的最小权值和对1e9+7取模的结果。
# 二叉树结点个数不超过1e5。
MOD = int(1e9 + 7)


class Solution:
    def getTreeSum(self, tree: TreeNode) -> int:
        def dfs(root: Optional["TreeNode"]) -> int:
            if not root:
                return 0
            left = dfs(root.left)
            right = dfs(root.right)
            return 2 * max(left, right) + 1

        return dfs(tree) % MOD


################################################################

INF = int(4e18)


class Solution:
    def minOperations(self, str: str) -> int:
        """
        每次选择两个相同的字符删除
        结尾任意新增一个小写字母
        最少删除多少次后
        所有字母不同
        n<=1000

        O(26*n)
        """
        counter = [0] * 26
        for char in str:
            ord_ = ord(char) - 97
            counter[ord_] += 1

        res = 0
        while max(counter) > 1:
            max_, maxI = -INF, -1
            for i, v in enumerate(counter):
                if v > max_:
                    max_ = v
                    maxI = i
            counter[maxI] -= 2

            min_, minI = INF, -1
            for i, v in enumerate(counter):
                if v < min_:
                    min_ = v
                    minI = i
            counter[minI] += 1
            res += 1

        return res

    def minOperations2(self, str: str) -> int:
        import string
        from sortedcontainers import SortedList

        counter = Counter(str)
        sl = SortedList([(counter[char], char) for char in string.ascii_lowercase])

        res = 0
        while sl[-1][0] > 1:
            count, char = sl.pop()
            sl.add((count - 2, char))
            count, char = sl.pop(0)
            sl.add((count + 1, char))
            res += 1
        return res


print(Solution().minOperations("abab"))
print(Solution().minOperations("aa"))
print(Solution().minOperations("aaa"))
print(Solution().minOperations("aaabb"))

print(Solution().minOperations2("abab"))
print(Solution().minOperations2("aa"))
print(Solution().minOperations2("aaa"))
print(Solution().minOperations2("aaabb"))
