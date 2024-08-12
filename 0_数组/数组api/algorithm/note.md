https://learn.microsoft.com/zh-cn/cpp/standard-library/algorithm?view=msvc-170
https://learn.microsoft.com/zh-cn/cpp/standard-library/algorithm-functions?view=msvc-170#adjacent_find
C++ 标准库算法处理迭代器范围通常是由其开始或末尾位置指定。
从 C++20 开始，<algorithm>中定义的大多数算法也以采用 range 的形式提供。 例如，可以调用 `ranges::sort(v1, greater<int>());，而不调用 sort(v1.begin(), v1.end(), greater<int>());`

C++ 标准库算法可以同时处理不同类型的容器对象。 两个后缀已用于传递与算法目的相关的信息：
`_if后缀指`示将算法用于对元素的值（而非元素本身）产生作用的函数对象。 例如，find_if 算法查找其值满足函数对象指定的条件的元素，而 find 算法搜索特定值。
`_copy后缀`指示算法通常修改复制的值，而不是复制修改的值。 换句话说，它们不会修改源范围的元素，而是将结果放入输出范围/迭代器中。 例如，reverse 算法反向排序范围中的元素，而 reverse_copy 算法将反向结果复制到目标范围。
