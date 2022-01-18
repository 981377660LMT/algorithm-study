# global 很明显就是声明代码块中的变量`使用外部全局的同名变量`

```Python
def dfs(index: int, curCost: int) -> None:
    global res
    if curCost > limit:
        return
    if index == target:

        res += 1
        res %= MOD
        return

    for nextCost in range(curCost, limit + 1, curCost):
        dfs(index + 1, nextCost)


res = 0
```

# nolocal 的使用场景就比较单一，它是`使用在闭包中的`，使变量使用外层的同名变量

```Python
    def maxUniqueSplit(self, s: str) -> int:
        def backtrack(index: int, splitCount: int) -> None:
            if index >= length:
                nonlocal maxSplit
                maxSplit = max(maxSplit, splitCount)
                return

            for i in range(index, length):
                substr = s[index : i + 1]
                if substr not in visited:
                    visited.add(substr)
                    backtrack(i + 1, splitCount + 1)
                    visited.remove(substr)
```
