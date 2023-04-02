https://ei1333.hateblo.jp/entry/2018/05/29/011140
https://nyaannyaan.github.io/library/lct/link-cut-base.hpp
https://wenku.baidu.com/view/7857b870aaea998fcc220ed8.html

动态树：LCT/Euler Tour Tree
支持 link 连接边/cut 断开边的操作，结合 splay 维护链

- LCT 一般维护路径点权
  如果要维护路径边权，可以在每两个点之间插入一个辅助顶点，将路径权值放到这个顶点上，查询时直接查询所有**辅助顶点间权值和**即可

- 动态维护连通性&双联通分量（可以说是并查集的升级，因为并查集只能连不能断）

---

有的时候删除结点不必真的删除，而是在`线段树上标记删除`
