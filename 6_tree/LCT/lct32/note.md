# Euler Tour Tree

Euler Tour Tree 也是动态树，支持 O(log2n)的换根，加边，删边操作。类似于 Link-Cut-Tree，但是底层原理是将树表示为 DFS 序，并将 DFS 序放到平衡树上去维护。
要了解更多，可以去阅读斯坦福大学的课件，地址如下：

http://web.stanford.edu/class/archive/cs/cs166/cs166.1146/lectures/05/Slides05.pdf

Euler Tour Tree 和 Link Cut Tree 的区别在于，后者是以 Splay 森林保存信息，每个森林中的树都代表一条树上的路径，因此 Link Cut Tree 非常适合用于处理路径问题。而前者是以一株平衡树进行维护的，因此可以很好的利用树的一些性质，比如处理子树问题。
