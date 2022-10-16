// class Solution:
//     def sumOfNumberAndReverse(self, num: int) -> bool:
//         for x in range(10**5 + 1):
//             if x + int(str(x)[::-1]) == num:
//                 return True
//         return False
function sumOfNumberAndReverse(num: number): boolean {
  for (let x = 0; x <= 10 ** 5; x++) {
    if (x + parseInt(x.toString().split('').reverse().join(''), 10) === num) {
      return true
    }
  }
  return false
}
