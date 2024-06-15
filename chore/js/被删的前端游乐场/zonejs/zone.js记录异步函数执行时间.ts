// 那它到底是怎麼做到的，zone.js採用猴子補丁（Monkey-patched）的暴力方式將JavaScript中的非同步方法都包了一層，使得這些非同步方法都將運行在zone的上下文中。每一個非同步的任務在zone.js都被當做為一個Task，並在Task的基礎上zone.js為開發者提供了執行前後的方法（hook）包括：

// onZoneCreated：產生一個新的zone對象時的函數。 zone.fork也會產生一個繼承至基類zone的新zone，形成一個獨立的zone上下文；
// beforeTask：zone Task執行前的函數；
// afterTask：zone Task執行完成後的函數；
// onError：zone運行Task時候的異常函數；

import 'zone.js'

{
  function f1(): void {
    console.log('F1')
  }

  function f2(): void {
    console.log('F2')
  }

  function MyZone(): void {
    f1() // Sync
    setTimeout(f2, 5000) //Async
  }

  const zZone = Zone.current.fork({
    name: 'zZone',
    onHasTask(parentDelegate, current, target, hasTask) {
      console.timeLog('a')
    }
  })

  console.time('a')
  zZone.run(MyZone)
}

// 為什麼Angular要使用 Zone.js
// 用來追縱非同步事件
// Click Event
// Error Handling
// Debugging
// Websockset
// etc...

export {}
