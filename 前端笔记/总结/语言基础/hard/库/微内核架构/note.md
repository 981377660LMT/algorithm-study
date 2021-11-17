1. 微内核架构的本质就是将可能需要不断变化的部分封装在插件中，从而达到快速灵活扩展的目的，而又不影响整体系统的稳定。
2. 西瓜视频播放器主张一切设计都是插件，小到一个播放按钮大到一项直播功能支持。

对于微内核的核心系统设计来说，它涉及三个关键技术：插件管理、插件连接和插件通信

1. 插件管理：核心系统需要知道当前有哪些插件可用，如何加载这些插件，什么时候加载插件。常见的实现方法是`插件注册表机制`。插件注册表含有每个插件模块的信息，包括它的名字、位置、加载时机（启动就加载，或是按需加载）等。

```JS
// packages/xgplayer/src/index.js
import Player from './player' // ①
import * as Controls from './control/*.js' // ②
import './style/index.scss' // ③
export default Player // ④

从 index.js 文件中，我们发现在第二行代码中使用了 import * as Controls from './control/*.js' 语句批量导入播放器的所有内置插件。该功能是借助 `babel-plugin-bulk-import` 这个插件来实现的。

```

2. 插件连接:插件连接是指插件如何连接到核心系统
   通常来说，核心系统必须指定插件和核心系统的连接规范，然后插件按照规范实现，核心系统按照规范加载即可。
   类似于 Vue 的插件

```JS
这里我们以简单的内置 loading 内置插件为例：
// packages/xgplayer/src/control/loading.js
import Player from '../player'

let loading = function () {
  let player = this;
  let util = Player.util;
  let container = util.createDom('xg-loading', `
    <svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewbox="0 0 100 100">
      <path d="M100,50A50,50,0,1,1,50,0"></path>
    </svg>
    `, {}, 'xgplayer-loading')
  player.root.appendChild(container)
}

Player.install('loading', loading)


// packages/xgplayer/src/player.js

class Player extends Proxy {
  static install (name, descriptor) {
    if (!Player.plugins) {
      Player.plugins = {}
    }
    Player.plugins[name] = descriptor
  }
}

当调用 Player.install 方法后，会把插件信息注册到 Player 类的 plugins 命名空间下。需要注意的是，`这里仅仅是完成插件的注册操作。在利用 Player 类创建播放器实例的时候，才会进行插件初始化操作`

在 Player 类构造函数中会调用 pluginsCall 方法来初始化插件

在西瓜视频播放器中，自定义插件只有两个步骤：
「1. 开发插件」
// pluginName.js
import Player from 'xgplayer';

let pluginName=function(player){
  // 插件逻辑
}

Player.install('pluginName',pluginName);

「2. 使用插件」
import Player from 'xgplayer';

let player = new Player({
  id: 'xg',
  url: '//abc.com/**/*.mp4'
})

```

**静态方法 install 插件 调用构造函数加载插件**

3. 插件通信
   「由于插件之间没有直接联系，通信必须通过核心系统，因此核心系统需要提供插件通信机制」。
   这种情况和计算机类似，计算机的 CPU、硬盘、内存、网卡是独立设计的配置，但计算机运行过程中，CPU 和内存、内存和硬盘肯定是有通信的，计算机通过主板上的总线提供了这些组件之间的通信功能。

   源码中是通过 player 实例提供的 `on、off 和 once` 三个方法来实现
   那么上述的三个方法来自哪里呢？通过阅读西瓜视频播放器的源码，我们发现上述方法是 Player 类通过继承 Proxy 类，在 Proxy 类中又通过构造继承的方式继承于来自 event-emitter 第三方库的 EventEmitter 类来实现的。

   在西瓜视频播放器初始化的时候，会通过调用 Video 元素的 addEventListener 方法来监听各种原生事件，在对应的事件处理函数中，会调用 emit 方法进行事件派发。
   在西瓜视频播放器销毁时，会调用 destroyFunc 方法，在该方法内部，`会继续调用 emit 方法来发射 destroy 事件`。之后，若其它插件有监听 destroy 事件，那么将会触发对应的事件处理函数，执行相应的清理工作。而对于`插件之间的通信，同样也可以借助 player 播放器对象上事件相关的 API 来实现`，这里就不再展开。
