回溯去重的几个类型

全排列元数组允许有相同元素，但是每个元素只能使用一次:
使用 **visited[i] 数组**来保证同一个数字不能使用两次，并且排序元数组后排除重复元素一前一后的情况

```JS
    for (let i = 0; i < nums.length; i++) {
      // 同一个数字不能用两次
      if (visited[i]) continue
      // 同样值的数字不能用两次 之前已经看过了,那就要防止[前一个1，后一个1]和[后一个1，前一个1]这种情况的出现 这里约定只允许[前一个1，后一个1]
      if (i > 0 && nums[i] === nums[i - 1] && visited[i - 1]) continue
      visited[i] = true
      bt(path.concat(nums[i]), visited)
      visited[i] = false
    }
```

组合求和元数组**无相同元素但是每个元素可以使用多次**:
使用 index 记录选择的元素并保证 **index 不减**

```JS
    for (let i = index; i < len; i++) {
      const next = candidates[i]
      path.push(next)
      bt(path, sum + next, i)
      path.pop()
    }
```

组合求和元数组/子集加强版**有相同元素但是每个元素只能使用一次**:
排序;
使用 index 每次加 1 保证不取到相同元素;
并且限制每个重复的元素只能在开头第一个被使用

```JS
 for (let i = index; i < len; i++) {
       // 限制每个重复的元素只能在开头第一个被使用
      if (i > index && candidates[i] === candidates[i - 1]) continue

      const next = candidates[i]
      path.push(next)
      // 注意这里i+1限制不能取到重复的元素
      bt(path, sum + next, i + 1, visited)
      path.pop()
    }
```

1. push 结果 时要 push 浅拷贝
2. 如果递归参数是引用要回溯 pop
3. 要求元素不能重用时要**排序**并带上 index(start)参数并且限制每个重复的元素只能在开头第一个被使用
   **如果不允许排序(491. 递增子序列),则不能限制每个重复的元素只能在开头第一个被使用,而是在每一轮 bt 中创建一个 set 记录用过的数字,即记录已经使用过的 next**

```JS
  const visited = new Set<number>()
  for (let i = start; i < nums.length; i++) {
    // if (i > start && nums[i] === nums[i - 1]) continue
    if (visited.has(nums[i])) continue
    visited.add(nums[i])
    path.push(nums[i])
    bt(i + 1, path)
    path.pop()
  }
```

**回溯法解决的问题都可以抽象为树形结构**
回溯法，一般可以解决如下几种问题：

组合问题：N 个数里面按一定规则找出 k 个数的集合
切割问题：一个字符串按一定规则有几种切割方式
子集问题：一个 N 个数的集合里有多少符合条件的子集
排列问题：N 个数按一定规则全排列，有几种排列方式
棋盘问题：N 皇后，解数独等等
//////////////////////////////////////////////////////
**尝试删除每一个位置的元素** 模板:
例题见 1_stack\括号\301. 删除无效的括号.ts

```JS
for (const item of queue) {
  for (let i = 0; i < item.length; i++) {
    // 尝试去掉每一个位置的元素
      const next = item.slice(0, i) + item.slice(i + 1)
      if (!visited.has(next)) {
        queue.push(next)
        visited.add(next)
    }
  }
}

```
