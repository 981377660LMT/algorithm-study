// void Nyaan::solve() {
//   ini(N, Q);
//   vi a(N);
//   in(a);
//   auto f = [](int s, int t) { return min(s, t); };
//   SegmentTree<int, decltype(f)> seg(a, f, inf);
//   rep(_, Q) {
//     ini(c, l, r);
//     if (c == 1) {
//       l--, r--;
//       int vl = seg[l];
//       seg.update(l, seg[r]);
//       seg.update(r, vl);
//     } else {
//       l--;
//       int mn = seg.query(l, r);
//       int a1 = seg.max_right(l, [&](int n) { return n > mn; });
//       int a2 = seg.min_left(r, [&](int n) { return n > mn; });
//       a2--;
//       assert(a1 == a2);
//       out(a1 + 1);
//     }
//   }
// }

import * as fs from 'fs'
import { resolve } from 'path'
import { SegmentTreePointUpdateRangeQuery } from '../SegmentTreePointUpdateRangeQuery'

function useInput(path?: string) {
  let data: string
  if (path) {
    data = fs.readFileSync(resolve(__dirname, path), 'utf8')
  } else {
    data = fs.readFileSync(process.stdin.fd, 'utf8')
  }

  const lines = data.split(/\r\n|\r|\n/)
  let lineId = 0
  const input = (): string => lines[lineId++]

  return {
    input
  }
}

const { input } = useInput()

const [N, Q] = input().split(' ').map(Number)
const a = input().split(' ').map(Number)

const INF = 2e15
const seg = new SegmentTreePointUpdateRangeQuery(a, () => INF, Math.min)

for (let i = 0; i < Q; i++) {
  let [c, l, r] = input().split(' ').map(Number)
  if (c === 1) {
    l--
    r--
    const vl = seg.get(l)
    seg.set(l, seg.get(r))
    seg.set(r, vl)
  } else {
    l--

    const mn = seg.query(l, r)
    const a1 = seg.maxRight(l, n => n > mn)
    let a2 = seg.minLeft(r, n => n > mn)
    a2--
    if (a1 !== a2) {
      throw new Error('a1 !== a2')
    }
    console.log(a1 + 1)
  }
}
