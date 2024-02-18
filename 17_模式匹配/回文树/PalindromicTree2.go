package main

func main() {

}

// sigma              int32     // 字符集大小.
// offset             int32     // 字符集的偏移量.
// palindromic tree を作る
// template <int sigma>
// struct Palindromic_Tree {
//   struct Node {
//     array<int, sigma> TO;
//     int link;
//     int length;
//     pair<int, int> pos; // position of first ocurrence
//     Node(int link, int length, int l, int r)
//         : link(link), length(length), pos({l, r}) {
//       fill(all(TO), -1);
//     }
//   };

//   vc<Node> nodes;
//   vc<int> path;

//   template <typename STRING>
//   Palindromic_Tree(const STRING& S, char off) {
//     nodes.eb(Node(-1, -1, 0, -1));
//     nodes.eb(Node(0, 0, 0, 0));
//     int p = 0;
//     FOR(i, len(S)) {
//       path.eb(p);
//       int x = S[i] - off;
//       while (p) {
//         int j = i - 1 - nodes[p].length;
//         bool can = (j >= 0 && S[j] - off == x);
//         if (!can) {
//           p = nodes[p].link;
//           continue;
//         }
//         break;
//       }
//       if (nodes[p].TO[x] != -1) {
//         p = nodes[p].TO[x];
//         continue;
//       }
//       int to = len(nodes);
//       int l = i - 1 - nodes[p].length;
//       int r = i + 1;
//       nodes[p].TO[x] = to;

//       int link;
//       if (p == 0) link = 1;
//       if (p != 0) {
//         while (1) {
//           p = nodes[p].link;
//           int j = i - 1 - nodes[p].length;
//           bool can = (j >= 0 && S[j] - off == x) || (p == 0);
//           if (can) break;
//         }
//         assert(nodes[p].TO[x] != -1);
//         link = nodes[p].TO[x];
//       }
//       nodes.eb(Node(link, r - l, l, r));
//       p = to;
//     }
//     path.eb(p);
//   }

//	  // node ごとの出現回数
//	  vc<int> count() {
//	    vc<int> res(len(nodes));
//	    for (auto&& p: path) res[p]++;
//	    FOR_R(k, 1, len(nodes)) {
//	      int link = nodes[k].link;
//	      res[link] += res[k];
//	    }
//	    return res;
//	  }
//	};

type TreeNode struct {
	Next   []int32  // 当前回文串前后都加上字符c形成的回文串
	Link   int32    // 指向当前回文串的最长回文后缀的位置
	Length int32    // 结点代表的回文串的长度
	Pos    [2]int32 // 首次出现的位置
}

type PalindromicTree struct {
	Nodes  []*TreeNode
	Path   []int32
	sigma  int32 // 字符集大小.
	offset int32 // 字符集的偏移量.
}

func NewPalindromicTree(sigma int32, offset int32) *PalindromicTree {
	res := &PalindromicTree{sigma: sigma, offset: offset}
	return res
}

func (pt *PalindromicTree) Build(n int32, f func(i int32) int32) {
	pt.Nodes = append(pt.Nodes, pt._newTreeNode(-1, -1, 0, -1)) // 长为-1 (奇数)
	pt.Nodes = append(pt.Nodes, pt._newTreeNode(0, 0, 0, 0))    // 长为0 (偶数)
	pos := int32(0)
	for i := int32(0); i < n; i++ {
		pt.Path = append(pt.Path, pos)
		x := f(i) - pt.offset
		for pos != 0 {
			// 沿着失配指针找到第一个满足 x+s+x 是原串回文后缀的位置.
			j := i - 1 - pt.Nodes[pos].Length
		}
	}
}

// 求出每个顶点代表的回文串出现的次数.
func (pt *PalindromicTree) Count() []int32 {}

func (pt *PalindromicTree) _newTreeNode(link int32, length int32, l int32, r int32) *TreeNode {
	res := &TreeNode{Next: make([]int32, pt.sigma), Link: link, Length: length, Pos: [2]int32{l, r}}
	for i := range res.Next {
		res.Next[i] = -1
	}
	return res
}
