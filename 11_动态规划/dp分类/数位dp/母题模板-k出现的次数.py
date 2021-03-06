#  统计区间[L,R]出现0123456789的各个数字总次数
# 每个结果包含十个用空格隔开的数字，第一个数字表示 0 出现的次数，第二个数字表示 1 出现的次数，以此类推。

# 1≤a,b≤231−1,
# 1≤N<100


from functools import lru_cache




@lru_cache(None)
def cal(upper: int, queryDigit: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, count: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位，出现次数为count，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == len(nums):
            return count

        res = 0
        up = nums[pos] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos + 1, count, True, (isLimit and cur == up))
            else:
                res += dfs(pos + 1, count + (cur == queryDigit), False, (isLimit and cur == up))
        return res

    nums = list(map(int, str(upper)))
    return dfs(0, 0, True, True)


while True:
    left, right = sorted(map(int, input().split()))
    if left == right == 0:
        break
    res = []
    for i in range(10):
        res.append(str(cal(right, i) - cal(left - 1, i)))
    print(' '.join(res))

# 6902 6902 6904 6984 7011 7523 16002 11851 5952 6574
# 1523 1523 1560 1619 2261 2523 7662 2533 2533 1923
# 16373 17028 25374 26275 26275 26267 18017 16217 16165 16264
# 17742 17742 12263 7844 7270 7513 7742 7742 7743 7742
# 2566 2571 2661 2251 1676 3436 5151 1571 2351 2566
# 2833 2398 2945 3394 3371 10942 3833 3398 3398 3398
# 31646 36726 41656 41330 40735 41655 41629 41549 41545 35299
# 24799 31846 34846 35800 35742 35143 34699 30145 24699 24746
# 17842 17842 18396 18842 26454 28952 28923 28843 27183 17843
# 8379 8379 9103 9153 8478 16190 18452 12046 8371 8379

# 6902 6902 6904 6984 7011 7523 16002 11851 5952 6574
# 1523 1523 1560 1619 2261 2523 7662 2533 2533 1923
# 16373 17028 25374 26275 26275 26267 18017 16217 16165 16264
# 7742 17742 12263 7844 7270 7513 7742 7742 7743 7742
# 2566 2571 2661 2251 1676 3436 5151 1571 2351 2566
# 2833 2398 2945 3394 3371 10942 3833 3398 3398 3398
# 31646 36726 41656 41330 40735 41655 41629 41549 41545 35299
# 24799 31846 34846 35800 35742 35143 34699 30145 24699 24746
# 17842 17842 18396 18842 26454 28952 28923 28843 27183 17843
# 8379 8379 9103 9153 8478 16190 18452 12046 8371 8379
