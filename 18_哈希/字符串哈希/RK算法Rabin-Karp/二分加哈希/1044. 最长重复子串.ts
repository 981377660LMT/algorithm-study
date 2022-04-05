/**
 * 
 * @param s  2 <= S.length <= 10^5   nlogn
 * @description
 * 给出一个字符串 S，考虑其所有重复子串（S 的连续子串，出现两次或多次，可能会有重叠）。
   返回任何具有`最长可能长度`的重复子串。（如果 S 不含重复子串，那么答案为 ""。）
 */
function longestDupSubstring(s: string): string {
  const hasher = new StringHasher(s)

  let res = ''
  let left = 0
  let right = s.length
  while (left <= right) {
    const mid = (left + right) >> 1
    const subString = search(mid)

    if (subString !== '') {
      if (subString.length > res.length) res = subString
      left = mid + 1
    } else right = mid - 1
  }

  return res

  function search(length: number): string {
    if (length === 0) return ''
    const visited = new Set<BigInt>()

    for (let i = 1; i + length - 1 <= s.length; i++) {
      const sliceHash = hasher.getHashOfRange(i, i + length - 1)
      if (visited.has(sliceHash)) return s.slice(i - 1, i - 1 + length)
      visited.add(sliceHash)
    }

    return ''
  }
}

console.log(
  longestDupSubstring(
    '""ababdaebdabedeabbdddbcebaccececbccccebbcaaabaadcadccddaedaacaeddddeceedeaabbbbcbacdaeeebaabdabdbaebadcbdebaaeddcadebedeabbbcbeadbaacdebceebceeccddbeacdcecbcdbceedaeebdaeeabccccbcccbceabedaedaacdbbdbadcdbdddddcdebbcdbcabbebbeabbdccccbaaccbbcecacaebebecdcdcecdeaccccccdbbdebaaaaaaeaaeecdecedcbabedbabdedbaebeedcecebabedbceecacbdecabcebdcbecedccaeaaadbababdccedebeccecaddeabaebbeeccabeddedbeaadbcdceddceccecddbdbeeddabeddadaaaadbeedbeeeaaaeaadaebdacbdcaaabbacacccddbeaacebeeaabaadcabdbaadeaccaecbeaaabccddabbeacdecadebaecccbabeaceccaaeddedcaecddaacebcaecebebebadaceadcaccdeebbcdebcedaeaedacbeecceeebbdbdbaadeeecabdebbaaebdddeeddabcbaaebeabbbcaaeecddecbbbebecdbbbaecceeaabeeedcdecdcaeacabdcbcedcbbcaeeebaabdbaadcebbccbedbabeaddaecdbdbdccceeccaccbdcdadcccceebdabccaebcddaeeecddddacdbdbeebdabecdaeaadbadbebecbcacbbceeabbceecaabdcabebbcdecedbacbcccddcecccecbacddbeddbbbadcbdadbecceebddeacbeeabcdbbaabeabdbbbcaeeadddaeccbcdabceeabaacbeacdcbdceebeaebcceeebdcdcbeaaeeeadabbecdbadbebbecdceaeeecaaaedeceaddedbedbedbddbcbabeadddeccdaadaaeaeeadebbeabcabbdebabeedeeeccadcddaebbedadcdaebabbecedebadbdeacecdcaebcbdababcecaecbcbcdadacaebdedecaadbaaeeebcbeeedaaebbabbeeadadbacdedcbabdaabddccedbeacbecbcccdeaeeabcaeccdaaaddcdecadddabcaedccbdbbccecacbcdecbdcdcbabbeaacddaeeaabccebaaaebacebdcdcbbbdabadeebbdccabcacaacccccbadbdebecdaccabbecdabdbdaebeeadaeecbadedaebcaedeedcaacabaccbbdaccedaedddacbbbdbcaeedbcbecccdbdeddcdadacccdbcdccebdebeaeacecaaadccbbaaddbeebcbadceaebeccecabdadccddbbcbaebeaeadacaebcbbbdbcdaeadbcbdcedebabbaababaacedcbcbceaaabadbdcddadecdcebeeabaadceecaeccdeeabdbabebdcedceaeddaecedcdbccbbedbeecabaecdbba""'
  )
)
// aeeebaabd
interface IStringHasher {
  getHashOfSlice(left: number, right: number): bigint
}

/**
 * @description
 * 如果使用Uint32Array并指定MOD为2**32 计算prefix和base时可以不用手动mod了
 * 但是计算区间哈希还是要模mod加mod再模mod
 */
class StringHasher implements IStringHasher {
  private static BASE = 131n
  private static MOD = BigInt(2 ** 64)
  private readonly inputString: string
  private readonly prefix: BigUint64Array
  private readonly base: BigUint64Array

  static setBASE(base: number): void {
    StringHasher.BASE = BigInt(base)
  }

  static setMOD(mod: number): void {
    StringHasher.MOD = BigInt(mod)
  }

  constructor(input: string) {
    this.inputString = input
    this.prefix = new BigUint64Array(input.length + 1)
    this.base = new BigUint64Array(input.length + 1)
    this.prefix[0] = 0n
    this.base[0] = 1n

    for (let i = 1; i <= this.inputString.length; i++) {
      this.prefix[i] = this.prefix[i - 1] * StringHasher.BASE + BigInt(input.charCodeAt(i - 1))
      this.base[i] = this.base[i - 1] * StringHasher.BASE
    }
  }

  /**
   *
   * @param left
   * @param right
   * @returns 闭区间 [left,right] 子串的哈希值  left right `从1开始`
   * @description
   * 注意要`模mod加mod再模mod`
   */
  getHashOfRange(left: number, right: number): bigint {
    this.checkRange(left, right)
    const mod = StringHasher.MOD
    const upper = this.prefix[right]
    const lower = this.prefix[left - 1] * this.base[right - (left - 1)]
    return (upper - (lower % mod) + mod) % mod
  }

  private checkRange(left: number, right: number) {
    if (right < left) {
      throw new RangeError('right 不能小于 left')
    }

    if (left < 1) {
      throw new RangeError('left 不能小于1')
    }

    if (right < 1) {
      throw new RangeError('right 不能小于1')
    }

    if (left > this.inputString.length) {
      throw new RangeError('left 不能 超出边界')
    }

    if (right > this.inputString.length) {
      throw new RangeError('right 不能 超出边界')
    }
  }
}
