https://nyaannyaan.github.io/library/dp/monge-d-edge-shortest-path-enumerate.hpp
https://maspypy.github.io/library/convex/monge.hpp
https://noshi91.github.io/algorithm-encyclopedia/d-edge-shortest-path-monge
https://topcoder-g-hatena-ne-jp.jag-icpc.org/spaghetti_source/20120915/1347668163.html

---

https://www.cnblogs.com/alex-wei/p/DP_Involution.html
王钦石二分，简称 wqs 二分，又称带权二分，斜率凸优化。
该算法常见于 限制选取物品个数 的 DP。它有很明显的标志，因此看起来比较套路。
设 f(i)表示**最多（或恰好，或至少）选取 i 个物品**时的答案，若为凸函数则适用 wqs 二分。通常的凸性难以证明，考场上可以猜结论：nk 做不了就是凸的。
