https://www.acwing.com/blog/content/7944/

DFS 做法才是数位 DP 的正解

数位 DP 问题一般给定一个区间 [L,R]，问区间满足的条件的数有多少个。数字范围很大。
可以利用前缀和来求解答案

数位 dfs 模板

1. 把数字拆成 nums
2. `dfs(pos,isLimit,*rest)` 开始 dfs(len(nums),True,...)
3. 结束条件为 pos==0
4. 每次选下一位数先确定 up `up = nums[pos - 1] if isLimit else 9`

```C++
// 假设数字 x 位数为 a1⋯an
// 两个必带参数：
// pos:pos: 表示数字的位数
// isLimit: 可以填数的限制（无限制的话(isLimit=True) 0∼9随便填，否则只能填到up）

// 四个选带参数：
// pre:表示上一个数是多少
// hasLeadingZero:前导零是否存在，表示是否开始选数了
// curSum:搜索到当前所有数字之和
// count: 某个数字出现的次数
int dfs(int pos, int pre, int lead, bool isLimit) {
    if (!pos) {
        边界条件
    }
    if (!limit && !lead && dp[pos][pre] != -1) return dp[pos][pre];
    int res = 0, up = limit ? a[pos] : 无限制位;
    for (int i = 0; i <= up; i ++) {
        if (不合法条件) continue;
        res += dfs(pos - 1, 未定参数, lead && !i, limit && i == up);
    }
    return limit ? res : (lead ? res : dp[pos][sum] = res);
}
int cal(int x) {
    memset(dp, -1, sizeof dp);
    len = 0;
    while (x) a[++ len] = x % 进制, x /= 进制;
    return dfs(len, 未定参数, 1, 1);
}
signed main() {
    cin >> l >> r;
    cout << cal(r) - cal(l - 1) << endl;
}

```

`经典计数问题`

```Python

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


class Solution:
    def digitsCount(self, d: int, low: int, high: int) -> int:
        return cal(high, d) - cal(low - 1, d)
```
