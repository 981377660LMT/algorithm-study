**多模式匹配(一个文本串加若干模式串)**
AC 自动机:trie+失配指针

(https://oi-wiki.org/string/ac-automaton/)

AC 自动机是 以 Trie 的结构为基础，结合 KMP 的思想 建立的。
AC 自动机`只是有 KMP 的一种思想`，实际上跟一个字符串的 KMPKMP 有着很大的不同。

简单来说，建立一个 AC 自动机有两个步骤：

1. 基础的 Trie 结构：将所有的模式串构成一棵 Trie。
2. KMP 的思想：对 `Trie 树上所有的结点构造失配指针`。
   然后就可以利用它进行多模式匹配了。
   如果一个点 i 的 fail 指针指向 j,那么 root 到 j 的字符串是 root 到 i 的字符串的**最长后缀**。

- 给定 k 个单词和一段包含 n 个字符的文章，求有多少个单词在文章里出现过。

https://zhuanlan.zhihu.com/p/137584630
如果要匹配汉字 需要把主串和所有模式串中涉及的字符 Unicode 离散化

> 参考
>
> - [AC 自动机详解&演示](https://www.bilibili.com/video/BV1iV411B73u?spm_id_from=333.337.search-card.all.click&vd_source=e825037ab0c37711b6120bbbdabda89e)
> - [AC 自动机](https://www.luogu.com.cn/blog/juruohyfhaha/ac-zi-dong-ji)
> - [[算法]轻松掌握 ac 自动机](https://www.bilibili.com/video/BV1uJ411Y7Eg?p=4)
