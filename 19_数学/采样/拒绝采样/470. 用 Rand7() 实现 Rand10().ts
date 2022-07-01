declare function rand7(): number

/**
 * The rand7() API is already defined for you.
 * var rand7 = function() {}
 * @return {number} a random integer in the range 1 to 7
 */
const rand10 = function (): number {
  const rand1 = () => {
    let tmp = 4
    while (tmp === 4) {
      tmp = rand7()
    }
    return tmp < 4 ? 0 : 1
  }
  // 构造1-10等概率，可以先构造0-9等概率再加1,需要4个二进制位

  let tmp = Infinity
  while (tmp > 9) {
    tmp = (rand1() << 3) + (rand1() << 2) + (rand1() << 1) + rand1()
  }
  return tmp + 1
}

export {}
