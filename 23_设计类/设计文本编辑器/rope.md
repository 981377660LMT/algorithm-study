rope 是 vim 用的，piece table 是 VSCode 用的

rope 是一种高效的数据结构，用于存储和操作非常长的可变字符串
它减少了应用程序的内存重新分配和数据复制的开销
适合的应用场景:将非常长的字符串上分成多个较小的字符串

关键:**省空间**

The following table outlines the computational complexity of various operations over strings and ropes of length **n**.

| Operation                                   | Regular JavaScript String                     | Rope         |
| ------------------------------------------- | --------------------------------------------- | ------------ |
| Initialization                              | **O(n)**                                      | **O(n)**     |
| Removal of **m** characters                 | **O(n)**                                      | **O(m)**     |
| Insertion of **m** characters               | **O(n)**                                      | **O(m)**     |
| Random access                               | **O(1)**                                      | **O(log n)** |
| Concatenation of a string with length **m** | Best Case **O(1)** / Worst Case **O(n+m)** \* | **O(1)**     |
| Extraction of substring with length **m**   | **O(m)**                                      | **O(m)**     |

\* Most JavaScript engines have certain optimizations in place for concatenating strings.

The Rope data structure really shines in outperforming regular JS strings when **_m <<< n_**. It's therefore best suited for
applications that perform small, but very frequent operations in very large strings. (e.g. text editors)
