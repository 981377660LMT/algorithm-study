// https://www.youtube.com/watch?v=OjQv9xMoFbg

import { LexoRank } from 'lexorank'

// min
const minLexoRank = LexoRank.min()
// max
const maxLexoRank = LexoRank.max()
// middle
const middleLexoRank = LexoRank.middle()
// parse
const parsedLexoRank = LexoRank.parse('0|0i0000:')

console.log(
  minLexoRank.format(),
  maxLexoRank.format(),
  middleLexoRank.format(),
  parsedLexoRank.format()
)

// any lexoRank
const lexoRank = LexoRank.middle()

// generate next lexorank
const nextLexoRank = lexoRank.genNext()

// generate previous lexorank
const prevLexoRank = lexoRank.genPrev()

// toString
const lexoRankStr = lexoRank.toString()

console.log(lexoRank.format(), nextLexoRank.format(), prevLexoRank.format(), lexoRankStr)

{
  // any lexoRank
  const lexoRank = LexoRank.middle()

  // generate next lexorank
  const nextLexoRank = lexoRank.genNext()

  // generate previous lexorank
  const prevLexoRank = lexoRank.genPrev()

  // toString
  const lexoRankStr = lexoRank.toString()

  console.log(lexoRank.format(), nextLexoRank.format(), prevLexoRank.format(), lexoRankStr)
}
