## 离散化，就是当我们只关心数据的大小关系时，用排名代替原数据进行处理的一种预处理方法

常见方式

1. **哈希表映射查询**：使用哈希表将所有数手动建立映射关系离散化。
   查询时用 `mapping里的关系`

```Python
sl = sorted(set(nums))
mapping = {sl[i]: i + 1 for i in range(len(sl))}
newNums = [mapping[num] for num in nums]
```

2. **有序数组二分查询**：查询时用 `bisect_right(allNums, num) + 1`

```Python
sl = sorted(set(nums))
newNums = [bisect_right(sl, num) + 1 for num in nums]
```

注意：离散化适合用于`单点查询`
`区间查询使用离散化`会很麻烦 一般用线段树/树状数组

- Discretize2D，注意 api 与 1D 统一
- 一维、二维前缀和，支持离散化
  (presumDense, presumSparse, presumDense2D, presumSparse2D)
  `api 设计类似 RectangleSum`
