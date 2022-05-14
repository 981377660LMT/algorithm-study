from collections import deque, Counter, defaultdict

# 请你通过若干次以下操作将字符串 s 转化成字符串 t ：
# 选择 s 中一个 非空 子字符串并将它包含的字符就地 `升序 排序`。

# 如果可以将字符串 s 变成 t ，返回 true 。否则，返回 false 。

# 总结:
# 在 t[0] 出现的字母一定是由 s 中第一个相同数字转换过来的
# 我们只需要检查比这个字母小的字母都出现在元字符位置之后
# 原因是如果是由后面相同数字转换过来，一定跨不过第一个数字，
# 所以只需要检查 s 中那个数字的前面有没有小于它的(移到左边没有阻拦)，如果有就不能转换！


class Solution:
    def isTransformable(self, s: str, t: str) -> bool:
        ls = list([int(ch) for ch in s])
        lt = list([int(ch) for ch in t])
        cs = Counter(ls)
        ct = Counter(lt)
        if cs != ct:
            return False

        indexes = defaultdict(deque)
        for i, x in enumerate(ls):
            indexes[x].append(i)

        # print(q_s)
        for x in lt:
            firstIndex = indexes[x].popleft()
            for smaller in range(x):
                if indexes[smaller] and indexes[smaller][0] < firstIndex:
                    return False

        return True


print(Solution().isTransformable(s="84532", t="34852"))
# 输出：true
# 解释：你可以按以下操作将 s 转变为 t ：
# "84532" （从下标 2 到下标 3）-> "84352"
# "84352" （从下标 0 到下标 2） -> "34852"

