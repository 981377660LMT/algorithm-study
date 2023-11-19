from typing import List


class XorTrieSimpleNode:
    __slots__ = "count", "children"

    def __init__(self):
        self.count = 0
        self.children: List["XorTrieSimpleNode"] = [None, None]  # type: ignore


class XorTrieSimple:
    __slots__ = "root", "bit"

    def __init__(self, upper: int):
        self.root = XorTrieSimpleNode()
        self.bit = upper.bit_length()

    def insert(self, num: int) -> "XorTrieSimpleNode":
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            if root.children[bit] is None:  # type: ignore
                root.children[bit] = XorTrieSimpleNode()
            root = root.children[bit]
            root.count += 1
        return root

    def remove(self, num: int) -> None:
        """需要保证num在trie中."""
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            root = root.children[bit]
            root.count -= 1

    def query(self, num: int) -> int:
        """查询num与trie中的最大异或值."""
        root = self.root
        res = 0
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root.children[needBit] is not None and root.children[needBit].count > 0:  # type: ignore
                res |= 1 << i
                bit = needBit
            root = root.children[bit]

        return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/maximum-strong-pair-xor-ii/
    # 给你一个下标从 0 开始的整数数组 nums 。如果一对整数 x 和 y 满足以下条件，则称其为 强数对 ：
    # |x - y| <= min(x, y)
    # 你需要从 nums 中选出两个整数，且满足：这两个整数可以形成一个强数对，并且它们的按位异或（XOR）值是在该数组所有强数对中的 最大值 。
    # 返回数组 nums 所有可能的强数对中的 最大 异或值。
    #
    # 排序后，变为 y <= 2x
    class Solution:
        def maximumStrongPairXor(self, nums: List[int]) -> int:
            nums.sort()
            res, left, n = 0, 0, len(nums)
            trie = XorTrieSimple(max(nums))
            for right in range(n):
                cur = nums[right]
                trie.insert(cur)
                while left <= right and cur > 2 * nums[left]:
                    trie.remove(nums[left])
                    left += 1
                res = max(res, trie.query(cur))
            return res
