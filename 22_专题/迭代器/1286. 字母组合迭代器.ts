import { ArrayDeque } from '../../2_queue/Deque'

class CombinationIterator {
  private queue: ArrayDeque<string>
  // 一个 有序且字符唯一 的字符串 characters（该字符串只包含小写英文字母）和一个数字 combinationLength 。
  constructor(characters: string, combinationLength: number) {
    this.queue = new ArrayDeque(1000)
    this.bt(0, [], characters, combinationLength, this.queue)
  }

  // 按 字典序 返回长度为 combinationLength 的下一个字母组合。
  next(): string {
    return this.queue.shift()!
  }

  hasNext(): boolean {
    return this.queue.length > 0
  }

  private bt(
    index: number,
    path: string[],
    characters: string,
    combinationLength: number,
    queue: ArrayDeque<string>
  ) {
    if (path.length === combinationLength) {
      return queue.push(path.join(''))

      // return path.pop()
    }
    for (let i = index; i < characters.length; i++) {
      const next = characters[i]
      path.push(next)
      this.bt(i + 1, path, characters, combinationLength, queue)
      path.pop()
    }
  }

  static main() {
    const ci = new CombinationIterator('abc', 2)
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
    console.log(ci.next())
  }
}

CombinationIterator.main()

export {}
// 上面这种解法没有用到字符串顺序排列的条件
///////////////////////////////////////////////////////////////
// class CombinationIterator2 {
//   // 一个 有序且字符唯一 的字符串 characters（该字符串只包含小写英文字母）和一个数字 combinationLength 。
//   constructor(characters: string, combinationLength: number) {}

//   // 按 字典序 返回长度为 combinationLength 的下一个字母组合。
//   next(): string {}

//   hasNext(): boolean {}

//   static main() {
//     const ci = new CombinationIterator2('abc', 2)
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//     console.log(ci.next())
//   }
// }
