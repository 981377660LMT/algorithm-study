"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 两个账户都有一些共同的邮箱地址，则两个账户必定属于同一个人。
// 即如果两个同名的人邮箱有相同的就将他们合并
const accountMerge = (accounts) => { };
console.log(accountMerge([
    ['John', 'johnsmith@mail.com', 'john00@mail.com'],
    ['John', 'johnnybravo@mail.com'],
    ['John', 'johnsmith@mail.com', 'john_newyork@mail.com'],
    ['Mary', 'mary@mail.com'],
]));
