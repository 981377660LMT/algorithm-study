trie+kmp

[<!-- 很复杂 没看了 -->](https://oi-wiki.org/string/ac-automaton/)

AC 自动机是 以 Trie 的结构为基础，结合 KMP 的思想 建立的。

简单来说，建立一个 AC 自动机有两个步骤：

1. 基础的 Trie 结构：将所有的模式串构成一棵 Trie。
2. KMP 的思想：对 Trie 树上所有的结点构造失配指针。
   然后就可以利用它进行多模式匹配了。
