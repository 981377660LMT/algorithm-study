# 每次你可以把一段连续的木版涂成一个给定的颜色，后涂的颜色覆盖先涂的颜色。
# 例如第一次把木版涂成RRRRR，第二次涂成RGGGR，第三次涂成RGBGR，达到目标。 用尽量少的涂色次数达到目标。
# 1≤n≤50。


from functools import lru_cache


INF = int(1e18)


@lru_cache(None)
def dfs(left: int, right: int) -> int:
    """ "[left,right]这一段还要涂色"""

    if left > right:
        return 0
    if left == right:
        return 1

    # 涂一次可以涂一整个区间，l或者r的那一次染色可以拿来帮助另一个染色
    if word[left] == word[right]:
        return min(dfs(left + 1, right), dfs(left, right - 1))

    res = INF
    for i in range(left, right):
        res = min(res, dfs(left, i) + dfs(i + 1, right))
    return res


word = input()
res = dfs(0, len(word) - 1)
dfs.cache_clear()
print(res)
