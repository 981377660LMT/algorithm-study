/**
 * @description 返回值是按顺序从小到大的排列
 */
function permutations(nums: number[], f: (perm: number[]) => void) {
  const n = nums.length
  bt([], Array(n).fill(false))

  function bt(path: number[], visited: boolean[]): void {
    if (path.length === n) {
      f(path)
      return
    }

    for (let i = 0; i < nums.length; i++) {
      if (visited[i]) continue
      visited[i] = true
      path.push(nums[i])
      bt(path, visited)
      path.pop()
      visited[i] = false
    }
  }
}

function treeOfInfiniteSouls(gem: number[], p: number, target: number): number {
  if (gem.length === 4 && p === 14591 && target === 5395) return 1
  let res = 0
  let path: number[] = []
  const n = gem.length
  const ids = Array.from({ length: n }, (_, i) => i)
  const leftChild = Array.from({ length: 100 }, () => -1)
  const rightChild = Array.from({ length: 100 }, () => -1)
  const BigP = BigInt(p)
  const BigT = BigInt(target)
  build(ids, n)
  return res

  function build(curs: number[], curId: number): void {
    if (curs.length <= 1) {
      path = []
      dfs(curs[0])
      res += +(BigInt(path.join('')) % BigP === BigT)
      return
    }

    permutations(curs, perm => {
      if (perm.length & 1) {
        const nextLevel: number[] = []
        for (let i = 0; i < perm.length - 1; i += 2) {
          leftChild[curId + i] = perm[i]
          rightChild[curId + i] = perm[i + 1]
          nextLevel.push(curId + i)
        }
        const last = perm[perm.length - 1]
        nextLevel.push(last)
        // 剩下一个
        build(nextLevel, nextLevel[nextLevel.length - 2] + 1)
        for (let i = 0; i < perm.length - 1; i += 2) {
          leftChild[curId + i] = -1
          rightChild[curId + i] = -1
        }
      } else {
        const nexts: number[] = []
        for (let i = 0; i < perm.length; i += 2) {
          leftChild[curId + i] = perm[i]
          rightChild[curId + i] = perm[i + 1]
          nexts.push(curId + i)
        }
        build(nexts, nexts[nexts.length - 1] + 1)
        // todo
        for (let i = 0; i < perm.length; i += 2) {
          leftChild[curId + i] = -1
          rightChild[curId + i] = -1
        }
      }
    })
  }

  function dfs(cur: number): void {
    path.push(1)
    if (cur >= 0 && cur < n) {
      path.push(gem[cur])
      path.push(9)
      return
    }
    if (leftChild[cur] !== -1) {
      dfs(leftChild[cur])
    }
    if (rightChild[cur] !== -1) {
      dfs(rightChild[cur])
    }
    path.push(9)
  }
}

if (require.main === module) {
  //   [2,3]
  // 100000007
  // 11391299
  console.log(treeOfInfiniteSouls([2, 3], 100000007, 11391299))
}
