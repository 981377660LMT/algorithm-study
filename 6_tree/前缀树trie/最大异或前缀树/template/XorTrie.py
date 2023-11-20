"""
XorTrie 最大异或前缀树
https://kazuma8128.hatenablog.com/entry/2018/05/06/022654
"""

from typing import List


class BinaryTrie:
    __slots__ = (
        "_maxLog",
        "_xEnd",
        "_vList",
        "_multiset",
        "_edges",
        "_size",
        "_endCount",
        "_maxV",
        "_lazy",
    )

    def __init__(self, max=1 << 30, addLimit=int(2e5 + 10), allowMultipleElements=True):
        """
        Example:

        ```
        n = len(nums)
        bt = BinaryTrie(max = 1 << 30, addLimit=n, allowMultipleElements=True)
        for num in nums:
            bt.add(num)
        res = 0
        for num in nums:
            bt.xor_all(num)
            res += bt.bisect_right(high) - bt.bisect_left(low)
            bt.xor_all(num)
        ```
        """
        maxLog = max.bit_length()
        self._maxLog = maxLog
        self._xEnd = 1 << maxLog
        self._vList = [0] * (maxLog + 1)
        self._multiset = allowMultipleElements
        n = maxLog * addLimit + 1
        self._edges = [-1] * (2 * n)
        self._size = [0] * n
        self._endCount = [0] * n
        self._maxV = 0
        self._lazy = 0

    def add(self, x: int) -> None:
        x ^= self._lazy
        v = 0
        for i in range(self._maxLog - 1, -1, -1):
            d = (x >> i) & 1
            if self._edges[2 * v + d] == -1:
                self._maxV += 1
                self._edges[2 * v + d] = self._maxV
            v = self._edges[2 * v + d]
            self._vList[i] = v
        if self._multiset or self._endCount[v] == 0:
            self._endCount[v] += 1
            for v in self._vList:
                self._size[v] += 1

    def discard(self, x: int) -> None:
        if not 0 <= x < self._xEnd:
            return
        x ^= self._lazy
        v = 0
        for i in range(self._maxLog - 1, -1, -1):
            d = (x >> i) & 1
            if self._edges[2 * v + d] == -1:
                return
            v = self._edges[2 * v + d]
            self._vList[i] = v
        if self._endCount[v] > 0:
            self._endCount[v] -= 1
            for v in self._vList:
                self._size[v] -= 1

    def erase(self, x: int, count=1):
        """删除count个x.count=-1表示删除所有x."""
        if not 0 <= x < self._xEnd:
            return
        x ^= self._lazy
        v = 0
        for i in range(self._maxLog - 1, -1, -1):
            d = (x >> i) & 1
            if self._edges[2 * v + d] == -1:
                return
            v = self._edges[2 * v + d]
            self._vList[i] = v
        if count == -1 or self._endCount[v] < count:
            count = self._endCount[v]
        if self._endCount[v] > 0:
            self._endCount[v] -= count
            for v in self._vList:
                self._size[v] -= count

    def count(self, x: int) -> int:
        if not 0 <= x < self._xEnd:
            return 0
        x ^= self._lazy
        v = 0
        for i in range(self._maxLog - 1, -1, -1):
            d = (x >> i) & 1
            if self._edges[2 * v + d] == -1:
                return 0
            v = self._edges[2 * v + d]
        return self._endCount[v]

    def bisectLeft(self, x: int) -> int:
        if x < 0:
            return 0
        if self._xEnd <= x:
            return len(self)
        v = 0
        ret = 0
        for i in range(self._maxLog - 1, -1, -1):
            d = (x >> i) & 1
            left = (self._lazy >> i) & 1
            lc = self._edges[2 * v]
            rc = self._edges[2 * v + 1]
            if left == 1:
                lc, rc = rc, lc
            if d:
                if lc != -1:
                    ret += self._size[lc]
                if rc == -1:
                    return ret
                v = rc
            else:
                if lc == -1:
                    return ret
                v = lc
        return ret

    def bisectRight(self, x: int) -> int:
        return self.bisectLeft(x + 1)

    def index(self, x: int) -> int:
        if x not in self:
            raise ValueError(f"{x} is not in BinaryTrie")
        return self.bisectLeft(x)

    def find(self, x: int) -> int:
        if x not in self:
            return -1
        return self.bisectLeft(x)

    def at(self, index: int) -> int:
        if index < 0:
            index += self._size[0]
        v = 0
        ret = 0
        for i in range(self._maxLog - 1, -1, -1):
            left = (self._lazy >> i) & 1
            lc = self._edges[2 * v]
            rc = self._edges[2 * v + 1]
            if left == 1:
                lc, rc = rc, lc
            if lc == -1:
                v = rc
                ret |= 1 << i
                continue
            if self._size[lc] <= index:
                index -= self._size[lc]
                v = rc
                ret |= 1 << i
            else:
                v = lc
        return ret

    def minimum(self) -> int:
        return self.at(0)

    def maximum(self) -> int:
        return self.at(-1)

    def xorAll(self, x: int) -> None:
        self._lazy ^= x

    def __iter__(self):
        q = [(0, 0)]
        for i in range(self._maxLog - 1, -1, -1):
            left = (self._lazy >> i) & 1
            nq = []
            for v, x in q:
                lc = self._edges[2 * v]
                rc = self._edges[2 * v + 1]
                if left == 1:
                    lc, rc = rc, lc
                if lc != -1:
                    nq.append((lc, 2 * x))
                if rc != -1:
                    nq.append((rc, 2 * x + 1))
            q = nq
        for v, x in q:
            for _ in range(self._endCount[v]):
                yield x

    def __str__(self) -> str:
        prefix = "BinaryTrie("
        content = list(map(str, self))
        suffix = ")"
        if content:
            content[0] = prefix + content[0]
            content[-1] = content[-1] + suffix
        else:
            content = [prefix + suffix]
        return ", ".join(content)

    def __getitem__(self, k: int) -> int:
        return self.at(k)

    def __contains__(self, x: int) -> bool:
        return not not self.count(x)

    def __len__(self) -> int:
        return self._size[0]

    def __bool__(self) -> bool:
        return not not len(self)

    def __ixor__(self, x: int) -> "BinaryTrie":
        self.xorAll(x)
        return self


if __name__ == "__main__":

    class Solution:
        def findMaximumXOR(self, nums: List[int]) -> int:
            trie = BinaryTrie(max=max(nums), addLimit=len(nums), allowMultipleElements=True)
            res = 0
            for num in nums:
                trie.add(num)
                trie.xorAll(num)
                res = max(res, trie.maximum())
                trie.xorAll(num)
            return res

    # 2935. 找出强数对的最大异或值 II
    # https://leetcode.cn/problems/maximum-strong-pair-xor-ii/
    # 给你一个下标从 0 开始的整数数组 nums 。如果一对整数 x 和 y 满足以下条件，则称其为 强数对 ：
    # |x - y| <= min(x, y)
    # 你需要从 nums 中选出两个整数，且满足：这两个整数可以形成一个强数对，并且它们的按位异或（XOR）值是在该数组所有强数对中的 最大值 。
    # 返回数组 nums 所有可能的强数对中的 最大 异或值。
    #
    # 排序后，变为 y <= 2x
    class Solution2:
        def maximumStrongPairXor(self, nums: List[int]) -> int:
            nums.sort()
            res, left, n = 0, 0, len(nums)
            trie = BinaryTrie(max=max(nums), addLimit=n, allowMultipleElements=True)
            for right in range(n):
                cur = nums[right]
                trie.add(cur)
                while left <= right and cur > 2 * nums[left]:
                    trie.discard(nums[left])
                    left += 1
                trie.xorAll(cur)
                res = max(res, trie.maximum())
                trie.xorAll(cur)
            return res

    # 1803. 统计异或值在范围内的数对有多少
    # https://leetcode.cn/problems/count-pairs-with-xor-in-a-range/description/
    class Solution3:
        def countPairs(self, nums: List[int], low: int, high: int) -> int:
            n = len(nums)
            bt = BinaryTrie(max=max(nums), addLimit=n, allowMultipleElements=True)
            for num in nums:
                bt.add(num)
            res = 0
            for num in nums:
                bt.xorAll(num)
                res += bt.bisectRight(high) - bt.bisectLeft(low)
                bt.xorAll(num)
            return res // 2
