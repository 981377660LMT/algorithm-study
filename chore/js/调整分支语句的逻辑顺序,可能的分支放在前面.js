// 作者：在下康某
// 链接：https://www.zhihu.com/question/288096930/answer/2590083542
// 来源：知乎
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
// !调整分支语句的逻辑顺序,可能的分支放在前面
//
// if (cond1) {
//   do_something1();
// } else if (cond2) {
//   do_something2();
// } else if (cond3) {
//   do_something3();
// }
// 如果 cond1，cond2，cond3 不相关，且绝大多数情况下都是 cond2，可以改写成：if (cond2) {
//   do_something2();
// } else if (cond1) {
//   do_something1();
// } else if (cond3) {
//   do_something3();
// }
