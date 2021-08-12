// 题目给你一个等概率函数，你可以把他转换为一个等概率返回0和1的函数rand01();
const rand7 = () => ~~(Math.random() * 7) + 1

const rand01 = () => {
  let tmp = 4
  while (tmp === 4) {
    tmp = rand7()
  }
  return tmp < 4 ? 0 : 1
}
// 构造1-10等概率，可以先构造0-9等概率再加1,需要4个二进制位
const rand10 = () => {
  let tmp = Infinity
  while (tmp > 9) {
    tmp = (rand01() << 3) + (rand01() << 2) + (rand01() << 1) + rand01()
  }
  return tmp + 1
}

console.log(rand10())

// 此类题以后先转换成等概率返回0和1
