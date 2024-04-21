# #include <cassert>
# #include <iostream>
# #include <vector>

# template<typename T>
# class Imos {
# public:
#     const int n;
#     std::vector<T> line;

#     Imos(int n) : n(n) {
#         this->line.resize(n + 1, 0);
#     }

#     // [left, right) に +x
#     void add(const int left, const int right, const T x) {
#         if (left == right) {
#             return;
#         }
#         assert(left < right);
#         this->line[left] += x;
#         this->line[right] -= x;
#     }

#     // [left, right) += x
#     // 加算位置が n 以上の場合は 0 に戻って加算される
#     // O(log n)
#     void add_circle(long long left, long long right, const T x) {
#         assert(left < right);

#         const long long num_loop = (right - left) / this->n;
#         this->add(0, this->n, x * num_loop);

#         // ループで終わり
#         if ((right - left) % this->n == 0) {
#             return;
#         }

#         left %= this->n;
#         right %= this->n;

#         if (left < right) {
#             this->add(left, right, x);
#         } else {
#             this->add(left, this->n, x);
#             this->add(0, right, x);
#         }
#     }

#     void build() {
#         for (int i = 1; i < (int) line.size(); ++i) {
#             this->line[i] += this->line[i - 1];
#         }
#     }

#     T access(const int i) const {
#         return this->line[i];
#     }

#     void dump() const {
#         for (int i = 0; i < this->n; ++i) {
#             if (i != 0) {
#                 std::cout << " ";
#             }
#             std::cout << this->access(i);
#         }
#         std::cout << std::endl;
#     }
# };
# CircularDiff-环形差分
