有多少段区间的和>=target
`差分 dict 适用于区间个数限定(扫描个数被限定了)，值域不定的情况`

```Python
class Solution:
    def solve(self, walks, target):
        pos = 0
        diff = defaultdict(int)
        for num in walks:
            if num > 0:
                diff[pos] += 1
                diff[pos + num] -= 1
            elif num < 0:
                diff[pos + num] += 1
                diff[pos] -= 1
            pos += num


        # 一般都这么写，维护前一个位置和当前变量
        res = 0
        curSum = 0
        prePos = 0
        for pos, cur in sorted(diff.items()):
            if curSum >= target:
                res += pos - prePos
            curSum += cur
            prePos = pos

        return res
```
