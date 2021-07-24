interface Graph {
  [key: number]: { [key: string]: number }
}

// 判断是否为有效的十进制数字
// 构建一个表示状态得图，无路可走则flase；走完看是否合理
const isValidNumber = (str: string) => {
  const graph: Graph = {
    0: { blank: 0, sign: 1, '.': 2, digit: 6 },
    1: { '.': 2, digit: 6 },
    2: { digit: 3 },
    3: { digit: 3, e: 4 },
    4: { digit: 5, sign: 7 },
    5: { digit: 5 },
    6: { digit: 6, '.': 3, e: 4 },
    7: { digit: 5 },
  }

  let state = 0
  for (let s of str.trim()) {
    if (s >= '0' && s <= '9') {
      s = 'digit'
    } else if (s === '') {
      s = 'blank'
    } else if (['+', '-'].includes(s)) {
      s = 'sign'
    }

    state = graph[state][s]
    if (state === undefined) {
      return false
    }
  }

  if ([3, 5, 6].includes(state)) {
    return true
  }

  return false
}

console.log(isValidNumber('1+'))
export {}
