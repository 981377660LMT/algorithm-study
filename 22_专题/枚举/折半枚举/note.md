meet in middle 分割处理的方式可以把数量级由 2^n 降到 2^n/2 级别

# 关于“最接近目标值的子序列和”的这一类问题：

1. 如果数组长度特别大，但是值小 (sum<=10^7)，我们可以使用**背包问题**的方式来解决，
   `1049. 最后一块石头的重量 II`
   dp[i]表示是否能组成容量为 i 的背包
2. 但是如果数组长度不大(n<=20)，但是数值特别大的话，使用**折半的 DFS**。(数组最大长度是 40，直接 dfs 会超时，选择双向 dfs，以空间换时间)
   `1755. 最接近目标值的子序列和`
   `5897. 将数组分成两个数组并最小化数组和的差`

   1. 将数组分成两部分
   2. 求出两个数组的所有子序列和并排序
   3. 有序的 twoSum 问题 => 双指针

   数据范围 n :
   n <= 20, 可以暴力，2^20
   n <= 40, 不能直接暴力，但是可以对半来，变成 2^20
   数据规模<= 10 算法可接受时间复杂度 O(n!)
   数据规模<= 20 算法可接受时间复杂度 O(2^n)

https://leetcode-cn.com/problems/partition-array-into-two-arrays-to-minimize-sum-difference/solution/zui-jie-jin-mu-biao-zhi-de-zi-xu-lie-he-m0sq3/

## 枚举的两种形式

1.  k 进制枚举+check (如果要处理有哪些数，复杂度`n*k^n`)

```Python
for state in range(3 ** n):
    curSum, absSum = 0, 0
    for i in range(n):
        mod = (state // (3 ** i)) % 3
        if mod == 1:
            curSum += nums[i]
            absSum += nums[i]
        elif mod == 2:
            curSum -= nums[i]
            absSum += nums[i]
    res[curSum] = max(res.get(curSum, 0), absSum)
```

2. dfs 枚举(如果要处理有哪些数，复杂度`k^n`)

```Python
 def dfs(index: int, curSum: int, absSum: int) -> None:
    if index == n:
        res[curSum] = max(res.get(curSum, 0), absSum)
        return
    dfs(index + 1, curSum, absSum)
    dfs(index + 1, curSum - nums[index], absSum + nums[index])
    dfs(index + 1, curSum + nums[index], absSum + nums[index])
```

3. 状压 dp 枚举子序列和(复杂度`2^n`)

```JS
function getSubArraySum(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(1 << n).fill(0)

  // 外层遍历数组每个元素，遍历到时，表示取该元素
  for (let i = 0; i < n; i++) {
    // 内层遍历从0到外层元素之间到每一个元素，表示能取到的元素，由于前面的结果已经计算过，因此可以直接累加
    for (let pre = 0; pre < 1 << i; pre++) {
      res[pre + (1 << i)] = res[pre] + nums[i]
    }
  }

  return res
}
```

---

折半搜索的两种形式：

- dfs + 剪枝, 推荐使用
- bfs(双向 bfs)
