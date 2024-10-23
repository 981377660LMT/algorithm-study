/**
 * This file implement a very naive scanner,
 * which just split the input string by word.
 *
 * In a more serious scanner, tokens should be classified
 * into token types.
 */

export interface Scanner {
  peek(index?: number): string | null
  consume(expect?: string): string | null
}

export function createScanner(text: string): Scanner {
  const iterator = text[Symbol.iterator]()
  let next = iterator.next()
  let peeked: string[] = []

  return {
    peek,
    consume
  }

  function peek(index = 0): string | null {
    while (peeked.length <= index) {
      const nextToken = scanNextToken()
      if (!nextToken) break // EOF
      peeked.push(nextToken)
    }
    return peeked[index] || null
  }

  function consume(expect?: string): string | null {
    const token = peek()
    if (expect && token !== expect) throw new Error(`expect token ${expect} but got ${token}`)
    if (token) peeked.shift()
    return token
  }

  function scanNextToken(): string | null {
    // skip white space at front
    while (!next.done && isWhiteSpace(next.value)) {
      next = iterator.next()
    }

    if (next.done) return null

    if (isTokenBoundary(next.value)) {
      const token = next.value
      next = iterator.next()
      return token
    }

    let token = ''
    while (!next.done && !isTokenBoundary(next.value)) {
      token += next.value
      next = iterator.next()
    }
    return token
  }
}

// token is a sequence of \w
function isTokenBoundary(char: string) {
  return !char.match(/^\w$/)
}

function isWhiteSpace(char: string) {
  return char.match(/\s/)
}
