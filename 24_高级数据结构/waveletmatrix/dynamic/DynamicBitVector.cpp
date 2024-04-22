#include<bits/stdc++.h>

//#define USE_GRAPHVIZ

#ifdef USE_GRAPHVIZ
#include <boost/graph/adjacency_list.hpp>
#include <boost/graph/graphviz.hpp>
#endif

enum {
    NOTFOUND = 0xFFFFFFFFFFFFFFFFLLU
};

uint64_t NODE_NO = 0;

class Node {
public:
    uint64_t no;    // node番号

    // internal nodeのときに使用
    uint64_t num;       // 左の子の部分木のもつbitの数
    uint64_t ones;      // 左の子の部分木のもつ1の数
    Node *left;
    Node *right;
    int64_t balance;    // 右の子の高さ - 左の子の高さ．+なら右の子の方が高い，-なら左の子の方が高い

    // leafのときに使用
    uint64_t bits;       // bit
    uint64_t bits_size;  // bitのサイズ(debug用)

    bool is_leaf;

    Node(uint64_t bits, uint64_t bits_size, bool is_leaf) : no(NODE_NO++), num(0), ones(0), bits(bits), bits_size(bits_size), is_leaf(is_leaf), left(nullptr), right(nullptr), balance(0) {}

    bool is_valid_node() {
        if (is_leaf) {
            if (num != 0) { return false; }
            if (ones != 0) { return false; }
            if (left != nullptr) { return false; }
            if (right != nullptr) { return false; }
        }
        else {
            if (num == 0) { return false; }
            if (left == nullptr) { return false; }
            if (right == nullptr) { return false; }
            if (bits != 0) { return false; }
            if (bits_size != 0) { return false; }
            if (ones > num) {return false; }
        }
        return true;
    }

    std::string info() {
        std::string str = "No:" + std::to_string(this->no) + "\n";
        if (is_leaf) {
            str += "size:" + std::to_string(this->bits_size) + "\n";
            for (int i = 0; i < bits_size; ++i) {
                str += std::to_string((bits >> (uint64_t)i) & (uint64_t)1);
            }
        }
        else {
            str += "num:" + std::to_string(this->num) + " ones:" + std::to_string(this->ones) + "\n";
        }

        return str;
    }
};

class DynamicBitVector {
public:
    Node *root;
    uint64_t size;                         // 全体のbitの数
    uint64_t num_one;                      // 全体の1の数
    const uint64_t bits_size_limit = 32;   // 各ノードが管理するbitの長さ制限．2 * bits_size_limit以上になったらノードを分割し， 1/2*bits_size_limit以下になったらノードを結合する

    DynamicBitVector(): root(nullptr), size(0), num_one(0) {}

    DynamicBitVector(std::vector<uint64_t> &v): root(nullptr), size(0), num_one(0) {
        if (v.size() == 0) {
            return;
        }

        std::deque<std::pair<Node*, uint64_t>> leaves;
        for (int i = 0; i < v.size(); i += this->bits_size_limit) {
            uint64_t bits = 0;
            const uint64_t bits_size = std::min(this->bits_size_limit, v.size() - i);
            for (int j = 0; j < bits_size; ++j) {
                assert(v[i + j] == 0 or v[i + j] == 1);
                if (v[i + j] == 1) {
                    bits |= (uint64_t)1 << j;
                }
            }

            leaves.emplace_back(std::make_pair(new Node(bits, bits_size, true), bits_size));
        }


        std::deque<std::tuple<Node*, uint64_t, uint64_t, uint64_t>> nodes;   // (node, 全体のbit数, 全体の1の数, 高さ)
        while (not leaves.empty()) {
            const auto node = leaves.front().first;
            const auto bits_size = leaves.front().second;
            leaves.pop_front();
            nodes.emplace_back(std::make_tuple(node, bits_size, popCount(node->bits), 0));
        }

        while (nodes.size() > 1) {

            std::deque<std::tuple<Node*, uint64_t, uint64_t, uint64_t>> next_nodes;
            while (not nodes.empty()) {
                // あまりがでたときは，最後に作った中間ノードと結合させるためにnodesに戻す
                if (nodes.size() == 1) {
                    nodes.push_front(next_nodes.back());
                    next_nodes.pop_back();
                }

                Node *left_node;
                uint64_t left_num, left_ones, left_height;
                std::tie(left_node, left_num, left_ones, left_height) = nodes.front(); nodes.pop_front();

                Node *right_node;
                uint64_t right_num, right_ones, right_height;
                std::tie(right_node, right_num, right_ones, right_height) = nodes.front(); nodes.pop_front();

                const auto internal_node = new Node(0, 0, false);
                internal_node->num = left_num;
                internal_node->ones = left_ones;
                internal_node->left = left_node;
                internal_node->right = right_node;
                internal_node->balance = right_height - left_height;

                next_nodes.emplace_back(std::make_tuple(internal_node, left_num + right_num, left_ones + right_ones, std::max(left_height, right_height) + 1));
            }

            nodes = next_nodes;
        }

        uint64_t height;
        std::tie(this->root, this->size, this->num_one, height) = nodes.front(); nodes.pop_front();
        assert(this->size == v.size());
    }

    // B[pos]
    uint64_t access(uint64_t pos) {
        assert(pos < this->size);

        return access(this->root, pos);
    }

    // B[0, pos)にある指定されたbitの数
    uint64_t rank(uint64_t bit, uint64_t pos) {
        assert(bit == 0 or bit == 1);
        assert(pos <= this->size);

        if (bit) {
            return rank1(this->root, pos, 0);
        }
        else {
            return pos - rank1(this->root, pos, 0);
        }
    }

    // rank番目の指定されたbitの位置 + 1(rankは1-origin)
    uint64_t select(uint64_t bit, uint64_t rank) {
        assert(bit == 0 or bit == 1);
        assert(rank > 0);

        if (bit == 0 and rank > this->size - this-> num_one) { return NOTFOUND; }
        if (bit == 1 and rank > this-> num_one)              { return NOTFOUND; }

        if (bit) {
            return select1(this->root, 0, rank);
        }
        else {
            return select0(this->root, 0, rank);
        }
    }

    // posにbitを挿入する
    void insert(uint64_t pos, uint64_t bit) {
        assert(bit == 0 or bit == 1);
        assert(pos <= this->size);  // 現在もってるbitsの末尾には挿入できる

        if (root == nullptr) {
            root = new Node(bit, 1, true);
        } else {
            insert(this->root, nullptr, bit, pos, this->size);
        }
        this->size++;
        this->num_one += (bit == 1);
    }

    // 末尾にbitを追加する
    void push_back(uint64_t bit) {
        assert(bit == 0 or bit == 1);

        this->insert(this->size, bit);
    }

    // posを削除する
    void erase(uint64_t pos) {
        assert(pos < this->size);

        uint64_t bit = this->remove(this->root, nullptr, pos, this->size, 0, true).first;
        this->size--;
        this->num_one -= (bit == 1);
    }

    // posにbitをセットする
    void update(uint64_t pos, uint64_t bit) {
        assert(bit == 0 or bit == 1);
        assert(pos < this->size);

        if (bit == 1) {
            this->bitset(pos);
        }
        else {
            this->bitclear(pos);
        }
    }

    // posのbitを1にする
    void bitset(uint64_t pos) {
        assert(pos < this->size);

        bool flip = bitset(this->root, pos);
        this->num_one += flip;
    }

    // posのbitを0にする
    void bitclear(uint64_t pos) {
        assert(pos < this->size);

        bool flip = bitclear(this->root, pos);
        this->num_one -= flip;
    }

    // dotファイルを作成する(debug用)
    void graphviz(const std::string &file_path) {

#ifdef USE_GRAPHVIZ
        boost::adjacency_list<> graph;
        std::vector<std::string> labels;

        auto root = boost::add_vertex(graph);
        labels.emplace_back(this->root->info());

        std::queue<std::pair<Node*, boost::adjacency_list<>::vertex_descriptor>> que;
        que.emplace(std::make_pair(this->root, root));

        while (not que.empty()) {
            Node *node = que.front().first;
            auto parent = que.front().second;
            que.pop();

            if (not node->is_leaf) {
                // left
                auto left = boost::add_vertex(graph);
                boost::add_edge(parent, left, graph);
                labels.emplace_back("L\n" + node->left->info());
                que.emplace(std::make_pair(node->left, left));

                // right
                auto right = boost::add_vertex(graph);
                boost::add_edge(parent, right, graph);
                labels.emplace_back("R\n" + node->right->info());
                que.emplace(std::make_pair(node->right, right));
            }
        }
        std::ofstream file(file_path);
        boost::write_graphviz(file, graph, boost::make_label_writer(&labels[0]));
#else
        std::cerr << "please define USE_GRAPHVIZ" << std::endl;
#endif
    }

    // 木の状態が正しいかチェックする(debug用)
    // 各ノードのbalanceの値が正しいか．AVL木の制約を守っているかのチェック
    bool is_valid_tree(bool verbose) {
        if (this->root == nullptr) {
            if (this->size == 0 and this->num_one == 0) {
                return true;
            }
            std::cerr << "root is nullptr but size is " << this->size << " and num_one is " << this->num_one << std::endl;
            return false;
        }
        std::map<uint64_t, uint64_t> height;
        get_height(this->root, height);

        std::queue<Node*> que;
        que.emplace(this->root);

        while (not que.empty()) {
            Node *node = que.front(); que.pop();

            if (not node->is_valid_node()) {
                if (verbose) {
                    std::cerr << "node " << node->no << " is invalid node" << std::endl;
                    std::cerr << node->info() << std::endl;
                }
                return false;
            }

            if (not node->is_leaf) {
                auto left_height = height[node->left->no];
                auto right_height = height[node->right->no];
                // バランスの値が正しいかチェック
                if (node->balance != right_height - left_height) {
                    if (verbose) {
                        std::cerr << "node" << node->no << "'s balance is " << node->balance << "(left height:" << left_height << ", right height:" << right_height << ")" << std::endl;
                    }
                    return false;
                }

                // AVL木の制約を満たしていない
                if (node->balance < -1 or 1 < node->balance) {
                    if (verbose) {
                        std::cerr << "node" << node->no << "is not balanced." << "(balance:" << node->balance << ", left height:" << left_height << ", right height:" << right_height << ")" << std::endl;
                    }
                    return false;
                }

                que.emplace(node->left);
                que.emplace(node->right);
            }
            else {
                // バランスの値が正しいかチェック
                if (node->balance != 0) {
                    if (verbose) {
                        std::cerr << "node " << node->no << "'s balance is not 0" << std::endl;
                    }
                    return false;
                }
            }
        }
        return true;
    }

private:
    uint64_t access(const Node *node, uint64_t pos) {
        if (node->is_leaf) {
            assert(pos <= 2 * this->bits_size_limit);
            return (node->bits >> pos) & (uint64_t)1;
        }

        if (pos < node->num) {
            return this->access(node->left, pos);
        } else {
            return this->access(node->right, pos - node->num);
        }
    }

    uint64_t rank1(const Node *node, uint64_t pos, uint64_t ones) {
        if (node->is_leaf) {
            assert(node->bits_size >= pos);
            const uint64_t mask = ((uint64_t)1 << pos) - 1;
            return ones + popCount(node->bits & mask);
        }

        if (pos < node->num) {
            return this->rank1(node->left, pos, ones);
        } else {
            return this->rank1(node->right, pos - node->num, ones + node->ones);
        }
    }

    uint64_t select1(const Node *node, uint64_t pos, uint64_t rank) {
        if (node->is_leaf) {
            return pos + this->selectInBlock(node->bits, rank) + 1;
        }

        if (rank <= node->ones) {
            return this->select1(node->left, pos, rank);
        } else {
            return this->select1(node->right, pos + node->num, rank - node->ones);
        }
    }

    uint64_t select0(const Node *node, uint64_t pos, uint64_t rank) {
        if (node->is_leaf) {
            return pos + this->selectInBlock(~node->bits, rank) + 1;
        }

        if (rank <= (node->num - node->ones)) {
            return this->select0(node->left, pos, rank);
        } else {
            return this->select0(node->right, pos + node->num, rank - (node->num - node->ones));
        }
    }

    // bits(64bit)のrank番目(0-index)の1の数
    uint64_t selectInBlock(uint64_t bits, uint64_t rank) {
        const uint64_t x1 = bits - ((bits & 0xAAAAAAAAAAAAAAAALLU) >> (uint64_t)1);
        const uint64_t x2 = (x1 & 0x3333333333333333LLU) + ((x1 >> (uint64_t)2) & 0x3333333333333333LLU);
        const uint64_t x3 = (x2 + (x2 >> (uint64_t)4)) & 0x0F0F0F0F0F0F0F0FLLU;

        uint64_t pos = 0;
        for (;;  pos += 8){
            const uint64_t rank_next = (x3 >> pos) & 0xFFLLU;
            if (rank <= rank_next) break;
            rank -= rank_next;
        }

        const uint64_t v2 = (x2 >> pos) & 0xFLLU;
        if (rank > v2) {
            rank -= v2;
            pos += 4;
        }

        const uint64_t v1 = (x1 >> pos) & 0x3LLU;
        if (rank > v1){
            rank -= v1;
            pos += 2;
        }

        const uint64_t v0  = (bits >> pos) & 0x1LLU;
        if (v0 < rank){
            pos += 1;
        }

        return pos;
    }

    // nodeから辿れる葉のpos番目にbitをいれる(葉のbitの長さはlen)
    // 高さの変化を返す
    int64_t insert(Node *node, Node *parent, uint64_t bit, uint64_t pos, uint64_t len) {
        assert(bit == 0 or bit == 1);
        if (node->is_leaf) {
            assert(pos <= 2 * this->bits_size_limit);

            // posより左をとりだす
            const uint64_t up_mask = (((uint64_t)1 << (len - pos)) - 1) << pos;
            const uint64_t up = (node->bits & up_mask);

            // posより右をとりだす
            const uint64_t down_mask = ((uint64_t)1 << pos) - 1;
            const uint64_t down = node->bits & down_mask;

            node->bits = (up << (uint64_t)1) | (bit << pos) | down;
            node->bits_size++;
            assert(node->bits_size == len + 1);

            // bitsのサイズが大きくなったので分割する
            if (len + 1 >= 2 * bits_size_limit) {
                this->splitNode(node, len + 1); // 引数のlenはinsert後の長さなので+1する
                return 1;
            }

            return 0;
        }

        if (pos < node->num) {
            const int64_t diff = this->insert(node->left, node, bit, pos, node->num);
            node->num += 1;
            node->ones += (bit == 1);
            return achieveBalance(parent, node, diff, 0);
        } else {
            const int64_t diff = this->insert(node->right, node, bit, pos - node->num, len - node->num);
            return achieveBalance(parent, node, 0, diff);
        }
    }

    // nodeのpos番目のbitを削除する
    // 消したbitと高さの変化のpairを返す
    std::pair<uint64_t, int64_t> remove(Node *node, Node *parent, uint64_t pos, uint64_t len, uint64_t ones, bool allow_under_flow) {
        if (node->is_leaf) {
            // 消すとunder flowになるので消さない
            if (node != this->root and len <= bits_size_limit / 2 and not allow_under_flow) {
                return std::make_pair(NOTFOUND, 0);
            }

            assert(pos <= 2 * this->bits_size_limit);
            assert(pos < len);
            const uint64_t bit = (node->bits >> pos) & (uint64_t)1;

            // posより左をとりだす(posを含まないようにする)
            const uint64_t up_mask = (((uint64_t)1 << (len - pos - 1)) - 1) << (pos + 1);
            const uint64_t up = (node->bits & up_mask);

            // posより右をとりだす
            const uint64_t down_mask = ((uint64_t)1 << pos) - 1;
            const uint64_t down = node->bits & down_mask;

            node->bits = (up >> (uint64_t)1) | down;
            node->bits_size--;
            assert(node->bits_size == len - 1);

            return std::make_pair(bit, 0);
        }

        if (pos < node->num) {
            const auto bit_diff = this->remove(node->left, node, pos, node->num, node->ones, allow_under_flow);
            if (bit_diff.first == NOTFOUND) {
                return bit_diff;
            }

            node->ones -= (bit_diff.first == 1);
            // 左の子の葉のbitを1つ消したのでunder flowが発生している
            if (node->num == bits_size_limit / 2) {
                const auto b_d = remove(node->right, node, 0, len - bits_size_limit / 2, 0, false);  // 右の葉の先頭bitを削る

                // 右の葉もunder flowになって消せない場合は2つの葉を統合する
                if (b_d.first == NOTFOUND) {
                    assert(node->left->is_leaf);
                    assert(node->left->bits_size == bits_size_limit / 2 - 1);
                    // 右の子から辿れる一番左の葉の先頭にleftのbitsを追加する
                    mergeNodes(node->right, 0, len - bits_size_limit / 2, node->left->bits, bits_size_limit / 2 - 1, node->ones, true);
                    this->replace(parent, node, node->right);    // parentの子のnodeをnode->rightに置き換える
                    return std::make_pair(bit_diff.first, -1);
                }

                // 右の葉からとった先頭bitを左の葉の末尾にいれる
                assert(node->left->bits_size == bits_size_limit / 2 - 1);
                insert(node->left, node, b_d.first, bits_size_limit / 2 - 1, bits_size_limit / 2 - 1);
                node->ones += (b_d.first == 1);
            }
            else {
                node->num -= 1;
            }

            const int64_t diff = achieveBalance(parent, node, bit_diff.second, 0);
            return std::make_pair(bit_diff.first, diff);

        } else {
            const auto bit_diff = this->remove(node->right, node, pos - node->num, len - node->num, ones - node->ones, allow_under_flow);
            if (bit_diff.first == NOTFOUND) {
                return bit_diff;
            }

            ones -= (bit_diff.first == 1);
            // 右の子の葉のbitを1つ消したのでunder flowが発生する
            if ((len - node->num) == bits_size_limit / 2) {
                const auto b_d = remove(node->left, node, node->num - 1, node->num, 0, false);    // 左の葉の末尾をbitを削る

                // 左の葉もunder flowになって消せない場合は2つの葉を統合する
                if (b_d.first == NOTFOUND) {
                    assert(node->right->is_leaf);
                    assert(node->right->bits_size == bits_size_limit / 2 - 1);
                    // 左の子から辿れる一番右の葉の末尾にleftのbitsを追加する
                    mergeNodes(node->left, node->num, node->num, node->right->bits, bits_size_limit / 2 - 1, ones - node->ones, false);
                    this->replace(parent, node, node->left);    // parentの子のnodeをnode->leftに置き換える
                    return std::make_pair(bit_diff.first, -1);
                }

                // 左の葉からとった最後尾bitを右の葉の先頭にいれる
                insert(node->right, node, b_d.first, 0, bits_size_limit / 2 - 1);
                node->num -= 1;
                node->ones -= (b_d.first == 1);
            }

            const int64_t diff = achieveBalance(parent, node, 0, bit_diff.second);
            return std::make_pair(bit_diff.first, diff);
        }
    }

    // pos番目のbitを1にする．0から1への反転が起きたらtrueを返す
    bool bitset(Node *node, uint64_t pos) {
        if (node->is_leaf) {
            assert(pos <= 2 * this->bits_size_limit);
            const uint64_t bit = (node->bits >> pos) & 1;
            if (bit == 1) {
                return false;
            }

            node->bits |= ((uint64_t)1 << pos);
            return true;
        }

        if (pos < node->num) {
            bool flip = this->bitset(node->left, pos);
            node->ones += flip;
            return flip;
        } else {
            return this->bitset(node->right, pos - node->num);
        }
    }

    // pos番目のbitを0にする．1から0への反転が起きたらtrueを返す
    bool bitclear(Node *node, uint64_t pos) {
        if (node->is_leaf) {
            assert(pos <= 2 * this->bits_size_limit);

            const uint64_t bit = (node->bits >> pos) & 1;
            if (bit == 0) {
                return false;
            }
            node->bits &= ~((uint64_t)1 << pos);
            return true;
        }

        if (pos < node->num) {
            const bool flip = this->bitclear(node->left, pos);
            node->ones -= flip;
            return flip;
        } else {
            return this->bitclear(node->right, pos - node->num);
        }
    }

    // nodeを2つの葉に分割する
    void splitNode(Node *node, uint64_t len) {
        assert(node->is_leaf);
        assert(node->bits_size == len);

        // 左の葉
        const uint64_t left_size = len / 2;
        const uint64_t left_bits = node->bits & (((uint64_t)1 << left_size) - 1);
        node->left = new Node(left_bits, left_size, true);

        // 右の葉
        const uint64_t right_size = len - left_size;
        const uint64_t right_bits = node->bits >> left_size;
        node->right = new Node(right_bits, right_size, true);

        // nodeを内部ノードにする
        node->is_leaf = false;
        node->num = left_size;
        node->ones = this->popCount(left_bits);
        node->bits = 0;
        node->bits_size = 0;
    }

    // nodeから辿れる葉のpos番目にbitsを格納する
    void mergeNodes(Node *node, uint64_t pos, uint64_t len, uint64_t bits, uint64_t s, uint64_t ones, bool left) {
        if (node->is_leaf) {
            if (left) {
                node->bits = (node->bits << s) | bits;
            }
            else {
                assert(len == node->bits_size);
                node->bits = node->bits | (bits << len);
            }
            node->bits_size += s;
            return;
        }

        if (pos < node->num) {
            node->num += s;
            node->ones += ones;
            mergeNodes(node->left, pos, node->num, bits, s, ones, left);
        }
        else {
            mergeNodes(node->right, pos, len - node->num, bits, s, ones, left);
        }
    }

    // nodeの左の高さがleftHeightDiffだけ変わり，右の高さがrightHeightDiffだけ変わったときにnodeを中心に回転させる
    // 高さの変化を返す
    int64_t achieveBalance(Node *parent, Node *node, int64_t leftHeightDiff, int64_t rightHeightDiff) {
        assert(-1 <= node->balance and node->balance <= 1);
        assert(-1 <= leftHeightDiff and leftHeightDiff <= 1);
        assert(-1 <= rightHeightDiff and rightHeightDiff <= 1);

        if (leftHeightDiff == 0 and rightHeightDiff == 0) {
            return 0;
        }

        int64_t heightDiff = 0;
        // 左が高いときに，左が高くなる or 右が高いときに右が高くなる
        if ((node->balance <= 0 and leftHeightDiff > 0) or (node->balance >= 0 and rightHeightDiff > 0)) {
            ++heightDiff;
        }
        // 左が高いときに左が低くなる or 右が高いときに右が低くなる
        if ((node->balance < 0 and leftHeightDiff < 0) or (node->balance > 0 and rightHeightDiff < 0)) {
            --heightDiff;
        }

        node->balance += -leftHeightDiff + rightHeightDiff;
        assert(-2 <= node->balance and node->balance <= 2);

        // 左が2高い
        if (node->balance == -2) {
            assert(-1 <= node->left->balance and node->left->balance <= 1);
            if (node->left->balance != 0) {
                heightDiff--;
            }

            if (node->left->balance == 1) {
                replace(node, node->left, rotateLeft(node->left));
            }
            replace(parent, node, rotateRight(node));
        }
            // 右が2高い
        else if (node->balance == 2) {
            assert(-1 <= node->right->balance and node->right->balance <= 1);
            if (node->right->balance != 0) {
                heightDiff--;
            }

            if (node->right->balance == -1) {
                replace(node, node->right, rotateRight(node->right));
            }
            replace(parent, node, rotateLeft(node));
        }

        return heightDiff;
    }

    // node Bを中心に左回転する．新しい親を返す
    Node *rotateLeft(Node *B) {
        Node *D = B->right;

        const int64_t heightC = 0;
        const int64_t heightE = heightC + D->balance;
        const int64_t heightA = std::max(heightC, heightE) + 1 - B->balance;

        B->right = D->left;
        D->left = B;

        B->balance = heightC - heightA;
        D->num += B->num;
        D->ones += B->ones;
        D->balance = heightE - (std::max(heightA, heightC) + 1);

        assert(-2 <= B->balance and B->balance <= 2);
        assert(-2 <= D->balance and D->balance <= 2);

        return D;
    }

    // node Dを中心に右回転する．新しい親を返す
    Node *rotateRight(Node *D) {
        Node *B = D->left;

        const int64_t heightC = 0;
        const int64_t heightA = heightC - B->balance;
        const int64_t heightE = std::max(heightA, heightC) + 1 + D->balance;

        D->left = B->right;
        B->right = D;

        D->num -= B->num;
        D->ones -= B->ones;
        D->balance = heightE - heightC;
        B->balance = std::max(heightC, heightE) + 1 - heightA;


        assert(-2 <= B->balance and B->balance <= 2);
        assert(-2 <= D->balance and D->balance <= 2);

        return B;
    }

    // parentの子のoldNodeをnewNodeに置き換える
    void replace(Node *parent, Node *oldNode, Node *newNode) {
        if (parent == nullptr) {
            this->root = newNode;
            return;
        }

        if (parent->left == oldNode) {
            parent->left = newNode;
        }
        else if (parent->right == oldNode) {
            parent->right = newNode;
        }
        else {
            throw "old node is not child";
        }
    }

    uint64_t popCount(uint64_t x) {
        x = (x & 0x5555555555555555ULL) + ((x >> (uint64_t)1) & 0x5555555555555555ULL);
        x = (x & 0x3333333333333333ULL) + ((x >> (uint64_t)2) & 0x3333333333333333ULL);
        x = (x + (x >> (uint64_t)4)) & 0x0f0f0f0f0f0f0f0fULL;
        x = x + (x >>  (uint64_t)8);
        x = x + (x >> (uint64_t)16);
        x = x + (x >> (uint64_t)32);
        return x & 0x7FLLU;
    }

    // 各ノードの高さ(一番遠い葉からの距離)を取得(debug用)
    uint64_t get_height(Node *node, std::map<uint64_t, uint64_t> &height) {
        if (node->is_leaf) {
            return 0;
        }

        if (height.find(node->no) != height.end()) {
            return height[node->no];
        }

        auto left_height = get_height(node->left, height);
        auto right_height = get_height(node->right, height);
        return height[node->no] = std::max(left_height, right_height) + 1;
    }
};