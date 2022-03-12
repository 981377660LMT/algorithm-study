Adjacent Swaps

1. 贪心，从后面找第一个相同的位置搬过来
   `1505. 最多 K 次交换相邻数位后得到的最小整数`
   `1850. 邻位交换的最小次数`
   `2193. 得到回文串的最少操作次数`

2. 每个数不断减 1 求前缀和，然后计算到中位数的距离之和
   `ZJ7 字母交换-在至多邻位交换m次之后，字符串中最多有多少个连续的位置上的字母相同`
   `1703. 得到连续 K 个 1 的最少相邻交换次数`

   ```Python
   def minMoves(indexes: List[int], target: int) -> int:
      """得到连续 target 个 相同字符 的最少相邻交换次数"""
      indexes = [num - i for i, num in enumerate(indexes)]
      preSum = [0] + list(accumulate(indexes))

      res = int(1e20)
      # left+k-1<len(ones)
      for left in range(len(indexes) + 1 - target):
          right = left + target - 1
          mid = (left + right) >> 1
          leftSum = indexes[mid] * (mid - left) - (preSum[mid] - preSum[left])
          rightSum = preSum[right + 1] - preSum[mid + 1] - indexes[mid] * (right - mid)
          res = min(res, leftSum + rightSum)

      return res
   ```

3. 相邻交换消除逆序对
