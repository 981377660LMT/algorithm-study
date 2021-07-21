"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// const reverseString = (s: string[]) => {
//   for (let index = 0; index < s.length; index++) {
//     s.splice(index, 0, s.pop()!)
//   }
//   return s
// }
const reverseString = (s) => {
    let l = 0;
    let r = s.length - 1;
    while (l < r) {
        ;
        [[s[l], s[r]]] = [[s[r], s[l]]];
        l++;
        r--;
    }
    return s;
};
console.log(reverseString(['h', 'e', 'l', 'l', 'o']));
