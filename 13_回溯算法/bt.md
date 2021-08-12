回溯去重的几个类型

全排列元数组允许有相同元素，但是每个元素只能使用一次:
使用 visited[i]来保证同一个数字不能使用两次，并且排序元数组后排除重复元素一前一后的情况

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
