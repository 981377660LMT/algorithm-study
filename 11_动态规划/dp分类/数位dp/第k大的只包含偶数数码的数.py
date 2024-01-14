# 第k大的只包含偶数数码的数(0,2,4,6,8)
# 1<=k<=1e12
# https://atcoder.jp/contests/abc336/tasks/abc336_c
# 类似:https://atcoder.jp/contests/abc186/tasks/abc186_c


from functools import lru_cache


def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
    assert 2 <= radix <= 10
    if num < 0:
        return "-" + toString(-num, radix)
    if num == 0:
        return "0"
    res = []
    while num:
        div, mod = divmod(num, radix)
        res.append(str(mod))
        num = div
    return "".join(res)[::-1] or "0"


def evenDigits1(k: int) -> int:
    """转换为5进制."""
    k -= 1
    s = toString(k, 5)
    return 2 * int(s)


def evenDigits2(k: int) -> int:
    """二分+数位dp"""

    ok = set([0, 2, 4, 6, 8])

    def check(mid: int) -> bool:
        @lru_cache(None)
        def dfs(pos: int, hasLeadingZero: bool, isLimit: bool) -> int:
            """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
            if pos == m:
                return int(not hasLeadingZero)

            res = 0
            up = nums[pos] if isLimit else 9
            for cur in range(up + 1):
                if cur not in ok:
                    continue
                if hasLeadingZero and cur == 0:
                    res += dfs(pos + 1, True, (isLimit and cur == up))
                else:
                    res += dfs(pos + 1, False, (isLimit and cur == up))
            return res

        nums = list(map(int, str(mid)))
        m = len(nums)
        res = dfs(0, True, True) + 1
        dfs.cache_clear()
        return res >= k

    left, right = 0, int(1e19)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1
    return left


if __name__ == "__main__":
    k = int(input())
    assert evenDigits1(k) == evenDigits2(k)
    print(evenDigits1(k))
