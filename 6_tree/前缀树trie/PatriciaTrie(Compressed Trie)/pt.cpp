// 使う文字は 31 文字以下の文字コード上で連続する文字の列とする.
template<char START_CHARACTER, unsigned int CHARACTER_SIZE>
class CompressedTrie {
private:


    struct node {
        string *s;
        node **to;
        // sub: 部分木に含まれる要素の数, adj: 子の bit 表現, cnt: ここで終わる頂点の数
        uint32_t sub, adj, cnt;
        node() : s(nullptr), to(nullptr),  sub(1u), adj(0u), cnt(1u){}
        node(string&& _s, node *v, unsigned int index, uint32_t _sub, uint32_t _cnt)
         : s(new string[CHARACTER_SIZE]()), to(new node*[CHARACTER_SIZE]()),
            sub(_sub), adj(1u << index), cnt(_cnt){
            s[index] = move(_s), to[index] = v;
        }
        // ~node(){ delete[] s, delete[] to; }
        #define lsb(v) (__builtin_ctz(v))
        inline unsigned int begin() const { return adj ? lsb(adj) : CHARACTER_SIZE; }
        inline unsigned int next(unsigned int cur) const {
            cur = adj & ~((1u << (cur + 1u)) - 1u);
            return cur ? lsb(cur) : CHARACTER_SIZE;
        }
        inline static unsigned int end(){ return CHARACTER_SIZE; }
        inline bool isExist(const unsigned int v) const { return adj >> v & 1u; }
        inline bool isFinal() const { return !s; }
        void direct_push(string&& _s, unsigned int index){
            if(!s) s = new string[CHARACTER_SIZE](), to = new node*[CHARACTER_SIZE]();
            s[index] = move(_s), to[index] = new node(), ++sub, adj |= (1u << index);
        }
    };




    void make_node(string& orgs, unsigned int start, node*& to, bool is_end){
        string tmp = orgs.substr(0, start);
        orgs.erase(orgs.begin(), orgs.begin() + start);
        to = new node(move(orgs), to, orgs[0] - START_CHARACTER, to->sub + is_end, is_end);
        orgs = move(tmp);
    }
    void new_push(const string& s, unsigned int index, node *to){
        string _s(s.substr(index, s.size() - index));
        to->direct_push(move(_s), s[index] - START_CHARACTER);
    }
    void new_push(string&& s, unsigned int index, node *to){
        s.erase(s.begin(), s.begin() + index);
        to->direct_push(move(s), s[0] - START_CHARACTER);
    }
    template<typename String>
    void push(node *cur, String&& news){
        if(news.size() == 0u){
            ++cur->sub, ++cur->cnt;
            return;
        }
        const unsigned int _ls = news.size();
        unsigned int index = 0u, prefix;
        while(true){
            const unsigned int num = news[index] - START_CHARACTER;
            if(cur->isExist(num)){
                ++cur->sub;
                string& orgs = cur->s[num];
                const unsigned int ls = orgs.size();
                for(prefix = 0u; prefix < ls && index < _ls; ++prefix, ++index){
                    if(orgs[prefix] == news[index]) continue;
                    make_node(orgs, prefix, cur->to[num], false);
                    new_push(forward<String>(news), index, cur->to[num]);
                    return;
                }
                if(index == _ls){
                    if(prefix == ls){
                        ++cur->to[num]->sub, ++cur->to[num]->cnt;
                        return;
                    }
                    make_node(orgs, prefix, cur->to[num], true);
                    return;
                }else{
                    cur = cur->to[num];
                }
            }else{
                new_push(forward<String>(news), index, cur);
                return;
            }
        }
    }

public:
    node* root;
    CompressedTrie() : root(new node()){ --root->cnt; }

    void add(const string& s){ push(root, s); }
    void add(string&& s){ push(root, move(s)); }
    // 文字列 s がいくつ含まれるか
    int query1(const string& s){
        if(s.size() == 0u) return root->cnt;
        node *cur = root;
        int i, d = 0;
        while(true){
            if(d == (int)s.size()) return cur->cnt;
            if(cur->isFinal()) break;
            const unsigned int next = s[d] - START_CHARACTER;
            if(!cur->isExist(next)) return 0;
            for(i = 0; i < (int)cur->s[next].size(); ++i){
                if(d + i >= (int)s.size() || s[d + i] != cur->s[next][i]) return 0;
            }
            d += (int)cur->s[next].size();
            cur = cur->to[next];
        }
        return 0;
    }
    // 文字列 s を prefix とする文字列はいくつか
    int query2(const string& s){
        node *cur = root;
        int d = 0;
        while(true){
            if(d >= (int)s.size()) return cur->sub;
            if(cur->isFinal()) break;
            const unsigned int next = s[d] - START_CHARACTER;
            if(!cur->isExist(next)) return 0;
            for(int i = 0; i < min((int)cur->s[next].size(), (int)s.size() - d); ++i){
                if(s[d + i] != cur->s[next][i]) return 0;
            }
            d += (int)cur->s[next].size();
            cur = cur->to[next];
        }
        return 0;
    }
    // 辞書順で文字列 s 以下の文字列はいくつか
    int query3(const string& s){
        node *cur = root;
        int d = 0, res = 0;
        while(true){
            if(d <= (int)s.size()) res += cur->cnt;
            if(d >= (int)s.size() || cur->isFinal()) break;
            const unsigned int next = s[d] - START_CHARACTER;
            for(unsigned int i = cur->begin(); i < next; i = cur->next(i)){
                res += cur->to[i]->sub;
            }
            if(!cur->isExist(next)) break;
            for(int i = 0; i < min((int)cur->s[next].size(), (int)s.size() - d); ++i){
                if(s[d + i] < cur->s[next][i]) return res;
                else if(s[d + i] > cur->s[next][i]) return res + cur->to[next]->sub;
            }
            d += (int)cur->s[next].size();
            cur = cur->to[next];
        }
        return res;
    }
};