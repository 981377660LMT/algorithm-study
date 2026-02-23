# 吉老师线段树/SegmentTreeBeats

处理区间最值操作 & 区间历史最值

hitonanode 实装：
https://codeforces.com/blog/entry/57319
https://rsm9.hatenablog.com/entry/2021/02/01/220408
https://hitonanode.github.io/cplib-cpp/segmenttree/acl_beats.hpp

ei13333 实装：
https://ei1333.github.io/library/structure/segment-tree/segment-tree-beats.hpp

nyann 实装：
https://nyaannyaan.github.io/library/segment-tree/segment-tree-beats-abstract.hpp

---

mapping 需要多返回一个 bool 值表示是否作用成功

https://sotanishy.github.io/cp-library-cpp/data-structure/segtree/segment_tree_beats.cpp

https://smijake3.hatenablog.com/entry/2019/04/28/021457

https://maspypy.github.io/library/ds/segtree/segtree_beats.hpp

---

https://blog.csdn.net/Maxwell_Newton/article/details/157326723?spm=1001.2014.3001.5502
势能线段树的要点就是：如果一个区间操作无法 O(1)确定对区间的影响，每个元素可进行的操作次数有限，那么我们可以寻找一个判断条件，在线段树递归时，根据这个条件判断区间内是否有需要暴力递归的元素，如果有则暴力递归，否则返回。

**代替：当运算不满足半群时,用 ODT/并查集+线段树暴力遍历区间维护**

---

hitonanode 以前分析过这个 https://rsm9.hatenablog.com/entry/2021/02/01/220408
在 atcoder 的模版上加一行就行 https://hitonanode.github.io/cplib-cpp/segmenttree/acl_beats.hpp

这个可以被 并查集/珂朵莉树 + 线段树代替，容易写一些
