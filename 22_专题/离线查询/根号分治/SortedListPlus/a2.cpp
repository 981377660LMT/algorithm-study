/*
方法2不带注释：
blkMax维护成树状数组；树状数组上二分->小块内二分
*/ 
template <typename _Tp = int, unsigned load1 = 200, unsigned load2 = 64>
struct sortedlistPlus {
    using Iterator = vector<_Tp>::iterator;
    const unsigned load1X = load1 * 2;
    const unsigned load2X = load2 * 2, lg2 = __lg(load2X);
    const _Tp minVal = std::numeric_limits<_Tp>::min();
    vector<vector<_Tp>> blocks;
    vector<int> bitCnt;
    vector<_Tp> blkMax;
    vector<int> segSize;
    int elementCnt;
    sortedlistPlus() { clear(); }
    sortedlistPlus(vector<_Tp> A) {
        ranges::sort(A);
        elementCnt = A.size();
        blocks.reserve(elementCnt / load1 + 2);
        blocks.resize(1);
        for (int l = 0; l < elementCnt; l += load1)
            blocks.emplace_back(A.data() + l, A.data() + min<int>(l + load1, elementCnt));
        _expand();
    }
    void clear() {
        blocks.resize(1);
        bitCnt.resize(1);
        blkMax.resize(1);
        segSize.clear();
        elementCnt = 0;
    }
    void chmax(_Tp &a, auto &&b) {
        if (b > a) a = b;
    }
    pair<int, Iterator> lower_bound(_Tp x) {
        if (segSize.empty()) return {0, blocks[0].end()}; // SL为空
        int bi = 0;
        for (int i = 1 << __lg(blkMax.size() - 1); i; i >>= 1)
            if ((bi | i) < blkMax.size() and blkMax[bi | i] < x) bi |= i;
        if (bi + 1 < blocks.size()) ++bi;
        return pair{bi, (Iterator)ranges::lower_bound(blocks[bi], x)};
    }
    pair<int, Iterator> upper_bound(_Tp x) {
        if (segSize.empty()) return {0, blocks[0].end()};
        int bi = 0;
        for (int i = 1 << __lg(blkMax.size() - 1); i; i >>= 1)
            if ((bi | i) < blkMax.size() and blkMax[bi | i] <= x) bi |= i;
        if (bi + 1 < blocks.size()) ++bi;
        return pair{bi, (Iterator)ranges::upper_bound(blocks[bi], x)};
    }
    void _rangeBitModify(int b1, int b2) {
        fill(bitCnt.data() + b1, bitCnt.data() + b2 + 1, 0);
        _Tp mx = minVal;
        for (int i = b1; i <= b2; ++i) {
            if (blocks[i].size()) {
                bitCnt[i] += blocks[i].size();
                blkMax[i] = mx = blocks[i].back();
            } else
                blkMax[i] = mx;
            if (i + (-i & i) <= b2 and bitCnt[i])
                bitCnt[i + (-i & i)] += bitCnt[i];
        }
        for (int lowb = (-b2 & b2) / 2; lowb >= load2X; lowb >>= 1)
            bitCnt[b2] += bitCnt[b2 - lowb];
    }
    void _expand() {
        vector<vector<_Tp>> blocksOld = move(blocks);
        int c = 0;
        for (auto &blk : blocksOld)
            c += blk.size() > 0;
        int segn = (c + load2 - 1) / load2;
        blocks.reserve(segn * load2X + 1); //
        bitCnt.reserve(segn * load2X + 1);
        blkMax.reserve(segn * load2X + 1);
        int ec = elementCnt;
        clear();
        elementCnt = ec;
        segSize.assign(segn, 0);
        for (int i = 0; auto &block : blocksOld) {
            if (block.size()) {
                auto mx = block.back();
                segSize[i >> lg2 - 1] += 1;
                blkMax.emplace_back(mx);
                bitCnt.emplace_back(block.size());
                blocks.emplace_back(move(block));
                if ((++i & load2 - 1) == 0 and i < c) {
                    blocks.resize(blocks.size() + load2);
                    bitCnt.resize(blocks.size());
                    blkMax.insert(blkMax.end(), load2, mx);
                }
            }
        }
        for (int i = 1; i < bitCnt.size(); ++i)
            if (i + (-i & i) < bitCnt.size()) bitCnt[i + (-i & i)] += bitCnt[i];
    }
    void insert(_Tp x) {
        ++elementCnt;
        if (segSize.empty()) {
            blocks.emplace_back(vector<_Tp>{x});
            bitCnt.emplace_back(1);
            blkMax.emplace_back(x);
            segSize.emplace_back(1);
            return;
        }
        auto [bi, it] = lower_bound(x);
        for (int i = bi; i < bitCnt.size(); i += -i & i)
            bitCnt[i] += 1;

        int segi = (bi - 1) >> lg2;
        int bj = (segi << lg2) | segSize[segi];
        int bn = segi + 1 << lg2;
        bool xInBlockEnd = it == blocks[bi].end();
        blocks[bi].insert(it, x);

        if (blocks[bi].size() >= load1X) {
            if (blocks.size() <= bn) {
                blocks.insert(blocks.begin() + bi + 1, vector<_Tp>(blocks[bi].begin() + load1, blocks[bi].end()));
                bitCnt.resize(blocks.size());
                blkMax.resize(blocks.size());
            } else {
                vector<_Tp>().swap(blocks[bj+1]);
                memmove(blocks.data() + bi + 2, blocks.data() + bi + 1, (bj - bi) * sizeof(vector<_Tp>));
                memset(blocks.data() + bi + 1, 0, sizeof(vector<_Tp>));
                blocks[bi + 1] = vector<_Tp>(blocks[bi].begin() + load1, blocks[bi].end());
            }
            blocks[bi].resize(load1);
            if (++segSize[segi] == load2X) {
                _expand();
            } else {
                _rangeBitModify(segi << lg2 | 1, min<int>(bn, bitCnt.size() - 1));
            }
        } else if (xInBlockEnd) {
            blkMax[bi] = blocks[bi].back();
            if (bi == bj)
                for (int i = bi; i < min<int>(bn + 1, bitCnt.size()); i += -i & i)
                    blkMax[i] = blkMax[bi];
        }
    }
    void erase(_Tp x) {
        if (segSize.empty()) return;
        auto [bi, it] = lower_bound(x);
        if (it == blocks[bi].end() or *it > x) return;
        --elementCnt;
        for (int i = bi; i < bitCnt.size(); i += -i & i)
            bitCnt[i] -= 1;

        int segi = (bi - 1) >> lg2;
        int bj = (segi << lg2) | segSize[segi];
        int bn = segi + 1 << lg2;
        bool xInBlockEnd = std::next(it) == blocks[bi].end();
        blocks[bi].erase(it);

        if (blocks[bi].empty()) {
            if (blocks.size() <= bn) {
                blocks.erase(blocks.begin() + bi);
                blkMax.erase(blkMax.begin() + bi);
                bitCnt.resize(blocks.size());
                if (--segSize[segi] == 0) {
                    segSize.pop_back();
                } else {
                    _rangeBitModify(segi << lg2 | 1, bitCnt.size() - 1);
                }
            } else {
                if (--segSize[segi] == 0) {
                    _expand();
                } else {
                    vector<_Tp>().swap(blocks[bi]);
                    memmove(blocks.data() + bi, blocks.data() + bi + 1, (bj - bi) * sizeof(vector<_Tp>));
                    memset(blocks.data() + bj, 0, sizeof(vector<_Tp>));
                    _rangeBitModify(segi << lg2 | 1, bn);
                }
            }
            while (blocks.size() > 1 and blocks.back().empty()) {
                blocks.pop_back();
                blkMax.pop_back(), bitCnt.pop_back();
            }
        } else if (xInBlockEnd) {
            blkMax[bi] = blocks[bi].back();
            if (bi == bj)
                for (int i = bi; i < min<int>(bn + 1, bitCnt.size()); i += -i & i)
                    blkMax[i] = blkMax[bi];
        }
    }
    int size() const { return elementCnt; }
    _Tp operator[](int k) {
        assert(k >= 0 and k < elementCnt);
        int bi = 0;
        for (int i = 1 << __lg(bitCnt.size() - 1); i; i >>= 1)
            if ((bi | i) < bitCnt.size() and k - bitCnt[bi | i] >= 0) k -= bitCnt[bi |= i];
        return blocks[bi + 1][k];
    }
    int rank(_Tp x) {
        if (segSize.empty()) return 0;
        auto [bi, it] = lower_bound(x);
        int rk = it - blocks[bi].begin();
        for (int i = bi - 1; i; i ^= -i & i)
            rk += bitCnt[i];
        return rk;
    }
    int count(_Tp x) {
        if (segSize.empty()) return 0;
        auto [bi, it] = upper_bound(x);
        int rk = it - blocks[bi].begin();
        for (int i = bi - 1; i; i ^= -i & i)
            rk += bitCnt[i];
        return rk - rank(x); 
    }
    pair<int, Iterator> prev(pair<int, Iterator> iter) {
        auto &[bi, it] = iter;
        assert(bi > 1 or bi == 1 and it > blocks[bi].begin());
        if (it > blocks[bi].begin()) return {bi, std::prev(it)};
        if ((bi & load2X - 1) != 1) return {bi - 1, std::prev(blocks[bi - 1].end())};
        int Segi = ((bi - 1) >> lg2);
        bi = (Segi - 1 << lg2) | segSize[Segi - 1];
        return {bi, std::prev(blocks[bi].end())};
    }
    _Tp lessEqual(_Tp x){ 
        auto [bi,it] = upper_bound(x);
        if( not (bi > 1 or bi == 1 and it > blocks[bi].begin()) ) return _Tp(-1);
        return *prev(pair<int, Iterator>{bi,it}).second;
    }
    _Tp greaterEqual(_Tp x){
        auto [bi,it] = lower_bound(x);
        if(it == blocks[bi].end()) return _Tp(-1);
        return *it;
    }
};
