from typing import List
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= order.length <= 10^5
# 1 <= limit <= 10
# 1 <= order[i] <= 10^6
# 书架上有若干本书，从左至右的书籍编号记于整型数组 order 中。
# 为保证书籍多样性，管理员想取走一些重复编号的书籍，要求满足以下条件：

# 剩余书本中相同编号的书本数量均不大于 limit
# !取走的书籍数量尽可能少
# 由于存在多种整理方案，请返回剩余书本编号的排列为「最小排列」的方案。

# !22_专题/字典序列删除/1081. 不同字符的最小子序列.py


class Solution:
    def arrangeBookshelf(self, order: List[int], limit: int) -> List[int]:
        stack = []
        visited = defaultdict(int)
        remain = Counter(order)

        need = {key: min(limit, remain[key]) for key in remain}  # !每个元素最后需要的个数

        for num in order:
            # !能不能把这个数字入栈
            remain[num] -= 1
            if visited[num] == need[num]:
                continue

            # !能不能删除栈顶元素
            while (
                stack
                and stack[-1] > num
                and remain[stack[-1]] > need[stack[-1]] - visited[stack[-1]]
            ):
                visited[stack.pop()] -= 1
            stack.append(num)
            visited[num] += 1

        return stack


print(Solution().arrangeBookshelf(order=[5, 5, 6, 5], limit=2))
print(Solution().arrangeBookshelf(order=[5, 5, 6, 5], limit=3))
print(Solution().arrangeBookshelf(order=[3, 3, 9, 8, 9, 2, 8], limit=1))
print(Solution().arrangeBookshelf(order=[2, 1, 2, 2, 1, 3, 3, 1, 3, 3], limit=2))
print(Solution().arrangeBookshelf(order=[2, 2, 1, 2, 2], limit=3))
# [2,1,2,2]
