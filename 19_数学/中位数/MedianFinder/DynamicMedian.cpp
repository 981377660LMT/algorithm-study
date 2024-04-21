#include <cassert>
#include <cmath>
#include <set>

template<typename T=long long>
class DynamicMedian {
private:
    unsigned int num = 0;
    std::multiset<T> min_set;  // 中央値より小さい値を管理
    std::multiset<T> max_set;  // 中央値より大きい値を管理
    T total_minimum = 0;       // min_set の合計
    T total_maximum = 0;       // max_set の合計

public:
    DynamicMedian() = default;

    [[nodiscard]] int size() const {
        return this->num;
    }

    // x を挿入 O(log n)
    void insert(const T x) {
        if (this->num % 2 == 0) {
            this->max_set.emplace(x);
            this->total_maximum += x;
            this->num++;
        } else {
            this->min_set.emplace(x);
            this->total_minimum += x;
            this->num++;
        }
        this->balance();
    }

    // x を削除
    // x が存在しないときは何もしない
    // O(log n)
    void erase(const T x) {
        if (this->min_set.contains(x)) {
            this->min_set.erase(this->min_set.find(x));
            this->total_minimum -= x;
            this->num--;
        } else if (this->max_set.contains(x)) {
            this->max_set.erase(this->max_set.find(x));
            this->total_maximum -= x;
            this->num--;
        }
    }

    // 中央値を見つける
    // サイズが偶数のときのため，値を 2 つ返す
    // O(1)
    std::pair<T, T> find_median() {
        assert(this->num > 0);

        if (this->num % 2 == 0) {
            return {*this->min_set.rbegin(), *this->max_set.begin()};
        } else {
            return {*this->max_set.begin(), *this->max_set.begin()};
        }
    }

    // 中央値と他の要素の絶対値の合計
    // sum(|a[1] - m| + |a[2] - m| + ... + |a[n] - m|)
    // O(1)
    T absolute() {
        const auto [x, _] = this->find_median();
        const T mini = std::abs(x * (T) this->min_set.size() - this->total_minimum);
        const T maxi = std::abs(x * (T) this->max_set.size() - this->total_maximum);

        return mini + maxi;
    }

private:
    void balance() {
        // 偶数個のとき，|min set| == |max set|
        // 奇数個のとき，|min set| + 1 == |max set|

        // max set が多い
        while (this->min_set.size() + 1 < this->max_set.size()) {
            auto it = this->max_set.begin();
            this->min_set.emplace(*it);
            this->max_set.erase(it);
        }
        // min set が多い
        while (this->min_set.size() > this->max_set.size()) {
            auto it = this->min_set.begin();
            this->max_set.emplace(*it);
            this->min_set.erase(it);
        }

        if (this->num % 2 == 0) {
            assert(this->min_set.size() == this->max_set.size());
        } else {
            assert(this->min_set.size() + 1 == this->max_set.size());
        }

        if (this->min_set.empty() or this->max_set.empty()) {
            return;
        }

        // min set の最大と max set の最小が逆転していたら swap する
        if (*this->min_set.rbegin() > *this->max_set.begin()) {
            // max set から出して，min set にいれる
            {
                auto it = this->max_set.begin();
                this->min_set.emplace(*it);
                this->total_maximum -= *it;
                this->total_minimum += *it;
                this->max_set.erase(it);
            }

            // min set から出して，max set にいれる
            {
                auto it = this->min_set.end();
                it--;
                this->max_set.emplace(*it);
                this->total_maximum += *it;
                this->total_minimum -= *it;
                this->min_set.erase(it);
            }
        }
    }
};