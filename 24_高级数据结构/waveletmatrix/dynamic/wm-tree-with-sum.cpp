#include <bits/stdc++.h>
#include "../dynamicBitVector/DynamicBitVector.hpp"

class WaveletNode {
public:
    std::weak_ptr<WaveletNode> parent;
    std::shared_ptr<WaveletNode> left;
    std::shared_ptr<WaveletNode> right;
    DynamicBitVector bitVector;

    WaveletNode(std::weak_ptr<WaveletNode> parent) : parent(parent), left(nullptr), right(nullptr) {}
    WaveletNode() : left(nullptr), right(nullptr) {}

};

class DynamicWaveletTree {
public:
    std::shared_ptr<WaveletNode> root;
    std::vector<std::weak_ptr<WaveletNode>> leaves;

    uint64_t size;
    const uint64_t maximum_element; // 最大の数値
    uint64_t bit_size;              // 文字を表すのに必要なbit数

public:
    DynamicWaveletTree(uint64_t maximum_element) : size(0), maximum_element(maximum_element + 1) {
        this->root = std::shared_ptr<WaveletNode>(new WaveletNode);
        this->bit_size = this->get_num_of_bit(maximum_element);
        this->leaves.resize(maximum_element);
    }

    DynamicWaveletTree(uint64_t maximum_element, const std::vector<uint64_t> &array) : size(array.size()), maximum_element(maximum_element + 1) {
        this->root = std::shared_ptr<WaveletNode>(new WaveletNode);
        this->bit_size = this->get_num_of_bit(maximum_element);
        this->leaves.resize(maximum_element);
        if (array.empty()) {
            return ;
        }
        build(this->root, 0, array);
    }

    void build(const std::shared_ptr<WaveletNode> &node, uint64_t depth, const std::vector<uint64_t> &array) {
        assert(not array.empty());

        std::vector<uint64_t> left, right;
        for (uint64_t c : array) {
            uint64_t p = 0;
            if (bit_size > depth) {
                p = bit_size - depth - 1;
            }
            const uint64_t bit = (c >> p) & 1;  // 上からdepth番目のbit
            node->bitVector.push_back(bit);

            if (bit == 0) {
                left.emplace_back(c);
            }
            else {
                right.emplace_back(c);
            }
        }

        // 葉に到達
        if ((int)bit_size == (int)depth) {
            uint64_t c = array.at(0);
            assert((c & 1) == node->bitVector.access(0));
            this->leaves[c] = node;
            return;
        }

        // left node
        if (not left.empty()) {
            if (node->left == nullptr) {
                std::shared_ptr<WaveletNode> left_node(new WaveletNode(node));
                node->left = left_node;
            }
            build(node->left, depth + 1, left);
        }

        // right node
        if (not right.empty()) {
            if (node->right == nullptr) {
                std::shared_ptr<WaveletNode> right_node(new WaveletNode(node));
                node->right = right_node;
            }
            build(node->right, depth + 1, right);
        }
    }

    uint64_t access(uint64_t pos) {
        assert(pos < this->size);

        auto node = this->root;
        uint64_t c = 0;
        for (int i = 0; i < bit_size; ++i) {
            const uint64_t bit = node->bitVector.access(pos);   // T[pos]のi番目のbit
            c = (c <<= 1) | bit;
            pos = node->bitVector.rank(bit, pos);
            if (bit == 0) {
                node = node->left;
            }
            else {
                node = node->right;
            }
        }

        return c;
    }

    // v[0, pos)のcの数
    uint64_t rank(uint64_t c, uint64_t pos) {
        assert(pos <= size);
        if (c >= maximum_element) {
            return 0;
        }

        auto node = this->root;
        for (uint64_t i = 0; i < bit_size; ++i) {
            const uint64_t bit = (c >> (bit_size - i - 1)) & 1;  // 上からi番目のbit
            pos = node->bitVector.rank(bit, pos);             // cのi番目のbitと同じ数値の数
            node = bit == 0 ? node->left : node->right;
        }

        return pos;
    }

    // i番目のcの位置 + 1を返す。rankは1-origin
    uint64_t select(uint64_t c, uint64_t rank) {
        assert(rank > 0);
        if (c >= maximum_element) {
            return NOTFOUND;
        }

        auto node = this->leaves[c].lock()->parent;
        for (int i = 0; i < bit_size; ++i) {
            uint64_t bit = ((c >> i) & 1);      // 下からi番目のbit

            auto n = node.lock();
            rank = n->bitVector.select(bit, rank);
            node = n->parent;
        }

        return rank;
    }

    // posにcを挿入する
    void insert(uint64_t pos, uint64_t c) {
        assert(pos <= this->size);

        auto node = this->root;
        for (uint64_t i = 0; i < bit_size; ++i) {
            const uint64_t bit = (c >> (bit_size - i - 1)) & 1;  //　上からi番目のbit
            node->bitVector.insert(pos, bit);
            pos = node->bitVector.rank(bit, pos);
            if (i == bit_size - 1) {
                break;
            }

            if (bit == 0) {
                if (node->left == nullptr) {
                    std::shared_ptr<WaveletNode> left(new WaveletNode(node));
                    node->left = left;
                }
                node = node->left;
            }
            else {
                if (node->right == nullptr) {
                    std::shared_ptr<WaveletNode> right(new WaveletNode(node));
                    node->right = right;
                }
                node = node->right;
            }
        }

        this->size++;
        this->leaves[c] = node;
    }

    // 末尾にcを追加する
    void push_back(uint64_t c) {
        this->insert(this->size, c);
    }

    // posを削除する
    uint64_t erase(uint64_t pos) {
        assert(pos < this->size);

        auto node = this->root;
        uint64_t c = 0;
        for (uint64_t i = 0; i < bit_size; ++i) {
            uint64_t bit = node->bitVector.access(pos);   // もとの数値のi番目のbit
            c = (c <<= 1) | bit;
            auto next_pos = node->bitVector.rank(bit, pos);
            node->bitVector.erase(pos);
            node = bit == 0 ? node->left : node->right;

            pos = next_pos;
        }

        this->size--;
        return c;
    }

    void update(uint64_t pos, uint64_t c) {
        this->erase(pos);
        this->insert(pos, c);
    }

    // v[begin_pos, end_pos)でk番目に小さい数値を返す(kは0-origin)
    // つまり小さい順に並べてk番目の値
    uint64_t quantileRange(uint64_t begin_pos, uint64_t end_pos, uint64_t k) {
        if ((end_pos > size || begin_pos >= end_pos) || (k >= end_pos - begin_pos)) {
            return NOTFOUND;
        }

        auto node = this->root;
        uint64_t val = 0;
        for (uint64_t i = 0; i < bit_size; ++i) {
            const uint64_t num_of_zero_begin = node->bitVector.rank(0, begin_pos);
            const uint64_t num_of_zero_end = node->bitVector.rank(0, end_pos);
            const uint64_t num_of_zero = num_of_zero_end - num_of_zero_begin;     // beginからendまでにある0の数
            const uint64_t bit = (k < num_of_zero) ? 0 : 1;                       // k番目の値の上からi番目のbitが0か1か

            if (bit == 0) {
                node = node->left;
                begin_pos = num_of_zero_begin;
                end_pos = num_of_zero_begin + num_of_zero;
            }
            else {
                node = node->right;
                k -= num_of_zero;
                begin_pos = begin_pos - num_of_zero_begin;
                end_pos = end_pos - num_of_zero_end;
            }

            val = ((val << 1) | bit);
        }


        node = this->root;
        uint64_t left = 0;
        for (uint64_t i = 0; i < bit_size; ++i) {
            const uint64_t bit = (val >> (bit_size - i - 1)) & 1;  // 上からi番目のbit
            left = node->bitVector.rank(bit, left);                // cのi番目のbitと同じ数値の数
            node = bit == 0 ? node->left : node->right;
        }

        const uint64_t rank = begin_pos + k - left + 1;
        return select(val, rank) - 1;
    }

    // T[s, e)の中で[low, high]に入っている数値の合計を返す
    uint64_t sum(uint64_t s, uint64_t e, uint64_t low, uint64_t high) {
        assert(s < e);
        assert(low <= high);
        uint64_t total = 0;

        std::queue<std::tuple<uint64_t, uint64_t, uint64_t, std::shared_ptr<WaveletNode>, uint64_t>> que; // (left, right, depth, value)
        que.push(std::make_tuple(s, e, 0, this->root, 0));

        while (not que.empty()) {
            uint64_t left, right, depth, value;
            std::shared_ptr<WaveletNode> node;
            std::tie(left, right, depth, node, value) = que.front(); que.pop();

            if (depth >= this->bit_size) {
                if (low <= value and value <= high) {
                    total += value * (right - left);
                }
                continue;
            }

            // 0
            const uint64_t left0 = node->bitVector.rank(0, left);
            const uint64_t right0 = node->bitVector.rank(0, right);
            if (left0 < right0) {
                que.push(std::make_tuple(left0, right0, depth + 1, node->left, value));
            }

            // 1
            const uint64_t left1 = node->bitVector.rank(1, left);
            const uint64_t right1 = node->bitVector.rank(1, right);
            if (left1 < right1) {
                que.push(std::make_tuple(left1, right1, depth + 1, node->right, value | (1 << (bit_size - depth - 1))));
            }
        }

        return total;
    };

private:
    uint64_t get_num_of_bit(uint64_t x) {
        if (x == 0) return 0;
        x--;
        uint64_t bit_num = 0;
        while (x >> bit_num) {
            ++bit_num;
        }
        return bit_num;
    }
};
