// !线性时间RMQ
// !https://ei1333.github.io/library/test/verify/yosupo-staticrmq-4.test.cpp
// https://qiita.com/okateim/items/e2f4a734db4e5f90e410#m-%E4%B8%8A%E3%81%AE-rmq
// https://qiita.com/okateim/items/a1b0523c1c862009e81f#4%E4%BA%BA%E3%81%AE%E3%83%AD%E3%82%B7%E3%82%A2%E4%BA%BA%E3%81%AE%E6%96%B9%E6%B3%95%E3%82%92%E7%94%A8%E3%81%84%E3%81%9F%E9%AB%98%E9%80%9F%E5%8C%96
// https://hotman78.github.io/cpplib/data_structure/RMQ.hpp
// https://noshi91.hatenablog.com/entry/2018/08/16/125415
// https://etaoinwu.com/blog/rmq-in-linear-time/
// https://etaoinwu.com/blog/level-ancestor-in-linear-time/
// https://etaoinwu.com/blog/lca-in-linear-time/
// https://www.cnblogs.com/whx1003/p/13996517.html
//
// !实际查询效率与只使用ST表差不多。
// 这几个东西的关系:
//
// !线性时间RMQ:
// 1.利用笛卡尔树将序列的RMQ转化为笛卡尔树上的LCA(`线性时间LCA`).
// !线性时间LCA:
// 1.利用树的欧拉序将LCA转化为±1RMQ;
// 2.±1RMQ(PM1RMQ)利用四毛子算法(Method of Four Russians/4人のロシア人の方法)将序列分块，块大小为`ceil(logn/2)`;
//   这里四毛子算法指的是一种思想,即`预处理小问题的答案并将其存储在表中，然后在解决大问题时通过查表来加速.`
// TODO
package main

func main() {

}
