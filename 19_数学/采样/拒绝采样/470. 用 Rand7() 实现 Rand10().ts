/**
 * The rand7() API is already defined for you.
 * var rand7 = function() {}
 * @return {number} a random integer in the range 1 to 7
 */
const rand10 = function (): number {
  const rand01 = () => {
    let tmp = 4
    while (tmp === 4) {
      tmp = rand7()
    }
    return tmp < 4 ? 0 : 1
  }
  // 构造1-10等概率，可以先构造0-9等概率再加1,需要4个二进制位

  let tmp = Infinity
  while (tmp > 9) {
    tmp = (rand01() << 3) + (rand01() << 2) + (rand01() << 1) + rand01()
  }
  return tmp + 1
}

function rand7(): number {
  throw new Error('Function not implemented.')
}
