![suffix tree](23549d9fc1adfd388f3df385c36cb772.jpg)

炫酷后缀树魔术
https://www.luogu.com.cn/blog/EternalAlexander/xuan-ku-hou-zhui-shu-mo-shu

和 SAM 一样强大

---

**反串的 SAM 的 parent 树就是后缀树**
parent 树有一个性质,父亲是孩子的最长后缀( edp 不同)。
而把串翻转过来之后,反串的 parent 树就满足 : 父亲是孩子的最长前缀 ( beginpos 不同 )。

---

https://37zigen.com/suffix-tree/

---

**处理两个或者多个字符非常有用**

---

https://etaoinwu.com/blog/%E6%84%9F%E6%80%A7%E7%90%86%E8%A7%A3-sam/
压缩后缀树：

1. 第一种压缩方式：压缩单链，形成后缀树：
   由于后缀 trie 每条边上的这个字符串都是原串的一个子串，我们可以只存储这个子串的左右位置，因此我们得到了一个空间复杂度正确的结构
   ![压缩后缀Trie1](%E5%8E%8B%E7%BC%A9%E5%90%8E%E7%BC%80Trie1.png)
2. 第二种压缩方式：合并 endPos 等价类，形成后缀自动机：
   ![压缩后缀Trie2](%E5%8E%8B%E7%BC%A9%E5%90%8E%E7%BC%80Trie2.png)

后缀自动机的 fail 树是反串后缀树
