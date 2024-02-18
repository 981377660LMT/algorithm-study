https://ei1333.github.io/library/string/palindromic-tree.hpp
https://math314.hateblo.jp/entry/2016/12/19/005919
https://baobaobear.github.io/post/20200416-pam/
https://mojashi.hatenablog.com/entry/2017/07/17/155520
https://miti-7.hatenablog.com/entry/2016/03/14/195721
https://zhuanlan.zhihu.com/p/581946378

- 性质
  它的状态表示原字符串的一个回文子串
  它的转移表示在一个子串前后添加一个相同字符后可以得到新回文子串
  它的 fail 指针指向一个子串的最长回文真后缀
- 作用

1. 统计每种回文的出现次数
2. 以第 i 个字符结尾的不同长度的回文串
3. `在线`构造回文树

- API:

1. 在线查询以每个位置结尾的最长回文串的长度
2. 在线查询以每个位置结尾的回文串个数

---

https://www.bilibili.com/video/BV1r54y1S7C1
回文树的关键技术：奇偶字典树+后缀链跳跃
![Alt text](image.png)
