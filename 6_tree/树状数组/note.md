区间问题：
简单求区间和，用「前缀和」
多次将某个区间变成同一个数，用「线段树」
其他情况，用「树状数组」

树状数组求解区间问题时
长处就是可以**用 logn 的复杂度动态更新单个值和区间查询**

````TS
// tree的0号位置不存值;初始化「树状数组」，要默认数组是从 1 开始
class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array(size + 1).fill(0)
  }

    /**
   *
   * @param x (离散化后)的树状数组索引
   * @param k 增加的值
   * @description
   * 单点修改：tree数组下标x处及其各个父节点的值加k
   */
  add(x: number, k: number) {
    if (x <= 0) throw Error('查询索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  /**
   *
   * @param x
   * @description
   * 区间查询：返回前x项的值
   */
  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  sumRange(left: number, right: number) {
    return this.query(right) - this.query(left-1)
  }

 /**
   *
   * @param x
   * @returns x 的二进制表示中，最低位的1和后面的0构成的数。
   * @example
   * ```js
   * console.log(3 & -3) // 1
   * console.log(4 & -4) // 4
   *
   * ```
   */
  private lowbit(x: number) {
    return x & -x
  }
}
````

离散化+查询更新
相同模板解决三道困难
[315. 计算右侧小于当前元素的个数
](https://leetcode-cn.com/problems/count-of-smaller-numbers-after-self/solution/shu-zhuang-shu-zu-jie-fa-by-cao-mei-nai-b0zbw/)[327. 区间和的个数
](https://leetcode-cn.com/problems/count-of-range-sum/solution/jstsshu-zhuang-shu-zu-jie-fa-by-cao-mei-0icur/)[493. 翻转对
](https://leetcode-cn.com/problems/reverse-pairs/solution/jstsshu-zhuang-shu-zu-jie-fa-by-cao-mei-uowff/)

# 注意

一般来说需要将用于查询/修改的所有值需要进行离散化(set+并排序，map 映射成树状数组的索引，相对大小不变)

```JS
  const set = new Set(nums)
  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)  // key+1是因为查询和修改的树状数组的索引应为正整数
  }
  // Map(4) { 1 => 1, 2 => 2, 5 => 3, 6 => 4 }
```

附上:逆序对问题三种解法

1. 手动维护一个有序的数组(java-treeset,python-sortedList) 是 O(n^2)还是 O(nlogn)取决于使用的数据结构**插入**的复杂度
2. 归并排序的性质
3. 树状数组终极解决法

**持久化树状数组**
