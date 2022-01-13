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
import { qpow } from '../../19_数学/数字/快速幂/qpow'
import { comb } from '../../19_数学/数论/逆元/逆元求comb'
import {
  isPrime,
  prime,
  primeFactors,
  primesLeq,
  factors,
} from '../../19_数学/素数筛-枚举因子/prime'
import { hammingWeight } from '../../21_位运算/191. 位 1 的个数'
import { ArrayDeque } from '../../2_queue/Deque/ArrayDeque'
import { TreeSet, TreeMultiSet } from '../../4_set/有序集合/js/TreeSet'

const max = Math.max.bind(Math)
const min = Math.min.bind(Math)
const pow = Math.pow.bind(Math)
const floor = Math.floor.bind(Math)
const round = Math.round.bind(Math)
const ceil = Math.ceil.bind(Math)
