1. 快排
2. 高精度
   C++ 没有大整数类，所以需要高精度运算
3. 差分
4. Trie 怎么存汉字？
   算法题一般会限制种类
   既然 Trie 可以存整数，那么 Trie 就可以存汉字(化为二进制)
5. STL
   vector:变长数组，倍增思想
   操作系统特性：系统为某个程序分配空间时，所需时间`与空间大小无关`，`与申请次数有关`
   所以边长数组要尽量少申请空间

   set/multiset：平衡树
   insert
   clear
   find
   count
   erase
   lower_bound
   upper_bound

   bitset 压位
   原来 C++ 里`一个布尔需要 1 字节`
   bitset 压缩成`一个字节存八位`

6. 有向无环图(DAG)一定存在拓扑序列
7. NIM 游戏先手必胜 <=> `a1^a2^a3...an 不为 0`
8. dp 常用模型
9. 贪心
10. 时空间复杂度分析
    https://www.acwing.com/file_system/file/content/whole/index/content/1120024/
    dp 时间复杂度=`状态数量*状态转移的计算量`
