// https://github.com/MitI-7/WaveletMatrix/blob/master/WaveletMatrix/WaveletMatrix.hpp
package main

import (
	"fmt"
	"math/bits"
	"math/rand"
)

// class WaveletMatrix {
// 	private:
// 			std::vector<SuccinctBitVector> bit_arrays;
// 			std::vector<uint64_t> begin_one;                    // 各bitに着目したときの1の開始位置
// 			std::map<uint64_t, uint64_t> begin_alphabet;        // 最後のソートされた配列で各文字の開始位置
// 			std::vector<std::vector<uint64_t>> cumulative_sum;  // 各bitに着目したときの累積和

// 			uint64_t size;                                 // 与えられた配列のサイズ
// 			uint64_t maximum_element;                      // 文字数
// 			uint64_t bit_size;                             // 文字を表すのに必要なbit数

// 	public:
// 			WaveletMatrix (const std::vector<uint64_t> &array) {
// 					assert(array.size() > 0);
// 					size = array.size();
// 					maximum_element =  *max_element(array.begin(), array.end()) + 1;
// 					bit_size = get_num_of_bit(maximum_element);
// 					if (bit_size == 0) {
// 							bit_size = 1;
// 					}

// 					for (uint64_t i = 0; i < bit_size; ++i) {
// 							SuccinctBitVector sv(size);
// 							bit_arrays.push_back(sv);
// 					}
// 					this->begin_one.resize(bit_size);
// 					this->cumulative_sum.resize(bit_size + 1, std::vector<uint64_t>(size + 1, 0));

// 					for (uint64_t j = 0; j < array.size(); ++j) {
// 							this->cumulative_sum.at(0).at(j + 1) = this->cumulative_sum.at(0).at(j) + array[j];
// 					}

// 					std::vector<uint64_t> v(array);
// 					for (uint64_t i = 0; i < bit_size; ++i) {

// 							std::vector<uint64_t> temp;
// 							// 0をtempにいれてく
// 							for (uint64_t j = 0; j < v.size(); ++j) {
// 									uint64_t c = v.at(j);
// 									uint64_t bit = (c >> (bit_size - i - 1)) & 1;  //　上からi番目のbit
// 									if (bit == 0) {
// 											temp.push_back(c);
// 											bit_arrays.at(i).setBit(0, j);
// 									}
// 							}

// 							this->begin_one.at(i) = temp.size();

// 							// 1をtempにいれてく
// 							for (uint64_t j = 0; j < v.size(); ++j) {
// 									uint64_t c = v.at(j);
// 									uint64_t bit = (c >> (bit_size - i - 1)) & 1;  //　上からi番目のbit
// 									if (bit == 1) {
// 											temp.push_back(c);
// 											bit_arrays.at(i).setBit(1, j);
// 									}
// 							}

// 							for (uint64_t j = 0; j < temp.size(); ++j) {
// 									this->cumulative_sum.at(i + 1).at(j + 1) = this->cumulative_sum.at(i + 1).at(j) + temp.at(j);
// 							}

// 							bit_arrays.at(i).build();
// 							v = temp;
// 					}

// 					// ソートされた配列内での各文字の位置を取得
// 					for (int i = v.size() - 1; i >= 0; --i) {
// 							this->begin_alphabet[v.at(i)] = i;
// 					}
// 			}

// 			// v[pos]
// 			uint64_t access(uint64_t pos) {
// 					if (pos >= this->size) { return NOTFOUND; }

// 					uint64_t c = 0;
// 					for (uint64_t i = 0; i < bit_arrays.size(); ++i) {
// 							uint64_t bit = bit_arrays.at(i).access(pos);   // もとの数値のi番目のbit
// 							c = (c <<= 1) | bit;
// 							pos = bit_arrays.at(i).rank(bit, pos);
// 							if (bit) {
// 									pos += this->begin_one.at(i);
// 							}
// 					}
// 					return c;
// 			}

// 			// i番目のcの位置 + 1を返す。rankは1-origin
// 			uint64_t select(uint64_t c, uint64_t rank) {
// 					assert(rank > 0);
// 					if (c >= maximum_element) {
// 							return NOTFOUND;
// 					}
// 					if (this->begin_alphabet.find(c) == this->begin_alphabet.end()) {
// 							return NOTFOUND;
// 					}

// 					uint64_t index = this->begin_alphabet.at(c) + rank;
// 					for (uint64_t i = 0; i < bit_arrays.size(); ++i){
// 							uint64_t bit = ((c >> i) & 1);      // 下からi番目のbit
// 							if (bit == 1) {
// 									index -= this->begin_one.at(bit_size - i - 1);
// 							}
// 							index = this->bit_arrays.at(bit_size - i - 1).select(bit, index);
// 					}
// 					return index;
// 			}

// 			// v[begin_pos, end_pos)で最大値のindexを返す
// 			uint64_t maxRange(uint64_t begin_pos, uint64_t end_pos) {
// 					return quantileRange(begin_pos, end_pos, end_pos - begin_pos - 1);
// 			}

// 			// v[begin_pos, end_pos)で最小値のindexを返す
// 			uint64_t minRange(uint64_t begin_pos, uint64_t end_pos) {
// 					return quantileRange(begin_pos, end_pos, 0);
// 			}

// 			// v[begin_pos, end_pos)でk番目に小さい数値のindexを返す(kは0-origin)
// 			// つまり小さい順に並べてk番目の値
// 			uint64_t quantileRange(uint64_t begin_pos, uint64_t end_pos, uint64_t k) {
// 					if ((end_pos > size || begin_pos >= end_pos) || (k >= end_pos - begin_pos)) {
// 							return NOTFOUND;
// 					}

// 					uint64_t val = 0;
// 					for (uint64_t i = 0; i < bit_size; ++i) {
// 							const uint64_t num_of_zero_begin = bit_arrays.at(i).rank(0, begin_pos);
// 							const uint64_t num_of_zero_end = bit_arrays.at(i).rank(0, end_pos);
// 							const uint64_t num_of_zero = num_of_zero_end - num_of_zero_begin;     // beginからendまでにある0の数
// 							const uint64_t bit = (k < num_of_zero) ? 0 : 1;                       // k番目の値の上からi番目のbitが0か1か

// 							if (bit) {
// 									k -= num_of_zero;
// 									begin_pos = this->begin_one.at(i) + begin_pos - num_of_zero_begin;
// 									end_pos = this->begin_one.at(i) + end_pos - num_of_zero_end;
// 							}
// 							else {
// 									begin_pos = num_of_zero_begin;
// 									end_pos = num_of_zero_begin + num_of_zero;
// 							}

// 							val = ((val << 1) | bit);
// 					}

// 					uint64_t left = 0;
// 					for (uint64_t i = 0; i < bit_size; ++i) {
// 							const uint64_t bit = (val >> (bit_size - i - 1)) & 1;  // 上からi番目のbit
// 							left = bit_arrays.at(i).rank(bit, left);               // cのi番目のbitと同じ数値の数
// 							if (bit) {
// 									left += this->begin_one.at(i);
// 							}
// 					}

// 					const uint64_t rank = begin_pos + k - left + 1;
// 					return select(val, rank) - 1;
// 			}

// 			// v[0, pos)のcの数
// 			uint64_t rank(uint64_t c, uint64_t pos) {
// 					assert(pos < size);
// 					if (c >= maximum_element) {
// 							return 0;
// 					}
// 					if (this->begin_alphabet.find(c) == this->begin_alphabet.end()) {
// 							return 0;
// 					}

// 					for (uint64_t i = 0; i < bit_size; ++i) {
// 							uint64_t bit = (c >> (bit_size - i - 1)) & 1;  // 上からi番目のbit
// 							pos = bit_arrays.at(i).rank(bit, pos);         // cのi番目のbitと同じ数値の数
// 							if (bit) {
// 									pos += this->begin_one.at(i);
// 							}
// 					}

// 					uint64_t begin_pos = this->begin_alphabet.at(c);
// 					return pos - begin_pos;
// 			}

// 			// v[begin_pos, end_pos)で[min, max)に入る値の個数
// 			uint64_t rangeFreq(uint64_t begin_pos, uint64_t end_pos, uint64_t min_c, uint64_t max_c) {
// 					if ((end_pos > size || begin_pos >= end_pos) || (min_c >= max_c) || min_c >= maximum_element) {
// 							return 0;
// 					}

// 					const auto maxi_t = rankAll(max_c, begin_pos, end_pos);
// 					const auto mini_t = rankAll(min_c, begin_pos, end_pos);
// 					return std::get<1>(maxi_t) - std::get<1>(mini_t);
// 			}

// 			// v[0, pos)でcより小さい文字の数
// 			uint64_t rankLessThan(uint64_t c, uint64_t begin, uint64_t end) {
// 					auto t = rankAll(c, begin, end);
// 					return std::get<1>(t);
// 			}

// 			// v[0, pos)でcより大きい文字の数
// 			uint64_t rankMoreThan(uint64_t c, uint64_t begin, uint64_t end) {
// 					auto t = rankAll(c, begin, end);
// 					return std::get<2>(t);
// 			}

// 			// v[begin, end)で(cと同じ値の数、cより小さい値の数、cより大きい値の数)を求める
// 			std::tuple<uint64_t, uint64_t, uint64_t> rankAll(const uint64_t c, uint64_t begin, uint64_t end) {
// 					assert(end <= size);
// 					const uint64_t num = end - begin;

// 					if (begin >= end) {
// 							return std::make_tuple(0, 0, 0);
// 					}
// 					if (c >= maximum_element || end == 0) {
// 							return std::make_tuple(0, num, 0);
// 					}

// 					uint64_t rank_less_than = 0, rank_more_than = 0;
// 					for (size_t i = 0; i < bit_size && begin < end; ++i) {
// 							const uint64_t bit = (c >> (bit_size - i - 1)) & 1;

// 							const uint64_t rank0_begin = this->bit_arrays.at(i).rank(0, begin);
// 							const uint64_t rank0_end = this->bit_arrays.at(i).rank(0, end);
// 							const uint64_t rank1_begin = begin - rank0_begin;
// 							const uint64_t rank1_end = end - rank0_end;

// 							if (bit) {
// 									rank_less_than += (rank0_end - rank0_begin);    // i番目のbitが0のものは除外される
// 									begin = this->begin_one.at(i) + rank1_begin;
// 									end = this->begin_one.at(i) + rank1_end;
// 							} else {
// 									rank_more_than += (rank1_end - rank1_begin);    // i番目のbitが1のものは除外される
// 									begin = rank0_begin;
// 									end = rank0_end;
// 							}
// 					}

// 					const uint64_t rank = num - rank_less_than - rank_more_than;
// 					return std::make_tuple(rank, rank_less_than, rank_more_than);
// 			}

// 			// T[s, e)で出現回数が多い順にk個の(値，頻度)を返す
// 			// 頻度が同じ場合は値が小さいものが優先される
// 			std::vector<std::pair<uint64_t, uint64_t>> topk(uint64_t s, uint64_t e, uint64_t k) {
// 					assert(s < e);
// 					std::vector<std::pair<uint64_t, uint64_t>> result;

// 					// (頻度，深さ，値)の順でソート
// 					auto c = [](const std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, uint64_t> &l, const std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, uint64_t> &r) {
// 							// width
// 							if (std::get<0>(l) != std::get<0>(r)) {
// 									return std::get<0>(l) < std::get<0>(r);
// 							}
// 							// depth
// 							if (std::get<3>(l) != std::get<3>(r)) {
// 									return std::get<3>(l) > std::get<3>(r);
// 							}
// 							// value
// 							if (std::get<4>(l) != std::get<4>(r)) {
// 									return std::get<4>(l) > std::get<4>(r);
// 							}
// 							return true;
// 					};

// 					std::priority_queue<std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, uint64_t>, std::vector<std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, uint64_t>>, decltype(c)> que(c);  // width, left, right, depth, value
// 					que.push(std::make_tuple(e - s, s, e, 0, 0));

// 					while (not que.empty()) {
// 							auto element = que.top(); que.pop();
// 							uint64_t width, left, right, depth, value;
// 							std::tie(width, left, right, depth, value) = element;

// 							if (depth >= this->bit_size) {
// 									result.emplace_back(std::make_pair(value, right - left));
// 									if (result.size() >= k) {
// 											break;
// 									}
// 									continue;
// 							}

// 							// 0
// 							const uint64_t left0 = this->bit_arrays.at(depth).rank(0, left);
// 							const uint64_t right0 = this->bit_arrays.at(depth).rank(0, right);
// 							if (left0 < right0) {
// 									que.push(std::make_tuple(right0 - left0, left0, right0, depth + 1, value));
// 							}

// 							// 1
// 							const uint64_t left1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, left);
// 							const uint64_t right1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, right);
// 							if (left1 < right1) {
// 									que.push(std::make_tuple(right1 - left1, left1, right1, depth + 1, value | (1 << (bit_size - depth - 1))));
// 							}
// 					}

// 					return result;
// 			};

// 			// T[begin_pos, end_pos)でx <= c < yを満たすcの和を返す
// 			uint64_t rangeSum(const uint64_t begin, const uint64_t end, const uint64_t x, const uint64_t y) {
// 					return rangeSum(begin, end, 0, 0, x, y);
// 			}

// 			// T[begin_pos, end_pos)でx <= c < yを満たす最大のcを返す
// 			uint64_t prevValue(const uint64_t begin_pos, const uint64_t end_pos, const uint64_t x, uint64_t y) {
// 					assert(end_pos <= size);
// 					const uint64_t num = end_pos - begin_pos;

// 					if (x >= y or y == 0) {
// 							return NOTFOUND;
// 					}
// 					if (y > maximum_element) {
// 							y = maximum_element;
// 					}

// 					if (begin_pos >= end_pos) {
// 							return NOTFOUND;
// 					}
// 					if (x >= maximum_element || end_pos == 0) {
// 							return NOTFOUND;
// 					}

// 					y--; // x <= c <= yにする

// 					std::stack<std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, bool>> s;   // (begin, end, depth, c, tight)
// 					s.emplace(std::make_tuple(begin_pos, end_pos, 0, 0, true));

// 					while (not s.empty()) {
// 							uint64_t b, e, depth, c;
// 							bool tight;
// 							std::tie(b, e, depth, c, tight) = s.top(); s.pop();

// 							if (depth == bit_size) {
// 									if (c >= x) {
// 											return c;
// 									}
// 									continue;
// 							}

// 							const uint64_t bit = (y >> (bit_size - depth - 1)) & 1;

// 							const uint64_t rank0_begin = this->bit_arrays.at(depth).rank(0, b);
// 							const uint64_t rank0_end = this->bit_arrays.at(depth).rank(0, e);
// 							const uint64_t rank1_begin = b - rank0_begin;
// 							const uint64_t rank1_end = e - rank0_end;

// 							// d番目のbitが0のものを使う
// 							const uint64_t b0 = rank0_begin;
// 							const uint64_t e0 = rank0_end;
// 							if (b0 != e0) { // 範囲がつぶれてない
// 									const uint64_t c0 = ((c << 1) | 0);
// 									s.emplace(std::make_tuple(b0, e0, depth + 1, c0, tight and bit == 0));
// 							}

// 							// d番目のbitが1のものを使う
// 							const uint64_t b1 = this->begin_one.at(depth) + rank1_begin;
// 							const uint64_t e1 = this->begin_one.at(depth) + rank1_end;
// 							if (b1 != e1) {
// 									if (not tight or bit == 1) {
// 											const auto c1 = ((c << 1) | 1);
// 											s.emplace(std::make_tuple(b1, e1, depth + 1, c1, tight));
// 									}
// 							}
// 					}

// 					return NOTFOUND;
// 			}

// 			// T[begin_pos, end_pos)でx <= c < yを満たす最小のcを返す
// 			uint64_t nextValue(const uint64_t begin_pos, const uint64_t end_pos, const uint64_t x, const uint64_t y) {
// 					assert(end_pos <= size);
// 					const uint64_t num = end_pos - begin_pos;

// 					if (x >= y or y == 0) {
// 							return NOTFOUND;
// 					}

// 					if (begin_pos >= end_pos) {
// 							return NOTFOUND;
// 					}
// 					if (x >= maximum_element || end_pos == 0) {
// 							return NOTFOUND;
// 					}

// 					std::stack<std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, bool>> s;   // (begin, end, depth, c, tight)
// 					s.emplace(std::make_tuple(begin_pos, end_pos, 0, 0, true));

// 					while (not s.empty()) {
// 							uint64_t b, e, depth, c;
// 							bool tight;
// 							std::tie(b, e, depth, c, tight) = s.top(); s.pop();

// 							if (depth == bit_size) {
// 									if (c < y) {
// 											return c;
// 									}
// 									continue;
// 							}

// 							const uint64_t bit = (x >> (bit_size - depth - 1)) & 1;

// 							const uint64_t rank0_begin = this->bit_arrays.at(depth).rank(0, b);
// 							const uint64_t rank0_end = this->bit_arrays.at(depth).rank(0, e);
// 							const uint64_t rank1_begin = b - rank0_begin;
// 							const uint64_t rank1_end = e - rank0_end;

// 							// d番目のbitが1のものを使う
// 							const uint64_t b1 = this->begin_one.at(depth) + rank1_begin;
// 							const uint64_t e1 = this->begin_one.at(depth) + rank1_end;
// 							if (b1 != e1) {
// 									const auto c1 = ((c << 1) | 1);
// 									s.emplace(std::make_tuple(b1, e1, depth + 1, c1, tight and bit == 1));
// 							}

// 							// d番目のbitが0のものを使う
// 							const uint64_t b0 = rank0_begin;
// 							const uint64_t e0 = rank0_end;
// 							if (b0 != e0) {
// 									if (not tight or bit == 0) {
// 											const uint64_t c0 = ((c << 1) | 0);
// 											s.emplace(std::make_tuple(b0, e0, depth + 1, c0, tight));
// 									}
// 							}
// 					}

// 					return NOTFOUND;
// 			}

// 			// T[s1, e1)とT[s2, e2)に共通して出現する要素を求める
// 			std::vector<std::tuple<uint64_t, uint64_t, uint64_t>> intersect(uint64_t _s1, uint64_t _e1, uint64_t _s2, uint64_t _e2) {
// 					assert(_s1 < _e1);
// 					assert(_s2 < _e2);

// 					std::vector<std::tuple<uint64_t, uint64_t, uint64_t>> intersection;

// 					std::queue<std::tuple<uint64_t, uint64_t, uint64_t, uint64_t, uint64_t, uint64_t>> que; // s1, e1, s2, e2, depth, value
// 					que.push(std::make_tuple(_s1, _e1, _s2, _e2, 0, 0));
// 					while (not que.empty()) {
// 							auto e = que.front(); que.pop();
// 							uint64_t s1, e1, s2, e2, depth, value;
// 							std::tie(s1, e1, s2, e2, depth, value) = e;

// 							if (depth >= this->bit_size) {
// 									intersection.emplace_back(std::make_tuple(value, e1 - s1, e2 - s2));
// 									continue;
// 							}

// 							// 0
// 							uint64_t s1_0 = this->bit_arrays.at(depth).rank(0, s1);
// 							uint64_t e1_0 = this->bit_arrays.at(depth).rank(0, e1);
// 							uint64_t s2_0 = this->bit_arrays.at(depth).rank(0, s2);
// 							uint64_t e2_0 = this->bit_arrays.at(depth).rank(0, e2);

// 							if (s1_0 != e1_0 and s2_0 != e2_0) {
// 									que.push(std::make_tuple(s1_0, e1_0, s2_0, e2_0, depth + 1, value));
// 							}

// 							// 1
// 							uint64_t s1_1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, s1);
// 							uint64_t e1_1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, e1);
// 							uint64_t s2_1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, s2);
// 							uint64_t e2_1 = this->begin_one.at(depth) + this->bit_arrays.at(depth).rank(1, e2);

// 							if (s1_1 != e1_1 and s2_1 != e2_1) {
// 									que.push(std::make_tuple(s1_1, e1_1, s2_1, e2_1, depth + 1, value | (1 << bit_size - depth - 1)));
// 							}
// 					}

// 					return intersection;
// 			};

// 	private:
// 			uint64_t get_num_of_bit(uint64_t x) {
// 					if (x == 0) return 0;
// 					x--;
// 					uint64_t bit_num = 0;
// 					while (x >> bit_num) {
// 							++bit_num;
// 					}
// 					return bit_num;
// 			}

// 			uint64_t rangeSum(const uint64_t begin, const uint64_t end, const uint64_t depth, const uint64_t c, const uint64_t x, const uint64_t y) {
// 					if (begin == end) {
// 							return 0;
// 					}

// 					if (depth == bit_size) {
// 							if (x <= c and c < y) {
// 									return c * (end - begin);   // 値 * 頻度
// 							}
// 							return 0;
// 					}

// 					const uint64_t next_c = ((uint64_t)1 << (bit_size - depth - 1)) | c;                   // 上からdepth番目のbitを立てる
// 					const uint64_t all_one_c = (((uint64_t)1 << (bit_size - depth - 1)) - 1) | next_c;     // depth以降のbitをたてる(これ以降全部1を選んだときの値)
// 					if(all_one_c < x or y <= c) {
// 							return 0;
// 					}

// 					// [begin, pos)のすべての要素は[x, y)
// 					if (x <= c and all_one_c < y) {
// 							return this->cumulative_sum.at(depth).at(end) - this->cumulative_sum.at(depth).at(begin);
// 					}

// 					const uint64_t rank0_begin = this->bit_arrays.at(depth).rank(0, begin);
// 					const uint64_t rank0_end = this->bit_arrays.at(depth).rank(0, end);
// 					const uint64_t rank1_begin = begin - rank0_begin;
// 					const uint64_t rank1_end = end - rank0_end;

//					return rangeSum(rank0_begin, rank0_end, depth + 1, c, x, y) +
//								 rangeSum(this->begin_one.at(depth) + rank1_begin, this->begin_one[depth] + rank1_end, depth + 1, next_c, x, y);
//			}
//	};
func main() {
	nums := []uint64{1, 2, 1, 4, 5, 1, 7, 8, 2, 10}
	wm := NewWaveletMatrixStaticOmni(uint64(len(nums)), func(i uint64) uint64 { return nums[i] })
	_ = wm
	for i := uint64(0); i < uint64(len(nums)); i++ {
		fmt.Println(wm.Get(i))
	}
	fmt.Println("------")
	// fmt.Println(wm.Kth(0, 1))
	// fmt.Println(wm.Kth(1, 1))
	// fmt.Println(wm.Kth(2, 1))
	// fmt.Println(wm.Kth(3, 1))
	fmt.Println(wm.Kth(4, 1), "wm.Kth(4, 1)")
	// fmt.Println(wm.Kth(0, 2))
	// fmt.Println(wm.Kth(1, 2))
	// fmt.Println(wm.Kth(2, 2))

	// fmt.Println("------")
	// fmt.Println(wm.KthSmallest(0, 10, 0))
	// fmt.Println(wm.KthSmallest(0, 10, 1))
	// fmt.Println(wm.KthSmallest(0, 10, 2))
	// fmt.Println(wm.KthSmallest(0, 10, 3))

	checkKth := func(k, v uint64) uint64 {
		cnt := uint64(0)
		for i := uint64(0); i < uint64(len(nums)); i++ {
			if nums[i] == v {
				if cnt == k {
					return i
				}
				cnt++
			}
		}
		return NOT_FOUND
	}
	_ = checkKth

	checkCountPrefix := func(end, v uint64) uint64 {
		cnt := uint64(0)
		for i := uint64(0); i < end; i++ {
			if nums[i] == v {
				cnt++
			}
		}
		return cnt
	}
	_ = checkCountPrefix

	for i := 0; i < 100; i++ {
		k := uint64(rand.Intn(10))
		v := uint64(rand.Intn(10))
		if wm.Kth(k, v) != checkKth(k, v) {
			fmt.Println(wm.Kth(k, v), checkKth(k, v), k, v)
			panic("error")
		}
	}

	for i := 0; i < 100_000; i++ {
		end := uint64(rand.Intn(len(nums) + 1))
		v := uint64(rand.Intn(10))
		if wm.CountPrefix(end, v) != checkCountPrefix(end, v) {
			fmt.Println(wm.CountPrefix(end, v), checkCountPrefix(end, v), end, v)
			panic("error")
		}
	}

	fmt.Println("ok")
	fmt.Println(wm.Kth(4, 1))
	fmt.Println(wm.CountPrefix(10, 2))
}

type WaveletMatrixStaticOmni struct {
	bvs           []*succinctBitVector
	beginOne      []uint64
	beginAlphabet map[uint64]uint64
	presSum       [][]uint64
	size          uint64
	maxElement    uint64
	bitSize       uint64
}

//

func NewWaveletMatrixStaticOmni(n uint64, f func(uint64) uint64) *WaveletMatrixStaticOmni {
	if n <= 0 {
		panic("n must be positive")
	}
	data := make([]uint64, n)
	for i := uint64(0); i < n; i++ {
		data[i] = f(i)
	}
	maxElement := uint64(0)
	for i := uint64(0); i < n; i++ {
		if data[i] > maxElement {
			maxElement = data[i]
		}
	}
	maxElement++
	bitSize := uint64(bits.Len64(maxElement))
	if bitSize == 0 {
		bitSize = 1
	}
	bvs := make([]*succinctBitVector, bitSize)
	for i := uint64(0); i < bitSize; i++ {
		bvs[i] = newSuccinctBitVector(n)
	}
	beginOne := make([]uint64, bitSize)
	beginAlphabet := make(map[uint64]uint64, n)
	presSum := make([][]uint64, bitSize+1)
	for i := uint64(0); i <= bitSize; i++ {
		presSum[i] = make([]uint64, n+1)
	}
	for j := uint64(0); j < n; j++ {
		presSum[0][j+1] = presSum[0][j] + data[j]
	}

	zero, one := make([]uint64, n), make([]uint64, n)
	for i := uint64(0); i < bitSize; i++ {
		bv := bvs[i]
		p, q := uint64(0), uint64(0)
		for j := uint64(0); j < n; j++ {
			c := data[j]
			bit := (c >> (bitSize - i - 1)) & 1
			if bit == 0 {
				zero[p] = c
				p++
			} else {
				one[q] = c
				q++
				bv.Set(j, true)
			}
		}
		beginOne[i] = p
		ps := presSum[i+1]
		for j := uint64(0); j < n; j++ {
			ps[j+1] = ps[j] + data[j]
		}
		bv.Build()
		data, zero = zero, data
		copy(data[p:], one[:q])
	}

	for i := int32(n - 1); i >= 0; i-- {
		beginAlphabet[data[i]] = uint64(i)
	}

	return &WaveletMatrixStaticOmni{
		bvs:           bvs,
		beginOne:      beginOne,
		beginAlphabet: beginAlphabet,
		presSum:       presSum,
		size:          n,
		maxElement:    maxElement,
		bitSize:       bitSize,
	}
}

func (wm *WaveletMatrixStaticOmni) Get(pos uint64) uint64 {
	if pos >= wm.size {
		return NOT_FOUND
	}
	c := uint64(0)
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := wm.bvs[i].Get(pos)
		c <<= 1
		if bit {
			c |= 1
		}
		pos = wm.bvs[i].Count(pos, bit)
		if bit {
			pos += wm.beginOne[i]
		}
	}
	return c
}

func (wm *WaveletMatrixStaticOmni) CountPrefix(end uint64, v uint64) uint64 {
	if end <= 0 {
		return 0
	}
	if end > wm.size {
		end = wm.size
	}
	if v >= wm.maxElement {
		return 0
	}
	beginPos, ok := wm.beginAlphabet[v]
	if !ok {
		return 0
	}
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := v >> (wm.bitSize - i - 1) & 1
		end = wm.bvs[i].Count(end, bit == 1)
		if bit == 1 {
			end += wm.beginOne[i]
		}
	}
	return end - beginPos
}

func (wm *WaveletMatrixStaticOmni) CountRange(start, end uint64, floor, higher uint64) uint64 {
	if end > wm.size || start >= end || floor >= higher || floor >= wm.maxElement {
		return 0
	}
	return wm.CountLessThan(start, end, higher) - wm.CountLessThan(start, end, floor)
}

func (wm *WaveletMatrixStaticOmni) CountLessThan(start, end uint64, v uint64) uint64 {
	_, less, _ := wm.CountAll(start, end, v)
	return less
}

func (wm *WaveletMatrixStaticOmni) CountMoreThan(start, end uint64, v uint64) uint64 {
	_, _, more := wm.CountAll(start, end, v)
	return more
}

// 区间[start, end)中等于v的个数、小于v的个数、大于v的个数.
func (wm *WaveletMatrixStaticOmni) CountAll(start, end uint64, v uint64) (uint64, uint64, uint64) {
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0, 0, 0
	}
	num := end - start
	if v >= wm.maxElement {
		return 0, num, 0
	}
	rankLessThan, rankMoreThan := uint64(0), uint64(0)
	for i := uint64(0); i < wm.bitSize && start < end; i++ {
		bit := v >> (wm.bitSize - i - 1) & 1
		rank0Begin := wm.bvs[i].Count(start, false)
		rank0End := wm.bvs[i].Count(end, false)
		rank1Begin := start - rank0Begin
		rank1End := end - rank0End
		if bit == 1 {
			rankLessThan += rank0End - rank0Begin
			start = wm.beginOne[i] + rank1Begin
			end = wm.beginOne[i] + rank1End
		} else {
			rankMoreThan += rank1End - rank1Begin
			start = rank0Begin
			end = rank0End
		}
	}
	rank := num - rankLessThan - rankMoreThan
	return rank, rankLessThan, rankMoreThan
}

// k: 0-indexed.
func (wm *WaveletMatrixStaticOmni) Kth(k uint64, v uint64) uint64 {
	k++
	if v >= wm.maxElement {
		return NOT_FOUND
	}
	var index uint64
	if tmp, ok := wm.beginAlphabet[v]; !ok {
		return NOT_FOUND
	} else {
		index = tmp + k
	}
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := v >> i & 1
		if bit == 1 {
			index -= wm.beginOne[wm.bitSize-i-1]
		}
		if index == 0 {
			return NOT_FOUND
		}
		index = wm.bvs[wm.bitSize-i-1].Kth(index, bit == 1)
	}
	if index == 0 {
		return NOT_FOUND
	}
	return index - 1
}

func (wm *WaveletMatrixStaticOmni) KthSmallest(start, end uint64, k uint64) uint64 {
	if end > wm.size || start >= end || k >= end-start {
		return NOT_FOUND
	}
	val := uint64(0)
	for i := uint64(0); i < wm.bitSize; i++ {
		numOfZeroBegin := wm.bvs[i].Count(start, false)
		numOfZeroEnd := wm.bvs[i].Count(end, false)
		numOfZero := numOfZeroEnd - numOfZeroBegin
		bit := uint64(0)
		if k >= numOfZero {
			bit = 1
		}
		if bit == 1 {
			k -= numOfZero
			start = wm.beginOne[i] + start - numOfZeroBegin
			end = wm.beginOne[i] + end - numOfZeroEnd
		} else {
			start = numOfZeroBegin
			end = numOfZeroBegin + numOfZero
		}
		val = (val << 1) | bit
	}

	left := uint64(0)
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := (val >> (wm.bitSize - i - 1)) & 1
		left = wm.bvs[i].Count(left, bit == 1)
		if bit == 1 {
			left += wm.beginOne[i]
		}
	}

	rank := start + k - left + 1
	return wm.Kth(val, rank) - 1
}

func (wm *WaveletMatrixStaticOmni) TopK(start, end uint64, k uint64) [][2]uint64 {
	res := [][2]uint64{}
	type item struct {
		width, left, right, depth, value uint64
	}
	// (频率、深度、值)排序
	pq := NewHeap[item](func(a, b item) bool {
		if a.width != b.width {
			return a.width < b.width
		}
		if a.depth != b.depth {
			return a.depth > b.depth
		}
		if a.value != b.value {
			return a.value > b.value
		}
		return true
	}, nil)

	pq.Push(item{end - start, start, end, 0, 0})
	for pq.Len() > 0 {
		element := pq.Pop()
		left, right, depth, value := element.left, element.right, element.depth, element.value
		if depth >= wm.bitSize {
			res = append(res, [2]uint64{value, right - left})
			if uint64(len(res)) >= k {
				break
			}
			continue
		}

		left0 := wm.bvs[depth].Count(left, false)
		right0 := wm.bvs[depth].Count(right, false)
		if left0 < right0 {
			pq.Push(item{right0 - left0, left0, right0, depth + 1, value})
		}

		left1 := wm.beginOne[depth] + wm.bvs[depth].Count(left, true)
		right1 := wm.beginOne[depth] + wm.bvs[depth].Count(right, true)
		if left1 < right1 {
			pq.Push(item{right1 - left1, left1, right1, depth + 1, value | (1 << (wm.bitSize - depth - 1))})
		}
	}

	return res
}

func (wm *WaveletMatrixStaticOmni) RangeSum(start, end uint64, floor, higher uint64) uint64 {
	if end > wm.size || start >= end || floor >= higher || floor >= wm.maxElement {
		return 0
	}
	return wm._rangeSum(start, end, 0, 0, floor, higher)
}

// 区间前驱.
func (wm *WaveletMatrixStaticOmni) PrevValue(start, end, floor, higher uint64) uint64 {
	if floor >= higher || higher == 0 {
		return NOT_FOUND
	}
	if start >= end {
		return NOT_FOUND
	}
	if floor >= wm.maxElement || end == 0 {
		return NOT_FOUND
	}
	higher--
	type item struct {
		start, end, depth, c uint64
		tight                bool
	}
	stack := make([]item, 0)
	stack = append(stack, item{start, end, 0, 0, true})
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		b, e, depth, c, tight := last.start, last.end, last.depth, last.c, last.tight
		if depth == wm.bitSize {
			if c >= floor {
				return c
			}
			continue
		}
		bit := (higher >> (wm.bitSize - depth - 1)) & 1
		rank0Begin := wm.bvs[depth].Count(b, false)
		rank0End := wm.bvs[depth].Count(e, false)
		rank1Begin := b - rank0Begin
		rank1End := e - rank0End

		b0 := rank0Begin
		e0 := rank0End
		if b0 != e0 {
			c0 := (c << 1) | 0
			stack = append(stack, item{b0, e0, depth + 1, c0, tight && bit == 0})
		}

		b1 := wm.beginOne[depth] + rank1Begin
		e1 := wm.beginOne[depth] + rank1End
		if b1 != e1 {
			if !tight || bit == 1 {
				c1 := (c << 1) | 1
				stack = append(stack, item{b1, e1, depth + 1, c1, tight})
			}
		}
	}

	return NOT_FOUND
}

// 区间后继.
func (wm *WaveletMatrixStaticOmni) NextValue(start, end, floor, higher uint64) uint64 {
	if floor >= higher || higher == 0 {
		return NOT_FOUND
	}
	if start >= end {
		return NOT_FOUND
	}
	if floor >= wm.maxElement || end == 0 {
		return NOT_FOUND
	}
	type item struct {
		start, end, depth, c uint64
		tight                bool
	}

	stack := make([]item, 0)
	stack = append(stack, item{start, end, 0, 0, true})
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		b, e, depth, c, tight := last.start, last.end, last.depth, last.c, last.tight
		if depth == wm.bitSize {
			if c < higher {
				return c
			}
			continue
		}
		bit := (floor >> (wm.bitSize - depth - 1)) & 1
		rank0Begin := wm.bvs[depth].Count(b, false)
		rank0End := wm.bvs[depth].Count(e, false)
		rank1Begin := b - rank0Begin
		rank1End := e - rank0End

		b1 := wm.beginOne[depth] + rank1Begin
		e1 := wm.beginOne[depth] + rank1End
		if b1 != e1 {
			c1 := (c << 1) | 1
			stack = append(stack, item{b1, e1, depth + 1, c1, tight && bit == 1})
		}

		b0 := rank0Begin
		e0 := rank0End
		if b0 != e0 {
			if !tight || bit == 0 {
				c0 := (c << 1) | 0
				stack = append(stack, item{b0, e0, depth + 1, c0, tight})
			}
		}
	}

	return NOT_FOUND
}

func (wm *WaveletMatrixStaticOmni) Intersect(s1, e1, s2, e2 uint64) [][3]uint64 {
	if s1 >= e1 || s2 >= e2 {
		return nil
	}
	intersection := make([][3]uint64, 0)
	queue := [][6]uint64{{s1, e1, s2, e2, 0, 0}} // s1, e1, s2, e2, depth, value
	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]
		s1, e1, s2, e2, depth, value := e[0], e[1], e[2], e[3], e[4], e[5]
		if depth == wm.bitSize {
			intersection = append(intersection, [3]uint64{value, e1 - s1, e2 - s2})
			continue
		}
		bit := (value >> (wm.bitSize - depth - 1)) & 1
		rank0Begin := wm.bvs[depth].Count(s1, false)
		rank0End := wm.bvs[depth].Count(e1, false)
		rank1Begin := s1 - rank0Begin
		rank1End := e1 - rank0End
		if bit == 1 {
			queue = append(queue, [6]uint64{rank0Begin, rank0End, rank1Begin, rank1End, depth + 1, value})
		}
		rank0Begin = wm.bvs[depth].Count(s2, false)
		rank0End = wm.bvs[depth].Count(e2, false)
		rank1Begin = s2 - rank0Begin
		rank1End = e2 - rank0End
		if bit == 1 {
			queue = append(queue, [6]uint64{rank0Begin, rank0End, rank1Begin, rank1End, depth + 1, value | (1 << (wm.bitSize - depth - 1))})
		}
	}

	return intersection
}

func (wm *WaveletMatrixStaticOmni) _rangeSum(start, end, depth, c, x, y uint64) uint64 {
	if start == end {
		return 0
	}
	if depth == wm.bitSize {
		if x <= c && c < y {
			return c * (end - start)
		}
		return 0
	}
	nextC := ((1 << (wm.bitSize - depth - 1)) | c)
	allOneC := (((1 << (wm.bitSize - depth - 1)) - 1) | nextC)
	if allOneC < x || y <= c {
		return 0
	}

	if x <= c && allOneC < y {
		return wm.presSum[depth+1][end] - wm.presSum[depth+1][start]
	}

	rank0Begin := wm.bvs[depth].Count(start, false)
	rank0End := wm.bvs[depth].Count(end, false)
	rank1Begin := start - rank0Begin
	rank1End := end - rank0End
	return wm._rangeSum(rank0Begin, rank0End, depth+1, c, x, y) + wm._rangeSum(wm.beginOne[depth]+rank1Begin, wm.beginOne[depth]+rank1End, depth+1, nextC, x, y)
}

const NOT_FOUND = ^uint64(0)

type succinctBitVector struct {
	size  uint64
	ones  uint64
	large []uint64 // 大块
	small []uint16 // 小块
	bv    []uint16 // bitVector
}

func newSuccinctBitVector(n uint64) *succinctBitVector {
	res := &succinctBitVector{size: n}
	res.bv = make([]uint16, (n+15)>>4+1)
	res.large = make([]uint64, n>>9+1)
	res.small = make([]uint16, n>>4+1)
	return res
}

func (sbv *succinctBitVector) Set(pos uint64, bit bool) {
	blockPos := pos >> 4
	offset := pos & 15
	if bit {
		sbv.bv[blockPos] |= 1 << offset
	} else {
		sbv.bv[blockPos] &= ^(1 << offset)
	}
}

func (sbv *succinctBitVector) Get(pos uint64) bool {
	blockPos := pos >> 4
	offset := pos & 15
	return (sbv.bv[blockPos]>>offset)&1 == 1
}

func (sbv *succinctBitVector) Build() {
	num := uint64(0)
	for i := uint64(0); i <= sbv.size; i++ {
		if i&511 == 0 {
			sbv.large[i>>9] = num
		}
		if i&15 == 0 {
			sbv.small[i>>4] = uint16(num - sbv.large[i>>9])
		}
		if i != sbv.size && i&15 == 0 {
			num += uint64(bits.OnesCount16(sbv.bv[i>>4]))
		}
	}
	sbv.ones = num
}

func (sbv *succinctBitVector) Count(end uint64, bit bool) uint64 {
	if bit {
		return sbv.large[end>>9] + uint64(sbv.small[end>>4]) + uint64(bits.OnesCount16(sbv.bv[end>>4]&(1<<(end&15)-1)))
	}
	return end - sbv.Count(end, true)
}

// !kth 从0开始.
func (sbv *succinctBitVector) Kth(k uint64, bit bool) uint64 {

	k++
	if !bit && k > sbv.size-sbv.ones {
		return NOT_FOUND
	}
	if bit && k > sbv.ones {
		return NOT_FOUND
	}

	// 大块内搜索
	largeIndex := uint64(0)
	{
		left, right := uint64(0), uint64(len(sbv.large))
		for right-left > 1 {
			mid := (left + right) >> 1
			var r uint64
			if bit {
				r = sbv.large[mid]
			} else {
				r = mid<<9 - sbv.large[mid]
			}
			if r < k {
				left = mid
				largeIndex = mid
			} else {
				right = mid
			}
		}
	}

	// 小块内搜索
	smallIndex := (largeIndex << 9) >> 4
	{
		left, right := (largeIndex<<9)>>4, min64(((largeIndex+1)<<9)>>4, uint64(len(sbv.small)))
		for right-left > 1 {
			mid := (left + right) >> 1
			r := sbv.large[largeIndex] + uint64(sbv.small[mid])
			if !bit {
				r = mid<<4 - r
			}
			if r < k {
				left = mid
				smallIndex = mid
			} else {
				right = mid
			}
		}
	}

	// bitVector内搜索
	rankPos := uint64(0)
	{
		beginBlockIndex := (smallIndex << 4) >> 4
		totalBit := sbv.large[largeIndex] + uint64(sbv.small[smallIndex])
		if !bit {
			totalBit = smallIndex<<4 - totalBit
		}
		for i := uint64(0); ; i++ {
			b := uint64(bits.OnesCount16(sbv.bv[beginBlockIndex+i]))
			if !bit {
				b = 16 - b
			}
			if totalBit+b >= k {
				block := uint64(sbv.bv[beginBlockIndex+i])
				if !bit {
					block = ^block
				}
				rankPos = (beginBlockIndex+i)<<4 + sbv._selectInBlock(block, k-totalBit)
				break
			}
			totalBit += b
		}
	}

	return rankPos
}

func (sbv *succinctBitVector) _selectInBlock(x uint64, rank uint64) uint64 {
	x1 := x - ((x & 0xAAAAAAAAAAAAAAAA) >> 1)
	x2 := (x1 & 0x3333333333333333) + ((x1 >> 2) & 0x3333333333333333)
	x3 := (x2 + (x2 >> 4)) & 0x0F0F0F0F0F0F0F0F
	pos := uint64(0)
	for ; ; pos += 8 {
		rankNext := (x3 >> pos) & 0xFF
		if rank <= rankNext {
			break
		}
		rank -= rankNext
	}
	v2 := (x2 >> pos) & 0xF
	if rank > v2 {
		rank -= v2
		pos += 4
	}
	v1 := (x1 >> pos) & 0x3
	if rank > v1 {
		rank -= v1
		pos += 2
	}
	v0 := (x >> pos) & 0x1
	if v0 < rank {
		rank -= v0
		pos += 1
	}
	return pos
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
func min64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
