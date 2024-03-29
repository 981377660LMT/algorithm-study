集训队论文笔记

## 2013 前

### 数值运算的优化:

- 消除除法
- 减少取模运算

```cpp
inline int mod(int a){a%=M;if(a<0)a+=M;return a;};
```

- 减少浮点数除法

```cpp
 x=a/b;
 y=c/d;
 z=e/f;
```

可以变为

```cpp
t=1/(b*d*f);
x=a*d*f*t;
y=c*b*f*t;
z=e*b*d*t;
```

## 2013

### 浅谈数据结构题的几个非经典解法

- cdq 分治(一种时间分治算法)
  对于“修改独立，允许离线" 的数据结构题，
  可以用一个 log 的代价，去掉原问题中的动态修改，变成没有动态修改的`静态简化版`问题
  按照中序遍历序进行分治才能保证每一个修改都是严格按照时间顺序执行的。
- 二进制分组(类似查询分块，但是二进制分组更优)
  对于“修改独立，强制在线" 的数据结构题，
  可以用一个 log 的代价，去掉原问题中的动态修改，变成没有动态修改的`静态简化版`问题
- 整体二分
  我们的处理复杂度，绝对不能和序列总长度 n 线性相关，只能和`当前处理序列区间`的长度相关
  - 多个查询，子矩阵第 k 大
  - 多个查询，区间第 k 大

## 2014

### 精细地实现程序-浅谈 OI 竞赛中的常数优化

- 随机数：系统的随机生成函数很慢，自己写随机函数可以用 xor shift。
- 动态内存静态化：系统的内存分配函数多次调用时速度很慢，解决方案是估计需要的内存量并提前开辟内存池，手工实现内存分配。
- 内存访问连续性：尽量防止数组非最后一维的连续变动，如有需要可以交换数组的两维。

## 2015

### 浅谈分块在一类在线问题中的应用

- 序列分块
- 树分块
- 定期重构(将操作分块，每隔若干个操作后就重建一次状态，查询的时候根据上一次重建和期间的所有修改计算答案)
