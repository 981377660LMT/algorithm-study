- `选择长为k的字典序最大/小的子序列`
  `221021天池-03. 整理书架.py`
  `1081. 不同字符的最小子序列.py`
  `2030. 含特定字母的最小子序列.py`

  - 思路：字典序的特点，**最大/小的子序列肯定是一个递增/递减的单调栈**

  - 需要保留和丢弃相邻的元素，因此使用栈这种在一端进行添加和删除的数据结构
    字典序最小 栈底肯定是最小的 维护单增的单调栈
    字典序最大 栈底肯定是最大的 维护单增减的单调栈

  - 带有`是否能加` 和 `是否能删` 的 条件限制
    **是否能加**: visited[char] 等于阈值时，最后这个肯定要删的，既然要删，不如早删
    **是否能删**: 剩余的元素 remain[stack[-1]]能够凑齐要求的个数时，可以删栈顶元素

    ```Python

    class Solution:
      def arrangeBookshelf(self, order: List[int], limit: int) -> List[int]:
          stack = []
          visited = defaultdict(int)
          remain = Counter(order)

          need = {key: min(limit, remain[key]) for key in remain}  # !每个元素最后需要的个数

          for num in order:
              # !能不能把这个数字入栈
              if visited[num] == need[num]:
                  remain[num] -= 1
                  continue

              # !能不能删除栈顶元素
              while (
                  stack
                  and stack[-1] > num
                  and remain[stack[-1]] > need[stack[-1]] - visited[stack[-1]]
              ):
                  top = stack.pop()
                  visited[top] -= 1

              stack.append(num)
              visited[num] += 1
              remain[num] -= 1

          return stack
    ```

- `选择长为k的字典序最大/小的子集(可任意排序)`
  k,n<=50
  https://atcoder.jp/contests/abc225/tasks/abc225_f
  从后向前 dp 选还是不选
