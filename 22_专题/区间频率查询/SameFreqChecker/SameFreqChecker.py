# O(1) 判断数据结构内所有元素的频率是否相同.
#
# 1. SameFreqCheckerAddOnly:
#    只能添加元素,不能删除元素
# 2. SameFreqChecker:
#    可以添加元素,删除元素
#
# !判断方法1：最大频率*元素种类数 == 总元素个数
# !判断方法2：另一个哈希表记录频率的频率，看这个哈希表的大小是否为1


from typing import Any


def max2(a: int, b: int) -> int:
    return a if a > b else b


class SameFreqCheckerAddOnly:
    __slots__ = "_counter", "_maxFreq", "_count"

    def __init__(self):
        self._counter = dict()
        self._maxFreq = 0
        self._count = 0

    def add(self, num: Any) -> None:
        pre = self._counter.get(num, 0)
        self._counter[num] = pre + 1
        self._maxFreq = max2(self._maxFreq, pre + 1)
        self._count += 1

    def check(self) -> bool:
        return self._maxFreq * len(self._counter) == self._count


class SameFreqChecker:
    __slots__ = "_counter", "_freqCounter"

    def __init__(self):
        self._counter = dict()
        self._freqCounter = dict()

    def add(self, num: Any) -> None:
        preC = self._counter.get(num, 0)
        self._counter[num] = preC + 1
        self._freqCounter[preC + 1] = self._freqCounter.get(preC + 1, 0) + 1
        if preC > 0:
            preF = self._freqCounter.get(preC, 0)
            if preF == 1:
                self._freqCounter.pop(preC)
            else:
                self._freqCounter[preC] = preF - 1

    def discard(self, num: Any) -> bool:
        preC = self._counter.get(num, 0)
        if preC == 0:
            return False
        if preC == 1:
            self._counter.pop(num)
        else:
            self._counter[num] = preC - 1
        preF = self._freqCounter.get(preC, 0)
        if preF == 1:
            self._freqCounter.pop(preC)
        else:
            self._freqCounter[preC] = preF - 1
        if preC > 1:
            self._freqCounter[preC - 1] = self._freqCounter.get(preC - 1, 0) + 1
        return True

    def check(self) -> bool:
        return len(self._freqCounter) == 1


if __name__ == "__main__":
    # 100289. 分割字符频率相等的最少子字符串
    # https://leetcode.cn/problems/minimum-substring-partition-of-equal-character-frequency/
    # 注意：一个 平衡 字符串指的是字符串中所有字符出现的次数都相同。

    class Solution:
        def minimumSubstringsInPartition(self, s: str) -> int:
            from functools import lru_cache

            INF = int(1e18)

            def min2(a: int, b: int) -> int:
                return a if a < b else b

            @lru_cache(None)
            def dfs(index: int) -> int:
                if index >= n:
                    return 0
                res = INF
                C = SameFreqCheckerAddOnly()
                for i in range(index, n):
                    C.add(s[i])
                    if C.check():
                        res = min2(res, 1 + dfs(i + 1))
                return res

            n = len(s)
            res = dfs(0)
            dfs.cache_clear()
            return res

    class BruteForce:
        def __init__(self) -> None:
            self.counter = dict()

        def add(self, num: Any) -> None:
            self.counter[num] = self.counter.get(num, 0) + 1

        def remove(self, num: Any) -> None:
            if num not in self.counter:
                return
            self.counter[num] -= 1
            if self.counter[num] == 0:
                self.counter.pop(num)

        def check(self) -> bool:
            return len(set(self.counter.values())) == 1

    def check() -> None:
        from random import randint

        C1 = SameFreqChecker()
        C2 = BruteForce()
        for _ in range(100000):
            v = randint(1, 26)
            C1.add(v)
            C2.add(v)
            assert C1.check() == C2.check()
            if randint(0, 1) == 0:
                C1.discard(v)
                C2.remove(v)
                assert C1.check() == C2.check()

        print("Pass")

    check()
