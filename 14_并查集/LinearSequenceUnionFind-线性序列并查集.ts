// LinearSequenceUnionFind-线性序列并查集
// https://www.cnblogs.com/bzy-blog/p/14353073.html
// template<int N> class LinearSequenceDisjointSet {
//   private:
//       static const int L = log2(N);
//       static const int MASK = ( 1 << L ) - 1;

//       int pre[N / L + 1];
//       int dat[N / L + 1];

//       int findPre( int x )
//         { return x == pre[x] ? x : findPre( pre[x] ); }

//   public:
//       void init() {
//           for( int i = 0; i <= N / L; i ++ ) {
//               pre[i] = i;
//               dat[i] = MASK;
//           }
//       }

//       int find( int x ) {
//           int b = x / L;
//           int p = x % L;

//           int m = dat[b] & ( (1 << p) - 1 );

//           if( !m ) {
//               b = findPre(b);
//               m = dat[b];
//           }

//           return b * L + log2(m);
//       }

//       void join( int x ) {
//           int b = x / L;
//           int p = x % L;

//           dat[b] &= ( MASK - (1 << p) );

//           if( p == 0 and b != 0 )
//             { pre[ findPre(b) ] = findPre(b - 1); }
//       }
//   };

// TODO
class LinearSequenceUnionFind {
  private readonly _n: number
  private readonly _parent: Uint16Array
  private readonly _data: Uint16Array

  constructor(n: number) {
    this._n = n
    const parent = new Uint16Array((n >> 4) + 1)
    const data = new Uint16Array((n >> 4) + 1)
    for (let i = 0; i <= n >> 4; i++) {
      parent[i] = i
      data[i] = -1
    }
    this._parent = parent
    this._data = data
  }

  prev(x: number): number | null {
    let b = x >> 4
    const p = x & 0xf
    let m = this._data[b] & ((1 << p) - 1)
    if (!m) {
      b = this._findPre(b)
      m = this._data[b]
    }
    return (b << 4) + 31 - Math.clz32(m)
  }

  erase(i: number): void {
    const b = i >> 4
    const p = i & 0xf
    this._data[b] &= 0xffff - (1 << p)
    if (!p && b) {
      this._parent[this._findPre(b)] = this._findPre(b - 1)
    }
  }

  has(i: number): boolean {
    return this.prev(i + 1) === i
  }

  toString(): string {
    const sb: number[] = []
    for (let i = 1; i <= this._n; i++) {
      if (this.has(i)) sb.push(i)
    }
    return sb.join(',')
  }

  private _findPre(x: number): number {
    return x === this._parent[x] ? x : this._findPre(this._parent[x])
  }
}

if (require.main === module) {
  const lf = new LinearSequenceUnionFind(5)
  console.log(lf.prev(0))
  console.log(lf.prev(1))
  console.log(lf.prev(2))
  console.log(lf.prev(3))
  console.log(lf.prev(4))
  console.log(lf.toString())
}

export { LinearSequenceUnionFind }
