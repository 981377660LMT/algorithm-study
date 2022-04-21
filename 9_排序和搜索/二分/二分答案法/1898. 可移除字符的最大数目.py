from typing import List

# 1 <= p.length <= s.length <= 105
# removable 中的元素 互不相同

# 返回你可以找出的 最大 k ，满足在移除字符后 p 仍然是 s 的一个子序列。


class Solution:
    def maximumRemovals(self, s: str, p: str, removable: List[int]) -> int:
        # 辅助函数，用来判断
        def check(k: int) -> bool:
            """移除 k 个下标后 p 是否是 s_k 的子序列"""
            ok = [True] * ns  # 标记删除，不是真的删除
            for i in range(k):
                ok[removable[i]] = False

            hit = 0
            for i in range(ns):
                if ok[i] and s[i] == p[hit]:
                    hit += 1
                    if hit == np:
                        return True
            return False

        ns, np = len(s), len(p)

        left, right = 0, len(removable)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maximumRemovals(s="abcacb", p="ab", removable=[3, 1, 0]))
# 输出：2
# 解释：在移除下标 3 和 1 对应的字符后，"abcacb" 变成 "accb" 。
# "ab" 是 "accb" 的一个子序列。
# 如果移除下标 3、1 和 0 对应的字符后，"abcacb" 变成 "ccb" ，那么 "ab" 就不再是 s 的一个子序列。
# 因此，最大的 k 是 2 。

