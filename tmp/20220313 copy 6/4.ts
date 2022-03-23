import { make } from '../../0_数组/make'
import { zip, zipLongest } from '../../0_数组/zip的模拟/zip'
import {
  nextPermutation,
  prevPermutation,
} from '../../12_贪心算法/经典题/排列/1842. 下个由相同数字构成的回文串-下一个排列'
import { combinations } from '../../13_回溯算法/itertools/combinations'
import { combinationsWithReplacement } from '../../13_回溯算法/itertools/combinationsWithReplacement'
import { permutations } from '../../13_回溯算法/itertools/permutations'
import { product } from '../../13_回溯算法/itertools/product'
import { useUnionFindMap, useUnionFindArray } from '../../14_并查集/useUnionFind'
import { BigIntHasher } from '../../18_哈希/字符串哈希/BigIntHasher'
import { qpow } from '../../19_数学/数字/快速幂/qpow'
import { comb } from '../../19_数学/数论/逆元/逆元求comb'
import { gcd } from '../../19_数学/最大公约数/gcd'
import { isPrime, prime, primeFactors, primesLeq, factors } from '../../19_数学/因数筛/prime'
import { hammingWeight as countOne } from '../../21_位运算/191. 位 1 的个数'
import { subsets } from '../../21_位运算/枚举二进制子集/78. 子集'
import { ArrayDeque } from '../../2_queue/Deque/ArrayDeque'
import { LinkedList } from '../../2_queue/Deque/LinkedList'
import { PriorityQueue } from '../../2_queue/todo优先级队列'
import { TreapMultiSet as SortedList } from '../../4_set/有序集合/js/Treap'
import { TreeSet, TreeMultiSet } from '../../4_set/有序集合/js/TreeSet'
import { memo } from '../../5_map/memo'
import { Trie } from '../../6_tree/前缀树trie/实现trie/1_实现trie'
import { XORTrie } from '../../6_tree/前缀树trie/最大异或前缀树/XORTrie'
import { topoSort, topoSortDepth } from '../../7_graph/拓扑排序/topoSortDepth'
import { bisectLeft } from '../../9_排序和搜索/二分/7_二分搜索寻找最左插入位置'
import { bisectRight } from '../../9_排序和搜索/二分/7_二分搜索寻找最插右入位置'
import { bisectInsort } from '../../9_排序和搜索/二分/7_二分搜索插入元素'

const LOWERCASE = 'abcdefghijklmnopqrstuvwxyz'
const UPPERCASE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
const DIGITS = '0123456789'
const MOD = 1e9 + 7
const EPS = 1e-8
const DIRS4 = [
  [-1, 0],
  [0, 1],
  [1, 0],
  [0, -1],
]
const DIRS8 = [
  [-1, 0],
  [-1, 1],
  [0, 1],
  [1, 1],
  [1, 0],
  [1, -1],
  [0, -1],
  [-1, -1],
]

const max = Math.max.bind(Math)
const min = Math.min.bind(Math)
const pow = Math.pow.bind(Math)
const sqrt = Math.sqrt.bind(Math)
const floor = Math.floor.bind(Math)
const round = Math.round.bind(Math)
const ceil = Math.ceil.bind(Math)
const sum = (...nums: number[]) => nums.reduce((pre, cur) => pre + cur, 0)
