https://github.com/upupming/algorithm/blob/master/template-typescript.md
https://github.com/harttle/contest.js

## todo

- 组件参考 elementUI
- code up 高频题

JS,TS 刷题时大数问题可用 BigInt 或者求余解决
遇到 1e9+7 的题目最好全部用 BigInt，避免产生溢出问题

如果数组长度不用变化，则可使用 **new TypedArray(分配连续内存的数组，长度固定)**
防止数组索引越界 grid[row+1]?.[col]

分析算法，找到算法的瓶颈部分，然后选取合适的数据结构和算法来优化

数据结构：线性/树/图
可视化
https://visualgo.net/zh
JavaScript 实现数据结构与算法

在 leetcode 中使用 JS
datastructures-js/queue 官方可用的堆队列库

```JS
/**
 * @description 特点
 * @description 应用场景
 * @description 时间复杂度
 * @description 空间复杂度
 */
```

时间复杂度:普通算法看最差(有一组数据 100%恶化) 随机算法看期望(没有一组数据能 100%恶化)
例如:单路快排 O(n^2) 双路快排 O(nlogn)

数据规模
10^9 是可接受的规模 10^9
C 语言大约是 4 秒
**要想在 1s 内解决问题，需要是 10^8 次的计算**
空间复杂度指的是存储了多少个变量(例如二维矩阵就是 O(n^2))

空间换时间，时间换空间
降维打击:数学方法

vscode 断点调试:
快捷键|名称|作用
-----|-----|-----
F5|继续|跳到下一个断点处
F10|单步跳过|一行一行执行代码
F11|单步调试|进入函数
Shift+F11|单步跳出|跳出函数

回顾与总结
与前端最重要的是**链表**与**树**(JSON/对象)
搞清楚**特点**与**应用场景**
并且用**不同语言实现**
分析时间/空间复杂度
刷 300 道以上 leetcode
分类总结套路

1.  并查集(union find)

    - 并查集适用于合并集合，查找哪些元素属于同一组，**有相同根的元素为一组**
    - 如果 a 与 b 属于同一组，那么就 uniona 与 b;以后找元素属于哪个组时只需要 find(这个元素)找到属于哪个根元素
    - 例如很多邮箱都属于同一个人，就 union 这些邮箱；回来分类时找根邮箱即可

2.  滑动窗口(sliding window)
    - 适合解决**定长**(子)数组的问题
    - 减少 while 循环
    - 有时需要在滑动窗口中**做记录**

```TS
class UnionFind<T> {
  private map: Map<T, T>

  constructor() {
    this.map = new Map()
  }

  // key不存在时(find返回key自身)，会设置key1指向k2
  union(key1: T, key2: T) {
    const root1 = this.find(key1)
    const root2 = this.find(key2)
    if (root1 !== root2) {
      this.map.set(root1, root2)
    }
  }

  // key不存在时，返回key自身
  find(key: T) {
    while (this.map.has(key)) {
      key = this.map.get(key)!
    }
    return key
  }
}

```

1.  字符串杂题:优美的方法是**转数组**然后 map filter reduce((pre,cur,index,arr)=>...)

2.  子数组问题(subArray)：双指针(滑动窗口)+Map 保存状态/画折线图(与求和相关的，例如最大和子串)

3.  python 内置

算法:思考问题的方式
对一组数据排序：快速排序算法(nlogn)
不能忽略算法使用的环境
数组的特征?大量重复的元素:**三路快排**
是否近乎有序?**插入排序**
取值范围有限(学生成绩)?**计数排序**
排序有没有额外要求?是否需要稳定排序?**归并排序**
数据存储状况?链表存储?**归并排序**
数据量太大，内存不够?**外排序**

代码规范，容错性...
不只是对错

**参与项目至关重要**
实习
自己开发的小应用
代码整理

没有思路:

1. 简单测试用力尝试
2. 暴力解法

对于基本问题(最大堆/归并快排)，白板编程

查找问题为什么不全用哈希表(O(1))?
因为哈希表只能查对应关系，**失去了顺序性**,例如某个排位的元素/数据集中最大最小值
可以用二叉搜索树(增删改查全为 O(logn))

一定要重视算法与数据结构 每天坚持刷题
归纳、分类、与总结

| 数据规模 | 算法可接受时间复杂度 |
| -------- | -------------------- |
| <= 10    | O(n!)                |
| <= 20    | O(2^n)               |
| <= 100   | O(n^4)               |
| <= 500   | O(n^3)               |
| <= 2500  | O(n^2)               |
| <= 10^6  | O(nlogn)             |
| <= 10^7  | n                    |
| <= 10^14 | O(sqrt(n))           |
| -        | O(logn)              |

面试题困难难度的题目常见的题型有：

DP
设计题
图
游戏

```C++
O(n)的递归
int function1(int x, int n) {
    int result = 1;  // 注意 任何数的0次方等于1
    for (int i = 0; i < n; i++) {
        result = result * x;
    }
    return result;
}
```

内存对齐
只要可以跨平台的编程语言都需要做内存对齐，Java、Python 都是一样的。
为什么会有内存对齐？

- 平台原因：不是所有的硬件平台都能访问任意内存地址上的任意数据，某些硬件平台只能在某些地址处取某些特定类型的数据，否则抛出硬件异常。为了同一个程序可以在多平台运行，需要内存对齐。
- 硬件原因：经过内存对齐后，CPU 访问内存的速度大大提升。

**git 提交规范**

feat: add hat wobble
^--^ ^------------^
| |
| +-> Summary in present tense.
|
+-------> Type: chore, docs, feat, fix, refactor, style, or test.

**feat**: (new feature for the user, not a new feature for build script)
**fix**: (bug fix for the user, not a fix to a build script)
**docs**: (changes to the documentation)
**style**: (formatting, missing semi colons, etc; no production code change)
**refactor**: (refactoring production code, eg. renaming a variable)
**test**: (adding missing tests, refactoring tests; no production code change)
**chore**: (updating grunt tasks etc; no production code change)

```JS
Math.floor 与 ~~(双按位非)的区别
~将input截取为32位(>=2^32就不成立)  谨慎使用  作用是是数字向0取整
很想Math.trunc 但是~~失败时返回0 而Math.trunc失败时返回NaN
```

不像 Math 的其他三个方法： Math.floor()、Math.ceil()、Math.round() ，`Math.trunc() 的执行逻辑很简单，仅仅是删除掉数字的小数部分和小数点`，不管参数是正数还是负数。
Math.hypot() 函数返回所有参数的平方和的平方根
Math.clz32() 函数返回一个数字在转换成 32 无符号整形数字的二进制形式后, 开头的 0 的个数
