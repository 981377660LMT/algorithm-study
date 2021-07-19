# js-algorithm

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

空间复杂度指的是存储了多少个变量(例如二维矩阵就是 O(n^2))

vscode 断点调试:
快捷键|名称|作用
-----|-----|-----
F5|继续|跳到下一个断点处
F10|单步跳过|一行一行执行代码
F11|单步调试|进入函数
Shift+F11|单步跳出|跳出函数

1. 栈
   - 应用场景:十进制转二进制，判断字符串括号是否有效，函数调用堆栈
   - 单调栈适合的题目是求解**下一个大于 xxx 或者下一个小于 xxx** 这种题目
   - Stack over flow :函数的临时变量是存储在栈区;没有出口的递归调用栈会超出最大深度
   - Maximum call stack size exceeded
2. 队列
   - 应用场景:JS 异步中的任务队列，计算最近请求次数
3. 链表(有时候直接用首部 Node 来表示链表)

   - 元素储存不连续，用 next 连在一起
   - 应用场景:JS 异步中的任务队列，计算最近请求次数，原型链，循环找节点父元素
   - 基本方法:

   1. 建立 newNode 最后返回 newNode.next
   2. 使用 p1=node1 p2=node2 指针来操作,移动链表需要用指针(代理)
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

4. 集合
   - 应用场景:JS 异步中的任务队列，计算最近请求次数
5. 字典
   - 保存状态
   - js 的 map 是有序字典 => LRU
6. 树(分层数据抽象模型)

   - Object+Array
   - 遍历方法 dfs(大多数情况) bfs(层序遍历)
   - 递归 (最长同值路径，相同的树，验证树的某些性质等...)
   - 应用场景:DOM 树，级联选择，树形控件
     深度优先遍历 JSON

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
9. 排序算法 (方法:先写一轮，再写多轮)
   js 本身的 sort 复杂度 nlog(n)

   - 冒泡 n^2
   - 选择 n^2
   - 插入 n^2
   - 归并 `n*log(n)` 火狐浏览器
   - 快速 `n*log(n)`
     记得添加边界条件

   ```JS
   if (arr.length <= 1) return arr
   ```

10. 搜索算法(找下标)

    - 顺序搜索 n
    - 二分搜索 log(n) 前提是数组有序

11. 分治思想:递归解决小问题 eg:归并排序 快速排序 反转二叉树 判断二叉树的性质
    - 分
    - 递归
    - 合

```JS
对于数的递归：
// 递归的终点
  if (!t1 && !t2) return true
```

12. 动态规划

    - 构造数列递推式

13. 贪心算法:期盼局部最优达到全局最优,但是并不一定全局最优 分饼干，买股票

    - 珍珑棋局
    - 构造数列递推式
    - an algorithm makes the optimal choice **at each step as going forward**.

14. 回溯算法(BackTrace):渐进式解决问题
    - 遇到岔路：走一条**路**没走通，走回来
    - 递归模拟所有的路
    - 例子:输出全排列
    - 可以引入 **memo** 来辅助回溯

回顾与总结
与前端最重要的是**链表**与**树**(JSON/对象)
搞清楚**特点**与**应用场景**
并且用**不同语言实现**
分析时间/空间复杂度
刷 300 道以上 leetcode
分类总结套路

15. 并查集(union find)

    - 并查集适用于合并集合，查找哪些元素属于同一组，**有相同根的元素为一组**
    - 如果 a 与 b 属于同一组，那么就 uniona 与 b;以后找元素属于哪个组时只需要 find(这个元素)找到属于哪个根元素
    - 例如很多邮箱都属于同一个人，就 union 这些邮箱；回来分类时找根邮箱即可

16. 滑动窗口(sliding window)
    - 适合解决**定长**(子)数组的问题
    - 减少 while 循环

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
