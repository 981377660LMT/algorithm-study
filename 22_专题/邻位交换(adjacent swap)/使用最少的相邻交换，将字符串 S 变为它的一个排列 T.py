# 贪心地对于每一个位置，选取最靠前且还未被使用的相同字母。
# 使用这一字母需要移动的次数等于它的初始位置前还未被使用的字母的总个数。

# 2193. 得到回文串的最少操作次数
# 通过双端队列维护每个字母出现的位置，
# 并用树状数组维护哪些位置的字母已经被交换出去。
# 复杂度 O(nlogn)。
from collections import Counter, defaultdict


class BIT:
    """单点修改"""

 def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

class Solution:
    def minMovesToMakePalindrome(self, s: str):
        n = len(s)
        nums = [ord(char) - 97 for char in s]
        res, bit = 0, BIT(n)
        if sum(x & 1 for x in Counter(nums).values()) > 1:
            return -1

        indexesMap = defaultdict(list)
        for i, num in enumerate(nums):
            indexesMap[num].append(i)

        costs, paris = [0] * n, []
        for num, indexes in indexesMap.items():
            count = len(indexes)
            if count & 1:
                costs[indexes[count // 2]] = n // 2 + 1
            for i in range((count) // 2):
                paris.append((indexes[i], indexes[~i]))

        # 以贪心地对于每一个位置，选取最靠前且还未被使用的相同字母。
        # 使用这一字母需要移动的次数等于它的初始位置前还未被使用的字母的总个数。
        for i, (left, right) in enumerate(sorted(paris)):
            costs[left], costs[right] = i + 1, n - i

        for i, cost in enumerate(costs):
            # 减去已经被交换的个数
            res += i - bit.query(cost)
            bit.add(cost, 1)

        return res
