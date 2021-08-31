# js-algorithm

JS,TS 刷题时大数问题可用 BigInt 或者求余解决

分析算法，找到算法的瓶颈部分，然后选取合适的数据结构和算法来优化

数据结构：线性/树/图
可视化
https://visualgo.net/zh
JavaScript 实现数据结构与算法

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

0. 数组
   - 数组是存放在连续内存空间上的相同类型数据的集合。
   - 在 JavaScript 中，可以在数组中保存不同类型值，可以越界访问数组，并且数组可以动态增长，不像其它语言，例如 C，创建的时候要决定数组的大小，如果数组满了，就要重新申请内存空间。
1. 栈
   - 应用场景:十进制转二进制，判断字符串括号是否有效，函数调用堆栈
   - 单调栈适合的题目是求解**下一个大于 xxx 或者下一个小于 xxx** 这种题目
   - Stack over flow :函数的临时变量是存储在栈区;没有出口的递归调用栈会超出最大深度
   - Maximum call stack size exceeded
2. 队列
   - 处理广度优先遍历(树:层序遍历/图:无权图的最短路径)
   - 优先队列(每个元素有优先级，优先级高的先出队:操作系统任务调度)
   - 应用场景:JS 异步中的任务队列，计算最近请求次数
3. 链表(有时候直接用首部 Node 来表示链表)

   - 基本操作:遍历/插入/删除
   - 基本方法:虚拟 dummy,穿针引线(p1p2),快慢指针,先穿再排再判空
   - 元素储存不连续，用 next 连在一起
   - 应用场景:JS 异步中的任务队列，计算最近请求次数，原型链，循环找节点父元素
   - 4 个基本方法:

   1. 建立 newNode 最后返回 newNode.next(**设立链表的虚拟头节点**)
      链表问题首先我们看返回值，**如果返回值不是原本链表的头的话，我们可以使用虚拟节点。**
   2. 使用 headP1=head1 headP2=head2 指针来操作,移动链表需要用指针(**穿针引线**)
   3. 快慢指针寻找中间/第 n 位元素
   4. 先穿再排再判空

      **数组读，链表写**
      基本操作:遍历链表/删除节点/**获取 json 某个节点值**

   ```JS
   const json = {
   a: { b: { c: 1 } },
   d: { e: 2 },
   }

   const path = ['a', 'b', 'c']

   // 与遍历链表异曲同工
   let p: { [k: string]: any } = json
   path.forEach(key => {
   p = p[key]
   })

   console.log(p)

   ```

   - | 操作       | 单链表 | 双链表 |
     | ---------- | ------ | ------ |
     | append     | 1      | 1      |
     | appendleft | 1      | 1      |
     | remove     | n      | 1      |
     | find       | n      | 1      |

     - 画图/注意边界条件/类比数组
     - 链表的价值就在于其不必要求物理内存的连续性，以及对插入和删除的友好

4. 集合
   - 应用场景:JS 异步中的任务队列，计算最近请求次数
   - 由于二分搜索树不能 insert 相同元素，所以可以用来实现集合
   - 查找**有无**
   - 有序集合(搜索树实现)与无序集合(哈希表实现)
   - 客户 ip 统计
   - 词汇量统计
5. 字典
   - 保存状态
   - js 的 map 是有序字典 => LRU
   - 查找**对应关系**
6. 树(分层数据抽象模型)

   - Object+Array
   - 递归的中止条件(没有则会报栈溢出)+递归过程
   - 遍历方法 dfs(大多数情况) bfs(层序遍历)
   - 递归 (最长同值路径，相同的树，验证树的某些性质等...)
   - 应用场景:DOM 树，级联选择，树形控件
     深度优先遍历 JSON
   - 如果想要在遍历时对节点**产生副作用**，需要再 dfs 中**传递 root 的引用**(例如 antd 添加 key,为每个节点添加 children)

     ```TS
     interface Dict<V = any> {
     [key: string]: V
     }

     // 深度优先
     const json = {
     a: { b: { c: 1 } },
     d: [1, 2],
     }

     const dfs = (n: Dict, memo: string[] = []) => {
     console.log(n, memo)
     Object.keys(n).forEach(k => dfs(n[k], memo.concat(k)))
     }

     dfs(json)
     export {}

     ```

     - 关于递归:尽量少维护全局状态，把状态作为递归参数传递到递归内部，每层递归自己维护当前层的状态
     - 完全二叉树:除了最后一层都满了，最后一层所有节点都在最左侧(堆就是完全二叉树)
     - 满二叉树: 所有层的结点数都满了
     - **二叉搜索树(BST)**:每个结点的键值大于左孩子，小于右孩子；左右孩子同理
       基本操作(logn):插入/查找/删除/最大最小/前驱后继/某个元素的排名/第 k 大小元素
       **中序遍历**可以得到递增的序列
     - **平衡二叉树(AVL)**是对二叉搜索树的优化；二叉搜索树在极端条件下会退化为链表(logn=>n)
     - 完全随机的数据，普通的二分搜索树很好(不平衡)
     - 读多写少的数据，AVL 树很好
     - 读少写多的数据，红黑树很好(红黑树牺牲了平衡性)
     - **红黑树**的统计性能更优(crud)
     - java 中的 treemap 和 treeset 基于红黑树实现

7. 图
   - 二元关系，道路航班
   - Object+Array **map 实现的邻接矩阵，邻接表,入度数组**
   - 深度优先广度优先 dfs bfs(queue)
   - 邻接矩阵 number[][],邻接表 Map<number,numner[]>
8. 堆
   - 特殊的完全二叉树(完全填满，最后一层从左向右填充)
   - 最大堆/最小堆
   - 数组表示堆
   - 位置为 index 的左侧子节点的位置是 `2*index+1`
   - 位置为 index 的侧右子节点的位置是 `2*index+2`
   - 父节点位置为`(index-1)/2`
   - 作用:时间复杂度 O(1)找出第 K 个最大最小值
   - 第 K 个最大元素就是 K 位高手里最弱的哪一个(最小堆堆顶) 时间复杂度 `k*log(n)`
   - **堆化操作(heapify)的时间复杂度是 O(N)**
9. 排序算法 (方法:先写一轮，再写多轮)
   js 本身的 sort 复杂度 nlog(n)
   基础:冒泡/选择/插入/希尔
   高阶:归并/堆/快排

在 V8 引擎中， 7.0 版本之前，数组长度小于 10 时， Array.prototype.sort() 使用的是插入排序，否则用快速排序。
在 V8 引擎 7.0 版本之后就舍弃了快速排序，**因为它不是稳定的排序算法**，**在最坏情况下，时间复杂度会降级到 O(n2)**。
而是采用了一种混合排序的算法：**TimSort** 。
这种功能算法最初用于 Python 语言中，严格地说它不属于以上 10 种排序算法中的任何一种，属于一种混合排序算法：
在数据量小的子数组中使用插入排序，然后再使用归并排序将有序的子数组进行合并排序，时间复杂度为 O(nlogn) 。

- 冒泡 n^2
- 选择 n^2
- 插入 n^2
- 归并 `n*log(n)` 火狐浏览器
- 快速 `n*log(n)`
  记得添加边界条件

```JS
if (arr.length <= 1) return arr
```

topK 问题解法:
快速排序 O(n) 要求数据一次性给出 (见对应文件夹)
堆(优先队列)O(nlogk) 但是数据很大/数据是流时 用堆可以实时更新 (维护容量为 k 的最大最小堆)

**只有归并排序要 O(n)空间复杂度，其他算法都可以原地完成**
**快排/堆排序/归并排序中只有归并是稳定的**(需要 mergeTwo 过程不移动相等的元素)

当数据量小于 10 时，使用插入排序(**稳定的**)优化归并排序

如果元素是一维的，稳定性没有意义
如果元素位置关系到其他属性，稳定性才有意义

1.  搜索算法(找下标)

    - 顺序搜索 n
    - 二分搜索 log(n) 前提是数组有序

2.  分治思想:递归解决小问题 eg:归并排序 快速排序 反转二叉树 判断二叉树的性质
    - 分
    - 递归
    - 合

```JS
对于数的递归：
// 递归的终点
  if (!t1 && !t2) return true
```

12. 动态规划

    - 初始化 dp 数组(**dp 的维度**取决于状态个数)
    - 构建状态转移方程
    - 空间复杂度的优化:如果第 i 行元素只依赖于 i-1 行的元素，可以不用整个数组而是只用一行

13. 贪心算法:期盼局部最优达到全局最优,但是并不一定全局最优 分饼干，买股票

    - 珍珑棋局
    - 构造数列递推式
    - an algorithm makes the optimal choice **at each step as going forward**.

14. 回溯算法(BackTrace):渐进式解决问题
    - 回溯法是暴力解法的一个主要实现手段/是**经典**人工智能的基础
    - 树形问题
    - 遇到岔路：走一条**路**没走通，走回来
    - 剪枝:预先排除不必要的路径
    - 递归模拟所有的路
    - 例子:输出全排列
    - 可以引入 **memo** 来辅助回溯
    - 如果最后保存结果到数组中，需要在 bt 完之后 **pop**,并且注意 bt 的终点需要**浅拷贝**

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

4. todo

- 力扣加加
- lucifer 的博客
- 哥谭 blog
- top100
- 剑指 offer
- 程序员面试金典
- 力扣加加的题
- 编程之法
- 大前端宝典
- go 收藏
- go 算法
- go mooc

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
