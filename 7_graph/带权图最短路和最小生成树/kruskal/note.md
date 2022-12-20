1. 准备好所有边 [u,v,weight] 按 wweight 升序排列
2. 并查集构造生成树

3. 给定无向图的边，求出一个最小生成树(**如果不存在,则求出的是森林中的多个最小生成树**)
   [1697. 检查边长度限制的路径是否存在-在线](1697.%20%E6%A3%80%E6%9F%A5%E8%BE%B9%E9%95%BF%E5%BA%A6%E9%99%90%E5%88%B6%E7%9A%84%E8%B7%AF%E5%BE%84%E6%98%AF%E5%90%A6%E5%AD%98%E5%9C%A8-%E5%9C%A8%E7%BA%BF.go)

拓展

1. 包含某条边的最小生成树
   https://blog.csdn.net/weixin_43261862/article/details/104000715
   第一种情况：先对整图跑一遍 MST，边跑边记录 MST 的边，如果求的边正好在 MST 中，那直接输出 W（MST 的权值）即可。
   第二种情况：如果不在 MST 中的话，那我们就把这条边添加进 MST 中，那么一定会成环，接下来我们只需要`把这个环中除去刚添加进的那条边之外的最大权值的那条边去掉`，那么就会出现一棵新的包含该边的生成树。(LCA 倍增求 **树节点两点路径上最大边权**)
2. 删除某条边的最小生成树(可并堆 Meldable Heap)
   https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2559
   https://blog.hamayanhamayan.com/entry/2017/06/07/234608
