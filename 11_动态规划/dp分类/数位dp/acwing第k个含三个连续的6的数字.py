# 只要某数字的十进制表示中有三个连续的 6，古代人也认为这是个魔鬼的数，比如 666,1666,6663,16666,6660666 等等。
# 古代典籍中经常用“第 k 小的魔鬼的数”来指代这些数，这给研究人员带来了极大的不便。
# 现在请编写一个程序，可以实现输入 k，输出对应的魔鬼数。
# 1<=k<=5e7

from functools import lru_cache

# 二分+数位dp


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, isLimit: bool, count: int, isOk: bool) -> int:
        """当前在第pos位,isLimit表示是否贴合上界"""
        if pos == len(nums):
            return isOk

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            nextCount = 0 if cur != 6 else (count + 1)
            res += dfs(pos + 1, (isLimit and cur == up), nextCount, isOk or (nextCount == 3))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, True, 0, False)


T = int(input())
for _ in range(T):
    k = int(input())

    def countNGT(mid: int) -> int:
        """[1,mid]里面有多少个含三个连续的6的数字"""
        return cal(mid) - cal(0)

    left, right = 0, int(7e10)
    while left <= right:
        mid = (left + right) // 2
        if countNGT(mid) < k:
            left = mid + 1
        else:
            right = mid - 1
    print(left)
