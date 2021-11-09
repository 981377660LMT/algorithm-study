/**
 * @param {string} s
 * @return {boolean}
 */
var checkRecord = function (s: string): boolean {
  // check if there are more than 2 As and 3 continuous Ls
  // return !/^.*(A.*A|L{3,}).*$/.test(s)
  return !s.includes('LLL') && s.split('A').length < 3
}

console.log(checkRecord('PPALLP'))
// A少于2 且不存在连续三个或以上的L  返回true
