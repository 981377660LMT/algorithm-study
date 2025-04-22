目录：
本课程会涉及现代游戏引擎的

- 体系架构
  游戏引擎的基础构建，如 MVVM
  - 游戏引擎分几层？
  - 看引擎代码先看 update 函数
- Rendering
  把东西放在 1/30s 内显示出来，渲染管线，各种算法如何拼接和组合
- Animation
  过渡怎么做、应变、一些列的动画如何组织起来，强调交互与玩法，让设计师能够理解
- Physics
  物理表达，刚体力学模拟运动，弹性力学，流体模拟
- Gameplay：前面就是世界的模拟器，而 gameplay 就是好玩的规则，规则如何让计算机理解，而且需要设计师去使用，所有的游戏就是规则体系；编程不一定是写代码，制作逻辑就算编程
- Misc. Systems：特效系统，寻路系统、相机系统
- Tool set：构建工具体系，如反射体系，在更新的时候，过去与未来的设计能够兼容
- Online Gaming：通过信息沟通，让每个平行宇宙发生的事情是一致的；异步同步算法、帧同步算法
- Advanced Technology：
  - Motion matching、PCG（自动生成）
  - 面向数据的编程、多线程的任务系统，程序在多核运行
  - lumen 光照技术、Nanite 面片技术

1. 游戏引擎导论 (Overview of Game Engine)
   游戏引擎是一个软件开发`框架`。
   现代游戏引擎是一个非常庞大的系统，它集合了整个计算机科学从底层硬件到应用端几乎全部的内容。某种意义上说，开发一个现代游戏引擎的难度不亚于开发一个现代操作系统。
   `Game Engine is Way Beyond Rendering`
   - 为什么要学习游戏引擎
     技术
   - 游戏发展史
     1. 早期的游戏并没有专门的引擎，开发者的目标是在有限的计算能力和存储空间中实现相对复杂的游戏玩法。公认的游戏引擎之父是 John Carmack，他开发的重返德军总部(Wolfenstein 3D, 1992)可以认为是第一款通过游戏引擎开发的游戏。
     2. **现代游戏引擎的出现则要归功于 GPU 的发展**。随着对图形渲染需求的增加，人们开发出了 GPU 这样的计算硬件`将图形运算从 CPU 逻辑运算分离出来`。通过对**高度并行化图形运算**的优化，人们第一次实现了全程的 3D 游戏。
     3. 除了游戏引擎外也有很多公司和团队专精某一个子系统的功能，包括物理仿真、声学仿真、渲染等等。
   - 游戏引擎的目标是在有限的计算资源和带宽下尽可能对现实世界进行**实时**模拟。`实时性`是游戏引擎和图形学其它研究的重要区别之一，大部分游戏特效的计算都要压缩在 1 ms 左右。
   - 游戏引擎和工业软件例如 CAD、建筑软件等有很多共通点，但是每个引擎都会根据自己的需求做定制 游戏引擎仍然是跑的最远的，因为在商业上被充分竞争。
2. 基础架构 1 – 引擎架构分层，整体 Pipeline (Layered Architecture of Game Engine)
   现代游戏引擎是一个非常复杂的系统并包含了海量的代码，但幸运的是游戏引擎一般会通过分层的方式将这些代码组织起来。因此了解游戏引擎的分层架构有助于对整个系统形成全面的认识。
   对游戏引擎进行分层的意义在于对不同类型的代码进行解耦，这样可以更好地管理这个系统的复杂度。

   1. 工具层：现代游戏引擎的界面，各种编辑器，直接和开发者进行交互
   2. 功能层：让这个世界看得见动起来
      `游戏引擎的核心是 tick 函数`

      - tickLogic：先把整个世界的物理规则算一遍，模拟出这个世界
      - tickRender：某个人看到的一副 2d 画面；会做裁剪、光照、阴影

   3. 资源层：各种数据，动画、模型、音乐，负责加载管理这些资源，为功能层提供弹药

      - 把资源的数据转换成引擎的高效数据：resource->assets；去掉无用信息
      - 数据之间的关联 reference 是最重要的
      - `guid：资产的全局唯一识别号`，相当于身份证
      - 需要一个实时的资产管理器，handle 系统
      - 资产会根据进度进行加载和卸载，guid 和 handle 就是解决这个问题
        延迟加载策略，比如材质从粗糙到细致的加载

   4. 核心层：功能层各个部分都会调用相同的很基础底层的代码，就像是工具箱、瑞士军刀

      - 数学库：数学效率
        - carmack’s 1/sqrt(x)
        - sse：cpu 并行运算向量
      - 数据结构：在核心层做一套自己的数据结构，没有内存碎片
      - 内存管理：开辟一大块内存自己管理，为追求最高的效率；cpu 的缓存越大，取出数据效率越高；
        内存管理三大步骤：把数据放在一起，尽可能地顺序访问数据，读写的时候尽可能一起去读写

   5. 平台层：最终发步到用户的设备上，平台不一样、输入不一样（键盘、手柄、体感）

      - 掩盖掉平台差异度
      - Render Hardware Interface(RHI)：重新定义一个 api，把各个硬件 sdk 封装起来

   6. 中间件和第三方库：会集成到上面的各个层中，有 sdk 可以直接编译进去；还有是工具，在引擎之外

   - Decoupling and Reducing Complexity
     底层与高层独立
     Lower layers are independent from upper layers.
     封装，高层不知道底层的实现
     Upper layers dont't know how lower layers are implemented.
   - Response for Evolving Demands
     底层的代码不要轻易改动，顶层可以；只允许上调用下
     Upper layers evolve fast, but lower layers are stable.
     ![alt text](image-1.png)

   1. 打开一个游戏引擎会直接看到各种类型的编辑器。这个直接和开发者进行交互的层称为工具层(tool layer)。
      工具层是为游戏开发者提供支持的一层。工具层允许游戏开发者直观地预览不同美术资源在游戏环境中的表现，对于游戏开发者和设计师有着重要的意义。
      DCC(Digital Content Creation)：别人开发的资产生产工具（比如 3dmax、maya、houdini）
      从 dcc 到 game engine 的流程叫做 `asset condition pipeline，有导入器、导出器`
      ![alt text](image.png)
   2. 在工具层下面包含功能层(function layer)，用来实现游戏的渲染、动画、交互等不同类型的功能。

      功能层的作用类似于一个时钟(`tick`)用来控制整个游戏世界的运行，在每个时钟周期内功能层需要完成整个游戏系统的全部运算
      **tickLogic(deltaTime) -> tickRender(deltaTime)**

   3. 在此基础上还需要资源层(resource layer)来管理各种各样的场景美术资源。
      游戏中的每一个对象都包含不同类型的资产(asset)，比如说几何模型、纹理、声音、动画等，同时每种资产都可能包括不同的数据格式。通过资源层我们将游戏对象的资产组织起来并通过全局唯一标识符(Globally Unique Identifier, GUID)进行识别和管理。

      - **资源层的核心在于管理不同资产的生命周期(life cycle)**。gc 很关键

   4. 再下一层是核心层(core layer)，它包括支持游戏渲染、动画、物理系统、内存管理等不同系统的核心代码。
      核心层数学库还大量使用了 `CPU 的 SIMD 指令`，这样可以极大地加速矩阵和向量的计算速度。
      **核心层最重要的功能是实现内存管理(memory management)，**

      **内存管理三条核心**：

      - put data together (把数据放在一起)
      - access data in order (按顺序访问数据)
      - allocate and de-allocate as a block (读写在一个块中)

   5. 最底层是平台层(platform layer)，一般包括各种图形 API、输入设备支持以及不同游戏平台的底层代码。
      平台层来提供对不同操作系统和硬件平台的支持
      Render Hardware Interface：重新定义一个 api，把各个硬件 sdk 封装起来
   6. 除此之外各种不同类型的中间件也会贯穿整个游戏引擎的不同层。

   - Takeaways

     1. 引擎是分层架构的：
        平台层、核心层、资源层、功能层、工具层
     2. 越底层的代码越稳定质量越高，越高层的代码设计得越开放越灵活，能适应不同的游戏
     3. 游戏引擎的核心是 tick 函数

3. 基础架构 2 – 数据组织和管理
   `组合`

   - 游戏世界是由各种各样的**游戏对象(game object, GO)**组成的
   - 现代游戏引擎中一般是通过**组件化(components)**的方式来描述一个对象(组件模式)
     以无人机为例，我们可以把任意形式的无人机拆分成 Transform、Motor、Model、AI 等组件，然后单独实现每个组件的功能。
   - 更新方法模式：两种 tick  
     object-based tick -> component-based tick
     component-based tick 相对要反直觉一点，但却有着更高的性能。component-based tick 相当于把系统的 tick()函数分解成流水线，这样可以提高系统的并行性并重复利用缓存。
   - 使用了 event 机制来实现 GO 之间的通信
     **很多 tick 是并行执行的，所以时序很重要，如果让 go 彼此互相通信，就会产生逻辑混乱，因此引入邮局(dag)；go 之间不能直接通信**
     **pretick 和 posttick 就是解决时序问题**
   - How to Manage Game Objects?
     稍微好一点的管理方式是把场景划分成一个均匀的网格，这样需要进行查询时只考虑相邻格子中的 GO 即可。`但如果 GO 在网格上不是均匀分布，这种管理方式的效率仍然很低。`
     合理的管理方法是使用带层级的网格，这种方式的本质是使用一个`四叉树`进行管理。这样进行查询时只需要对树进行遍历即可。
     除了四叉树之外，现代游戏引擎进行场景管理时还会使用 `BVH、BSP、octree、scene graph 等不同类型的数据结构。`
     ![alt text](image-2.png)
   - FAQ
     - 如果 tick 时间过长怎么办：一般会用 deferred tick，把爆炸的计算分散到各个帧里面算，做延迟 tick(时间分片)。
     - tick 时渲染线程和逻辑线程怎么同步：tickrender 主要做数据准备，ticklogic 会比 tickrender 早一些。

4. 渲染系统 1 – 渲染数据组织 (Basic of Game Rendering)
   ![outline](image-3.png)

   渲染研究更关注于算法理论的正确性而对于实时性没有太多的要求，而在游戏渲染中实时性则至关重要。对`实时渲染`的关注构成了游戏渲染和渲染理论之间的主要差别。

   - realtime -> 30fps
   - interactive -> 100fps
   - offline rendering
   - out-of-core rendering (电影 cg 渲染，几天渲染一帧)

   渲染的四个难点

   1. 不同类型的渲染对象(场景复杂)
   2. 要考虑渲染过程的硬件实现
   3. 人们对于游戏画质和帧率的要求逐渐提高

      - 帧率稳定，30fps 算 realtime，10fps 算可交互
      - 1080P、4K 甚至 8K 分辨率

   4. 渲染系统只能占掉 10-20%的计算资源

   - Rendering Pipeline and Data
     目前游戏引擎渲染的主流方法仍是**基于光栅化(rasterization)的渲染管线。**
     ![alt text](image-6.png)
     投影、光栅化、着色、纹理采样：
     首先我们需要把场景中的物体投影到 NDC 上，然后分别计算平面上每一个像素对应的渲染对象。
     接下来对于每一个像素需要调用相应的 shader 计算像素的颜色。
     在调用 shader 时往往还需要通过纹理采样的方法进行反走样等处理。
     > mipmap：纹理的不同分辨率，用来解决纹理采样的问题
   - Understand the Hardware
     渲染计算的特点是有大量的像素需要进行计算，而像素之间的计算则往往是相互独立的。因此人们设计出了 `GPU(显卡) 来执行图形渲染计算`，这样还解放了 CPU 的计算资源。可以说现代 GPU 的发展也推动了整个游戏行业渲染技术的进步。
     SIMD 是指在运行程序时可以`把一条指令同时执行在不同的数据上`，目前现代 CPU 对于 SIMD 有着很好的支持，这种技术在高性能计算等领域中有着广泛的应用；而 SIMT 则是把`同一条指令分配到大量的计算核心上同时执行`，现代 GPU 的计算过程更类似于 SIMT。

     > simd(single instruction Multiple Data)：单指令、多数据的数据运算，就是四维向量，也就是说一个指令同时完成四个加减法
     > simt(single instruction Multiple Threads)：如果计算内核很小，但有很多个，一个指令让很`多核`去计算，极可能使用同样的代码，每个核访问自己的数据，所以跑的非常快

     GPU 费米架构：
     ![Intel显卡费米架构](image-11.png)
     在现代 GPU 架构中有着大量重复的内核，每一组内核称为一个 GPC(Graphics Processing Cluster)。在每个 GPC 内部存放着大量的 SM(Stream Multiprocessor)，而每个 SM 中还有着大量的 CUDA 核心用来执行数学运算，当 SM 接收到指令进行计算时会把运算分配给 CUDA 核心进行并行计算。同时 GPU 上还有 share memory 用来实现 GPU 上不同核心以及 GPU 和 CPU 之间的通信。

     `GPU 和 CPU 之间通信的代价是非常大的，因此在渲染系统中会尽量把数据通信设计为单向的。这样 GPU 只需要读取 CPU 发送的数据而无需反向传输渲染的结果。`
     ![尽可能单向传输数据](image-12.png)

     Cache：内存固然快，但依然比缓存慢百倍 (数据要放在一起)

   - Renderable
     在进行渲染时我们只需要考虑那些需要进行渲染的 GO，它们称为可渲染对象(renderable)。
     一般来说我们可以把整个可渲染对象拆分成若干个 block，每个 block 有着自身的网格、材质等渲染信息。
     对于网格数据，我们需要存储网格上所有的顶点坐标以及每个面包含节点的编号。同时我们往往还需要为每个顶点单独存储一个法向来`处理曲面发生突变`的情况。

     > 存储方式：顶点缓冲区（vertex buffer）和索引缓冲区（index buffer）
     > https://doodlewind.github.io/learn-wgpu-cn/beginner/tutorial4-buffer/

     ![alt text](image-13.png)
     对于材质(matrials)数据，我们需要定义常见材质的渲染模型。在现代游戏引擎中往往还会集成大量的 PBR 材质以渲染出更加逼真的图像。
     除此之外，我们还需要考虑材质的纹理(texture)。纹理对于材质的定义以及最终渲染呈现的效果起着至关重要的作用。
     当然我们还需要考虑 shader，在进行渲染时需要把编译好的 shader 连同数据一起提交的 GPU 上进行计算。

   - Render Objects in Engine
     接下来我们就可以对 GO 进行渲染了。
     将顶点数据、材质、纹理交给 gpu，再让 gpu 按 shader 执行，就会渲染出东西来。
     实际工程中我们往往需要把一个完整的网格拆分成不同的`submesh，每个submesh有着自己的材质和纹理而整个网格共享一套顶点和面片信息(instance pool，享元模式)`。这样利用 submesh 的概念就可以绘制出更加逼真的图像。
     为了更高效地利用 GPU，我们还可以把场景中的`submesh按照材质进行排序。这样可以保证渲染时具有相同材质的submesh会放在一起进行绘制，从而降低GPU切换资源的开销`
     在很多游戏场景中还存在着大量相似甚至是完全相同的 GO。对于这种情况可以通过`GPU batch rendering`的方法把这些 GO 组织在一起，然后把同一 batch 中的对象一次性交给 gpu 绘制出来，进一步提升场景渲染的效率。
   - Visibility Culling
     `可见性剔除(visibility culling)`是渲染系统中非常实用的技术，它的思想是在送入渲染管线前首先判断场景中的每个可渲染对象是否在相机视野中，然后只对视野范围内的对象进行渲染。
     `visibility culling的核心是把可渲染对象使用bounding box(包围盒)进行表示，然后通过bounding box来迅速判断物体是否在视锥范围内。`我们通过 bounding box 将场景中的物体组织起来，这样在渲染时只需要通过对它们进行遍历就可以快速地实现 visibility culling。

     ![bounding box](image-7.png)

     在现代游戏引擎中，**BVH(Bounding Volume Hierarchy) 是应用最为广泛的 bounding box。BVH 的一大特点是它可以在场景中物体发生运动时通过对节点的操作来动态地修改树的结构**，这样无需每次都重新建树从而大大提高了计算效率。

     在游戏设计中，**PVS(potential visibility set)**是一种非常实用的技术。`它的思想是把整个场景划分为若干个相对独立的区域，不同区域之间通过 portal 进行连接。`先用 BSP 将房间分成小格子，每个格子通过 Portal 链接，在其中一个格子时，只看到部分的格子，从而减少绘制。
     ![每个房间能看见哪几个房间?](image-8.png)

     利用现代 GPU 的强大并行计算能力我们可以**显卡来做裁剪，生成遮挡物的深度图，然后选择是否丢掉，这样的技术称为 GPU based culling。**

   - Texture Compression
     纹理压缩

     我们在前面的章节介绍过纹理对于渲染出逼真的物体起着重要的作用。通常情况下纹理会通过一张二维贴图进行表示，并且在`计算机中使用 JPG 或是 PNG 这样的压缩格式进行存储`。而在游戏引擎中则无法使用这些常用的图像压缩格式，这主要是因为 `JPG 这样的压缩算法不支持快速的随机图像坐标访问(Random Access)，而且它们往往具有过大的计算复杂度无法进行实时的压缩与解压。`

     在渲染系统中最常用的纹理压缩算法是`block compression`，它的思想是统计每个 4×4 区域内纹理图像最大和最小值然后通过插值的方法进行查询。
     PC 端：BC7 算法
     `移动端：ASTC 算法`

   - Authoring Tools of Modeling

     游戏中的模型是怎么获得的呢？

     1. 几何建模：最经典的建模方法是使用 3ds Max、Maya、blender 等建模软件来绘制 3D 模型。
     2. 雕刻：近几年基于雕刻的建模软件也获得了非常多的应用。
     3. 扫描：随着人工智能和三维重建技术的发展，我们甚至可以从实物通过扫描的方法来重建出非常精细的网格。
     4. 程序化生成：除此之外，还有一些自动化建模工具来自动生成地形等场景的网格。

     ![alt text](image-10.png)

   - Cluster-Based Mesh Pipeline
     随着现代 GPU 计算能力的提高以及人们对于画质需求的不断增长，在 3A 大作中的模型往往都具有`百万级甚至千万级的网格`
     为了渲染出具有如此高精度的网格就需要使用 **mesh shader** 相关的技术。**mesh shader 的核心思想是把网格上的一小块区域视为一个 meshlet，每个 meshlet 都具有固定数量的三角形。**
     mesh shader 可以生成几乎无限的细节，而且可以根据相机和物体的相对位置关系动态地调整网格的精度。
     虚幻 5 中的 Nanite 技术可以认为是更加成熟的 mesh shader。

   - Take Away
     1. 游戏引擎的渲染深度依赖于硬件架构设计(gpu)
     2. submesh 来提高渲染效率
     3. visibility culling 尽可能少地渲染(do nothing)
     4. 尽量用 GPU 做计算，CPU 只负责数据准备 (gpu driven)

   Q&A:

   - 引擎有必要自写渲染管线吗？
     最好不造轮子

5. 渲染系统 2 – 光照 (Materials、Shaders and Lighting)
   渲染是研究光与材质相互作用的学科，因此本节课从光线、材质以及 shader 三个方面介绍现代游戏引擎中各种经典实时算法的原理。

   - The Rendering Equation
     渲染的本质是求解渲染方程(the rendering equation)，它由 James Kajiya 于 1986 年提出。
     渲染的难点可以分为一下三部分：
     如何计算入射光线、如何考虑材质以及如何实现全局光照。
   - Starting from Simple

     1. 环境光+主光+环境光贴图
     2. Blinn-Phong 模型是最简单的光照模型，它包括了漫反射、高光反射和环境光反射三种光照效果。问题是能量不守恒。
     3. 对于阴影问题，现代游戏引擎的主流方法是 shadow map。
        shadow map 的处理流程是在光源位置设置一个新的相机并渲染出一张`深度图`，然后在实际相机进行渲染时对每个点检测它到光源处的深度。

   - Pre-computed Global Illumination
     全局光照可以显著地提升画面的渲染效果。

     - 球面谐波函数(spherical harmonics, SH)是实时渲染中表示环境光照的经典方法
     - lightmap ：我们可以将场景中每个点的光照离线烘焙到一张纹理图上，然后在渲染时读取纹理值来获得 SH 表达的环境光照。计算 lightmap 是非常耗时的，但通过 lightmap 可以实现非常逼真的场景效果，而且在实际渲染时 lightmap 可以实现场景的实时渲染。

   - Physical-Based Material

   - Image-Based Lighting

   - Classic Shadow Solution

   - Moving Wave of High Quality
     随着各种 shader 模型的提出以及硬件计算性能的进步，上面介绍的实时渲染算法已经不能完全满足人们对画质的需求。
     `实时光线追踪(real-time ray tracing)`就是一个很好的案例。随着显卡性能的提升我们可以把光线追踪算法应用在实时渲染中从而获得更加真实的光照和反射效果。
     在虚幻 5 引擎中还使用了 `virtual shadow map` 来生成更加逼真的阴影。`(原理类似雪碧图)`
   - Shader Management
     本节课最后讨论了游戏引擎中的 shader 管理问题。在现代 3A 游戏中每一帧的画面上可能都有上千个 shader 在运行。
     这些大量的 shader 一方面来自于美术对场景和角色的设计，另一方面不同材质在不同光照条件下的反应也使得程序员需要将不同情况下的 shader 组合到一起，并通过宏的方式让程序自行选择需要执行的代码。

6. 渲染系统 3 – 天空，地形，后处理等 (Specail Rendering)

   - Landscape

     1. 地形的几何表示

        - 表示地形最简单的方法是使用高度场(heightfield)。我们可以把地形看做是平面上具有不同高度的函数，然后通过在平面进行均匀采样来近似它。这种方法在遥感等领域仍然有着很多的应用。
          高度场的缺陷在于当我们需要表示大规模的地形或者需要更精细的地形时所需的采样点数会成倍的增长。

        - Adaptive Mesh Tessellation
          在游戏引擎中由于玩家观察的`视野(field of view, FOV)`是有限的，实际上我们不需要对所有的网格进行加密采样，只需关注视野中的地形即可。在这种思想下人们提出了两条加密采样原则：
          1. 根据距离和视野来调整网格的疏密，对于不在视野范围内或是距离观察点比较遥远位置的地形无需使用加密的网格
          2. 近处地形的误差尽可能小而远处的误差可以大一些
             ![alt text](image-14.png)
        - Triangle-based subdivision
          对三角网格进行加密可以通过三角网剖分算法来实现。对于均匀分布的网格，其中每个三角形都是等腰直角三角形。
        - **QuadTree-Based subdivision**
          `在游戏行业中更常用的高度场表达方式是使用四叉树来表达地形。`这种方法更符合人的直觉，同时也可以直接使用纹理的存储方式来存储这种四叉树的结构。
          ![alt text](image-15.png)

     2. 地形的纹理
        在现代游戏引擎中大量使用了虚拟纹理(virtual texture)的技术来提高渲染性能。

        - Camera-Relative Rendering：
          当渲染物体与相机的距离达到一定程度时就需要考虑浮点数的计算精度问题，如果不进行处理会导致严重的抖动和穿模现象。
          想要缓解这种问题可以将相机设置为世界坐标的中心，这样的处理方法称为 camera-relative rendering。

   - Sky and Atmosphere

7. 渲染系统 4 – 渲染管线 (pipeline)

- 走样（Aliasing）分为好几种，一个是锯齿，另一个是高光的闪烁，还有就是 texture 采样精度不足导致的纹理扭曲。所有的 Aliasing 都是采样率不足导致的。应对方法是`超采样`。

8. 动画系统 1 – 骨骼动画

- Introduction
  Challenges

  - 可交互性和动态变化的动画：
    在游戏中不能预设玩家的行为
    游戏中的动画要和很多 gameplay 互动
    受制于周围的环境
  - 实时：
    每一帧都要计算
    动画数据很大
  - 真实度：
    更生动
    表情

- 2D Animation Techniques in Games

  - live2D：通过把一个人物拆分成多个组件，例如头发，眼睛，衣服等，然后把所有的图元生成一些控制网格，通过编辑这些控制网格来编辑 keyframe 动画。

- 3D Animation Techniques in Games

  - 自由度(degrees of freedom, DoF)
    对于刚体而言描述它的运动需要 3 个平动和 3 个旋转一共 6 个自由度
  - Rigid Hierarchical Animation：层次结构刚体动画
    问题是骨骼转的时候 mesh 彼此会穿插
  - Per-vertex Animation：顶点动画

- Skinned Animation Implementation

- Math of 3D Rotation

  1. 欧拉角：三维旋转矩阵
     ![alt text](image-16.png)
     欧拉角的主要缺陷如下：

     万向锁(gimbal lock)及相应的自由度退化问题；
     很难对欧拉角进行插值；
     很难通过欧拉角对旋转进行叠加；
     很难描述绕 x,y,z 轴之外其它轴的旋转。
     由于这些缺陷的存在，游戏引擎中几乎不会直接使用欧拉角来表达物体的旋转。

  2. **四元数(Quaternion)**
     在游戏引擎中更常用的旋转表达方式是四元数(quaternion)
     ![alt text](image-17.png)
     任意的三维旋转可以通过一个单位四元数来表示。当我们需要对点 v 进行旋转时，只需要先把 v 转换成一个纯四元数，然后再`按照四元数乘法进行变换，最后取出虚部作为旋转后的坐标即可：`

- Joint Pose

- Animation Compression

- Animation DCC Process

9. 动画系统 2 – 高级动画技术：动画树、IK 和表情动画

- Animation Blending
  在实际游戏中我们还需要将不同类型的动画混合起来以实现更加自然的运动效果
- Blend Space

- Action State Machine (ASM)
  动作状态机(action state machine, ASM)。
- Animation Blend Tree
  动画树(animation blend tree)
- Inverse Kinematics (IK)
  反向运动学(inverse kinematics, IK)
- Facial Animation

- Animation Retargeting

10. 物理系统 1 – 碰撞和刚体

游戏中的物体分为几类：

- 静态 Actor：无法移动的例如墙，地板等
- 动态 actor：一些动态的物体，可以符合物理碰撞规律的，例如弹珠撞到了一个石头，石头会动。
- trigger：当任何一个 acter 前面，其他的物体会作出相应的反应，例如自动门。
- 反物理规律的(kinematics)：例如人在推箱子，但是推力没有设置好，箱子一下子就飞出去了。或者是有一些游戏机关，地板在不断的上下移动。让玩家上去。

- Collision Detection(碰撞检测)：
  - Broad Phase
    只利用物体的 bounding box 来快速筛选出`可能发生`碰撞的物体
    使用 BVH 空间划分快速检测碰撞 使用 sort and sweep 来做碰撞检测，该方法优势在于只要把物体都提前排序以后，只移动少部分物体，效率会非常高。
    ![alt text](image-18.png)
    ![Alt text](image-4.png)
  - Narrow Phase
    筛选出可能发生碰撞的物体后就需要对它们进行`实际的碰撞检测`，这个阶段称为 narrow phase。除了进一步判断刚体是否相交外，在 narrow phase 中一般还需要去计算交点、相交深度以及方向等信息。
    目前在 narrow phase 中一般会使用相交测试、Minkowski 距离以及分离轴等方法。
    - 简单物体求交： 圆、胶囊等求交相对简单
    - 凸包的求交使用 Minkowski 和 与 Minkowski difference 来判断。Minkowski difference 肯定会经过圆点，使用 GJK 算法找到圆点 另一种算法是通过判断是否存在一条边能把两个物体分开
    - 分离轴定理(separating axis theorem, SAT)同样是一种计算凸多边形交的算法
      它的思想是`平面上任意两个互不相交的图形我们必然可以找到一条直线将它们分隔在两端`。对于凸多边形还可以进一步证明必然存在以多边形顶点定义的直线来实现这样的分隔，因此判断凸多边形相交就等价于寻找这样的分隔直线。
      ![alt text](image-19.png)
- Collision Resolution(碰撞解决)：
  完成碰撞检测后就需要对发生碰撞的刚体进行处理，使它们`相互分开`。目前刚体的碰撞主要有三种处理思路，分别是 penalty force、velocity constraints 以及 position constraints

  - velocity constraints
    目前物理引擎中主流的刚体碰撞处理算法`是基于 Lagrangian 力学的求解方法，它会把刚体之间的碰撞和接触转换为系统的约束，然后求解约束优化问题`。

- Scene Query
  对场景中的物体进行一些查询，这些查询操作也需要物理引擎的支持
  - Raycast
    查询与射线相交的最近物体
    实际上在光线追踪中就大量使用了 raycast 的相关操作，而在物理引擎中 raycast 也有大量的应用，比如说`子弹击中目标`就是使用 raycast 来实现的
  - Sweep
  - Overlap
  - Collision Group
- Efficiency, Accuracy, and Determinism

  - Simulation Optimization
    分块：把场景中的物体划分为若干个 island，当 island 内没有外力作用时就对它们进行休眠
  - Continuous Collidion Detection(CCD)
    连续碰撞检测
    当物体运动的速度过快时可能会出现一个物体之间穿过另一个物体的现象(tunneling)，此时可以使用 CCD 的相关方法来进行处理。

    解决穿墙问题：

    - 把墙做厚一点
    - CCD 方法：做一个保守估计，即物体和环境碰撞的一个安全距离是多少，在安全距离外可以任意移动，`在安全距离内就会把 substep 调密`，做更精细的检测

  - Determinism Simulation
    在进行物理仿真时还需要考虑仿真结果的确定性。也就是说对物理的模拟，相同的操作结果相同（为了使得联网游戏中每个玩家看到的世界相同）

11. 物理系统 2 – 布料模拟

- Character Controller
  角色控制器是一个反物理的系统，摩擦力几乎是无限大的
  - 在构建角色控制器时一般会使用`简化后的形状来包裹角色`，这样便于处理各种场景之间的互动。
  - 玩家控制的角色撞到了墙壁上 => 碰到环境比如墙体不能前进时，会滑一下，修改角色的运动方向
  - 上台阶的时候，把 capsule 往上添加 offset
  - 在不同姿态 controller 体积要变换
  - 当 controller 和一个物体站在一起时，会在逻辑上将二者绑定起来
- Ragdoll
  布娃娃系统
- Cloth
  在布料仿真中往往还会为网格上的每个顶点赋予一定位移的约束，从而获得更符合人直觉的仿真结果。
- Destruction
  破坏系统
  - Chunk Hierarchy
    把碎片组成一个树状结构
  - Connectivity Graph
    使用一张图来表示不同碎片之间的连接关系,当冲击大于边上的值时就会发生物体的破碎。
  - Damage Calculation
  - Fracturing with Voronoi Diagram
    在物理引擎中一般会使用 Voronoi 图(Voronoi diagram)这样的技术来对原始的物体区域进行划分
  - Destruction in Physics System
  - Issues with Destruction
    当碎块碎了之后会执行很多回调函数，比如出发音效、粒子效果、navigation 更新
  - Popular Destruction Implementations
    目前很多商业引擎都有现成的破坏系统。
- Vehicle
  载具系统
- Advanced: PBD/XPBD

12. 粒子和声效系统
13. 工具链 1 – 基础框架
    ![Alt text](image-5.png)

    - Tool Chain
      工具链(tool chain)是沟通游戏引擎用户以及更底层 run time(渲染系统、物理引擎、网络通信等)之间的桥梁。对于商业级游戏引擎来说，工具链的工程量往往要比 runtime 大得多。
      ![alt text](image-20.png)
    - Complicated Tool

      - GUI

        1. immediate mode
           用户的操作会直接调用 GUI 模块进行绘制，让用户立刻看到操作后的效果
        2. retained mode (主流)
           用户的操作不会直接进行绘制，而是会把用户`提交的指令先存储到一个 buffer 中`，然后在引擎的绘制系统中再进行绘制

        ![alt text](image-21.png)

      - Design Pattern
        GUI 设计中常用的设计模式
        1. MVC
           `M 和 V 初步解耦，避免直接V操作M，但M可以直接作用于V`
           ![alt text](image-22.png)
           MVC 是经典的人机交互设计模式。MVC 的思想是把用户(user)、视图(view)和模型(model)进行分离，当`用户想要修改视图时只能通过控制器(controller)进行操作并由控制器转发给模型`，从而避免用户直接操作数据产生各种冲突。
        2. MVP
           `M和V解耦的更彻底,View 完全不知道 Model 的存在`
           ![alt text](image-23.png)
           V 和 M 之间的通信则通过展示者(presenter)来实现
           `presenter 容易臃肿，因为既要 query view 的 api，又要 query model 的 api`
           多维表格的 **View = connect(Component, Presenter)**
        3. **MVVM(游戏引擎中大量使用的 UI 设计模式)**
           微软提出，P 那层变成了 ViewModel，**变成了 bounding 机制**
           View 绑定到 ViewModel 上，ViewModel 于 Model 双向通信
           ![alt text](image-24.png)
           缺点：难以 debug
      - Load and Save

        - Serialization and Deserialization
          使用序列化(serialization)的技术来将各种不同的数据结构或是 GO 转换成二进制格式，而当需要加载数据时则需要通过反序列化(deserialization)从二进制格式恢复原始的数据
          最简单的序列化方法是把数据打包成 text 文件。
          目前常用的 text 文件格式包括 `txt、json、yaml、xml` 等。

          > text 可以作为 debug mode

          ![alt text](image-27.png)

          text 文件可以方便开发人员理解存储数据的内容，但计算机对于文本的读取和处理往往是比较低效的。当需要序列化的数据不断增长时就需要使用更加高效的存储格式，通常情况下我们会使用二进制格式来对数据进行存储。

          ![二进制](image-25.png)
          ![FBX Binary](image-26.png)
          和 text 文件相比，二进制文件往往只占用非常小的存储空间，而且对数据进行读取也要高效得多。因此在现代游戏引擎中一般都会使用二进制文件来进行数据的保存和加载。

      - Asset Reference(资产引用，享元)
        ![alt text](image-28.png)
        在很多情况下游戏的资产是重复的，此时`为每一个实例单独进行保存就会浪费系统的资源。`
        因此，在现代游戏引擎中会使用资产引用(asset reference)的方式来管理各种重复的资产。实际上资产的引用和去重是游戏引擎工具链最重要的底层逻辑之一。
        在调整和修改数据时直接进行复制很可能会破坏 GO 之间的关联而且容易造成数据的冗余，`因此在现代游戏引擎中对于数据引入了继承(inheritance)的概念。`数据之间的继承可以很方便地派生出更多更复杂的游戏对象，从而方便设计师和艺术家实现不同的效果。

        > <<Game Design Patterns>> 里的类型对象(Type Object)

    - How to Load Asset

      - Parsing
        在上一节我们主要是考虑如何对数据进行保存，而游戏引擎中工具链的一大难点在于如何加载不同的资产，即反序列化的过程。反序列化的过程可以理解为对文件进行解析(parsing)，文件中的不同字段往往有着不同的关键字以及域。我们需要对整个文件进行扫描来获得整个文件的结构。
        对文件完成解析后可以得到一棵由< key-value >对组成的树来表达不同类型的数据。
        ![alt text](image-29.png)
      - Endianness (大端小端)
        ![alt text](image-30.png)
        甜豆浆还是咸豆浆，没什么意思
        **对二进制文件进行反序列化和解析时需要额外注意 endianness 的问题。**在不同的硬件和操作系统上同样的二进制文件可能会被解析为不同的数据，这对于跨平台的应用需要额外注意。
      - Version Compatibility
        资产的兼容性问题
        在版本更迭中最常见的情况是数据的域发生了修改，新版本的数据定义可能会添加或删去老版本定义的域。
        ![Add or Remove Field](image-31.png)

        为了处理这种问题可以`手动为数据添加版本号`，在加载数据时根据版本号来控制加载过程。
        ![alt text](image-32.png)

        更好的处理方法是使用 guid 来进行管理。如 Google 就提出了使用 protocol buffer 来`为每一个域赋予一个 uid，在进行反序列化时只需要对域的 uid 进行比较即可。`
        ![alt text](image-33.png)

    - How to Make Robust Tools
      鲁棒性最基本的要求是`允许程序从崩溃中进行恢复，从而还原初始的开发状态。`
      为了实现这样的功能我们需要**将用户所有的行为抽象为原子化的命令(command)，通过命令的序列来表示整个开发的过程。**
      ![Command Pattern](image-34.png)

      **Undo/Redo**

      对 command 类进行抽象时需要为`每一个 command 实例赋予单调的 UID` 从而保证顺序的正确性，同时每一个 command 定义都需要实现 Invoke()和 Revoke() 方法表示执行命令以及恢复到执行命令前的状态。除此之外还需要实现 Serialize() 和 Deserialize() 方法来控制生成的数据序列化以及反序列化过程。

      整个 command 系统可以划分为三种不同类型的指令，包括**添加数据、删除数据以及更新数据。**实际上几乎所有的 command 都可以视为这三种基本指令的组合。
      ![alt text](image-35.png)

    - How to Make a Tool Chain
      ![alt text](image-36.png)
      而对于工具链来说，一个基本要求是要保证不同工具之间的沟通以及整个系统的可拓展性。我们不希望每个工具程序都使用单独的一套数据定义方式，这会导致整个工具链系统过于庞大而且难以进行维护。

      **关键是：同构**

      因此我们需要**去寻找不同工具中的一些共性，并把这些共同的数据封装为基本的单元**。利用对这些基本单元的组合来描述更加复杂的数据结构。

      - schema (DSL)
        schema 是一种对数据进行描述的结构，它描述了具体的数据结构是由哪些基本单元构成的。在工具链系统中所有流动的数据都要通过 schema 来进行描述，从而保证不同的程序都可以对数据进行解读。
        如何处理工具链中各个资源格式不同的问题： 使用 schema。将所有复杂的数据都拆分成一些“原子数据”，`schema 更像是一个分子式，是一个描述物体的格式`，schema 通常是一个 xml，而且要有继承关系，例如军人的 schema 可以继承自人的 schema。同时还需要能够相互 reference 数据。能够把数据关联在一起。

        `目前游戏引擎中的schema系统主要有两种实现方式，其一是单独实现schema的定义，而另一种则是使用高级语言进行定义。`

      - Three Views for Engine Data
        数据结构的三种视图
        - Runtime View (运行)
          在 runtime 中一般会以运行和计算效率为第一要务。
          ![alt text](image-37.png)
        - Storage View (存储)
          而在进行存储时则要游戏考虑数据的读写速度和空间需求。
          ![alt text](image-38.png)
        - Tool View (debug)
          而在面向开发者的工具程序中需要根据不同使用者的背景和需求来设计不同的数据表现形式。
          ![alt text](image-39.png)

    - What You See is What You Get
      所见即所得(what you see is what you get, WYSIWYG)是我们设计构建整个工具链系统的核心精神，它的目标是保证设计师和艺术家在工具链中的设计结果能完美地重现在实际的游戏场景中。

      - Play in Editor
        在编辑器中进行游玩时同样有两种实现方式
        1. `直接在编辑器中进行游戏`
           直接在编辑器中进行游戏可以无缝地对当前游戏场景进行编辑，但需要注意在进行编辑时不要污染游戏场景中的数据。
        2. `基于编辑器当前的状态生成一个新的游戏窗口(沙盒)进行游戏`
           这种设计方式可以保证编辑器中的数据与实际游戏中的数据保持相互独立，避免出现数据污染的情况。在大型游戏引擎的开发中一般会使用这种模式。

    - Plugin (可扩展性)
      现代游戏引擎的本质是一个开发平台

      ![api](image-40.png)
      Editor Call Plugin Interface
      Plugin Call Editor Api

      ![alt text](image-41.png)
      PluginManager、PluginInterface、API

      ![alt text](image-42.png)
      ![alt text](image-43.png)

      ![alt text](image-44.png)

14. 工具链 2 – Applications & Advanced Topic

    - Glance of Game Production
      设计师和艺术家往往需要使用大量的`第三方工具来辅助进行游戏角色和场景的建模`
      因此对于工具链来说，它需要实现和不同开发工具的通信也要考虑不同用户的使用需求。除此之外 WYSIWYG 原则又要求开发者在引擎中的体验必须和实际游戏时完全一致的，这对工具链的设计提出了更大的挑战。
    - World Editor
      世界编辑器(world editor)是整合了游戏引擎中几乎所有功能的平台。
      ![alt text](image-45.png)
      - Editor Viewport
        最重要
        注意引擎有一些 editor-only 的代码，不要带在 release 版本中
        ![alt text](image-47.png)
        兼容多个 view
      - Editable Object
        在同一个游戏场景中可能有成千上万个对象，因此我们需要一些高效的管理工具。`如何管理？`
        在编辑器中往往会使用`树状结构或是分组`的方式来管理场景中的对象，有时还会根据对象自身的特点设计相应的管理工具。
        ![alt text](image-48.png)
        当开发者选中某个对象时需要使用 schema 来获取该对象自身的信息。
      - Content Browser
        管理开发过程中设计的各种美术、场景资产
        ![alt text](image-49.png)
      - Editing
        1. 通过鼠标选取物体
           **在渲染流程中添加一个额外的选取帧，为图像上每一个像素赋予一个物体编号，这样使用鼠标进行选取时只需根据物体编号进行查询即可。** (code-inspector-plugin)
           还有一种方法是 raycasting。利用鼠标的位置来发射光线并通过与物体 bounding box 求交来来选择物体。这种实现的缺陷在于当物体比较复杂时 bounding box 是不能完全反应物体的几何形状的，此时使用 ray casting 的效率可能会比较低。
        2. 几何变换，包括平移、旋转和缩放等
           Object Transform Editing
        3. 高度笔刷
        4. instance brush (在设计好的地形上添加各种装饰件)
      - 规则系统(rule system)
        ![alt text](image-50.png)
    - Editor Plugin Architecture
      ![alt text](image-46.png)
      **对插件种类进行划分**
      我们可以`按照游戏对象的种类(网格、粒子、动画等)对插件进行分类，也可以根据对象的内容(NPC、建筑、人群等)进行分类。`因此现代游戏引擎的编辑器往往需要支持这种矩阵式的分类方法，允许用户根据喜好来选择和定制插件。
      Double Dispatch

      在对插件系统进行整合时还需要考虑不同插件之间的版本问题。不同的版本之间可以按照覆盖(covered)或是分布式(distributed)的方式进行协作。
      ![alt text](image-51.png)

      **实用主义，而不是追求优雅。任何能满足需求的架构都是好架构。**

      随着编辑器以及各种插件之间版本的迭代，插件系统一般还需要考虑版本控制的问题。

    - Design Narrative Tools
      除了游戏资产的设计外，**叙事(story telling)**在整个游戏开发流程中同样是非常重要的一环。叙事可以看做一个线性的过程，相关的游戏资产需要在一个`时间轴`上按照顺序进行调度。
      ![alt text](image-52.png)
      在虚幻引擎中使用了 sequencer 来跟踪游戏对象及其属性在时间轴上的变化。当我们把不同的对象利用 sequencer 在时间轴上组织起来就实现了简单的叙事。
    - Reflection and Gameplay
      通过反射我们可以让游戏引擎在运行阶段获取操作对象具有的各种属性。
      **通过宏来添加反射控制**
      在小引擎中还使用了`代码渲染(code rendering)`的技术来自动生成相关的代码。`(模板引擎)`，例如 Mustache。
      能自动化生成的代码，尽量自动化。
    - Collaborative Editing (下一代游戏引擎的发展方向)

      - 在虚幻引擎中提出了 `OFPA` 的策略来生成大量的文件来进行管理。
      - 协同编辑最终的状态是，每一个人可以实时看到别人的编辑结果，这其实是一个网络同步问题。这需要将对所有命令原子化。

      锁
      如何处理 undo/redo

15. GamePlay 1 – 玩法系统基础机制

- 玩法系统的核心是事件机制(event mechanism)
  pub-sub：

  - 事件定义(event definition)
    在现代游戏引擎中会使用反射和代码渲染的方式来**允许设计师自定义事件类型和相关的数据。**
  - 注册回调(callback registration)
    Invoke(激活)
    回调函数的一大特点在于它的注册和执行往往是分开的。这一特点可能会导致**回调函数调用时相关的一些 GO 可能已经结束了生命周期，因此回调函数的安全性是整个事件系统的一大难题。**
    解法 1：
    为了处理这样的问题我们可以使用强引用(strong reference)这样的机制来锁定相关 GO 的生命周期。强引用会保证所有和回调函数`相关的资源在回调函数调用前不会被回收，从而确保系统的安全。`
    解法 2：
    需要`Empty Service`
  - 事件派送(event dispatching)
    由于游戏中每一时刻往往存在着成千上万个 GO 和相应的回调函数，`我们需要一个非常高效的分发系统才能保证游戏的实时性。`

    - 最简单的分发机制是把消息**瞬时(immediate)**发出去。这种方式的缺陷在于它会阻塞前一个函数的执行，从而`形成一个巨大的调用栈使得系统难以调试`
    - 现代游戏引擎中更常用的分发方式是使用**事件队列(event queue)**来收集当前帧上所有的事件，然后在下一帧再进行分发和处理。一般是 RingBuffer 循环队列。
      由于 event queue 中有不同类型的事件，因此我们还需要结合序列化和反序列化的操作来进行存储
      nameSpace：往往会同时有`多个不同的 event queue 来处理不同类型的事件`，每个 queue 用来保存一些相对独立的事件。这种方式可以便于我们进行调试，也可以提升系统的性能。
      ![alt text](image-53.png)
      当然 event queue 也有一些自身的问题。首先 event queue `无法保证 event 执行的顺序`，同时对于一些实时性的事件 event queue 可能会导致执行的延误。
      ![alt text](image-54.png)

      解决方案：

      - 硬编码
      - immediate、preTick、postTick
      - 采用 Task 机制(如果不满足依赖，直接 return 依赖项，见腾讯文档实践)

- 游戏逻辑

除此之外游戏引擎面对的用户往往是设计师而不是程序员，对于设计师来说直接使用编程语言来设计游戏可能是过于复杂的。
因此`在现代游戏引擎中往往会使用脚本语言来实现游戏的开发工作，`它的优势在于它可以直接在虚拟机上运行。例如 Lua。
![alt text](image-55.png)

- 可视化脚本(visual scripting)
  ![alt text](image-56.png)
  虚幻引擎的 blueprint

  ![control flow](image-57.png)

  缺点：

  - 难进行协作
  - 节点过多时整个脚本在视觉上可能会非常繁琐，让人难以阅读

- 角色(character)、控制(control)以及镜头(camera) `3C系统`
  双人成行
  - 角色系统`一般需要一个非常复杂的状态机模型`来描述角色状态的切换和转移
    ![alt text](image-58.png)
  - 控制系统的核心问题是要兼容玩家不同的输入设备。
    ![alt text](image-59.png)
  - 镜头会直接描述玩家视角中看到的场景和事物。
    ![POV & FOV](image-60.png)
    位置、张角
    镜头系统的一大难点在于如何`根据角色的状态来调整相机的相关参数`，使游戏画面更接近于人眼的真实反映。
    同时镜头系统也需要考虑`各种相机特效的实现，包括镜头抖动、动态模糊等。`
    在复杂场景中往往还会有`多个相机`同时存在，因此我们也需要管理这些相机和相关参数。

16. 游戏中的人工智能 1 – Basic Artificial Intelligence

- 导航(navigation)
  ![寻路三步骤](image-61.png)

  1. Map Representation
     地图是玩家和 NPC 可以行动的区域
     - 网络图(waypoint network)
       ![alt text](image-62.png)
       优点：容易实现
       缺点：倾向于沿路径中心前进而无法利用两边的通道；需要迭代。此在现代游戏中路网的应用并不是很多
     - 网格(grid)
       ![alt text](image-63.png)
       优点：容易实现，支持动态更新，也便于调试
       缺点：内存占用大，难表示重叠区域之间的连接关系
     - `寻路网格(navigation mesh)`，标配
       在寻路网格中可通行的区域会**使用多边形来进行覆盖**，这样可以方便地表达不同区域直接相互连接的拓扑关系。
       ![alt text](image-64.png)
       优点：支持 3 维 ，动态更新
       缺点：生成算法复杂
     - 八叉树(sparse voxel octree)
       如果要制作三维空间中的地图则可以考虑八叉树这样的数据结构。
       ![alt text](image-65.png)
  2. Path Finding
     首先都需要把游戏地图转换为拓扑地图，然后再使用相应的算法进行寻路。
     A Star
     ![alt text](image-66.png)
     在网格地图中常用的启发函数包括 `Manhattan 距离`等。
     而在寻路网格中则可以使用`欧氏距离`作为启发函数。
  3. Path Smoothing
     直接使用寻路算法得到的路径往往包含`各种各样的折线不够光滑`，因此我们还需要使用一些路径平滑的算法来获得更加光滑的路径。
     游戏导航中比较常用 **funnel 算法**来对折线路径进行平滑，它不仅可以应用在二维平面上也可以应用在寻路网格上。
     ![alt text](image-67.png)

     > funnel: 漏斗

     https://lizb0907.github.io/2021/03/26/PathFinding04/

  ***

  NavMesh Generation
  如何从游戏地图上生成寻路网格是一个相对困难的问题

  对于动态的环境我们可以把巨大的场景地图划分为若干个 tile。当某个 tile 中的环境发生改变时`只需要重新计算该处的路径`就可以得到新的路径。

- Steering 算法
  在实际游戏中角色可能包含自身的运动学约束使得我们`无法严格按照计算出的路径进行运动`，这一点对于各种载具尤为明显。因此我们还需要结合 steering 算法来调整实际的行进路径

  steering 算法可以按照行为分为以下几种：追赶和逃脱(seek/flee)、速度匹配(velocity match)以及对齐(align)。

  1. Seek/Flee
     根据自身和目标当前的位置来调整自身的加速度从而实现追赶或是逃脱的行为，像游戏中的`跟踪、躲避或是巡逻`等行为
  2. Velocity Match
     利用当前自身和目标的`相对速度以及匹配时间来进行控制`，使得自身可以按指定的速度到达目标位置
  3. Align
     align 则是从角度和角加速度的层面进行控制，使得自身的朝向可以接近目标

- 群体模拟(crowd simulation)
  在游戏场景中往往会具有大量的 `NPC`，如何控制和模拟群体性的行为是现代游戏的一大挑战。
  目前游戏中群体行为模拟的方法主要可以分为三种：微观模型(microscopic models)、宏观模型(macroscopic models)以及混合模型(mesoscopic models)。

  - 微观模型
    对群体中每一个个体进行控制从而模拟群体的行为，通常情况下我们可以`设计一些规则来控制个体的行为`
  - 宏观模型
    在场景中`设计一个势场或流场`来控制群体中每个个体的行为。
  - 混合模型
    综合了微观和宏观两种模型的思路，它首先把`整个群体划分为若干个小组，然后在每个小组中对每个个体使用微观模型的规则来进行控制`。这样的方法在各种 `RTS 游戏中有着广泛的应用。`

  ***

  Collision Avoidance
  群体模拟中的一大难点在于如何`保证个体之间不会出现碰撞的问题`。

  - Force-based
    比较常用的方法是对每个个体施加一定的力来控制它的运动，这样就可以操纵群体的运动行为。
  - `速度障碍(velocity obstacle, VO)`
    VO 的思想是当两个物体将要发生碰撞时相当于在速度域上形成了一定的障碍，因此需要调整自身的速度来避免相撞。
    当参与避让的个体数比较多时还需要进行一些整体的优化，此时可以使用 ORCA 等算法进行处理。

- 感知(sensing)
  感知(sensing)是游戏 AI 的基础，根据获得信息的不同我们可以把感知的内容分为内部信息(internal information)和外部信息(external information)。
  内部信息包括 AI 自身的位置、HP 以及各种状态。这些信息一般可以被 AI 直接访问到，而且它们是 AI 进行决策的基础。
  外部信息则主要包括 AI 所处的场景中的信息，它会随着游戏进程和场景变化而发生改变。
  `外部信息的一种常用表达方式是 influence map(势力图)，`场景的变化会直接反映在 influence map 上。当 AI 需要进行决策时会同时考虑自身的状态并且查询当前的 influence map 来选择自身的行为。
  http://www.aisharing.com/archives/80
  ![alt text](image-68.png)
- **决策(decision making)系统**
  经典的决策系统包括有限状态机(finite state machine, FSM)和行为树(behavior tree, BT，决策树)两种

  - FSM
    有限状态机的缺陷在于现代游戏中 `AI 的状态空间可能是非常巨大的`，因此状态之间的转移会无比复杂。
    为了克服有限状态机过于复杂的问题，人们还提出了`hierarchical finite state machine(HFSM)这样的模型`。在 HFSM 中我们把整个复杂的状态机分为若干层，不同层之间通过有向的接口进行连接，这样可以增加模型的可读性。
  - **BT**
    ![alt text](image-69.png)
    在现代游戏中更为常用的决策算法是行为树，它的决策行为更接近人脑的决策过程。

    叶子节点：行为节点(execution node)，表示 AI 执行的过程，它包括条件判断以及具体执行的动作两种节点。
    ![alt text](image-70.png)
    非叶子节点：控制节点(control node)，用来控制叶子节点的执行顺序，包括序列节点(sequence)、选择节点(selector)、并行节点(parallel)以及 decorator 四种节点。
    ![alt text](image-71.png)

    1. sequence 是表示对当前节点的子节点`依次`进行访问和执行，一般可以用来表示 AI 在当前状态下的行为计划。有一个节点 fail 则会终止遍历
       ![alt text](image-72.png)
    2. selector 同样会遍历当前节点的子节点，但不同于 sequence 的地方是如果某个子节点返回 True 则会终止遍历
       ![alt text](image-73.png)
    3. parallel 节点会`同时`执行所有的子节点，大于一定数量的子节点返回 True 则会终止遍历
       ![alt text](image-74.png)

    ![行为树](image-75.png)

    行为树的 Tick
    原则：`从根出发，到根结束` (类似问题：react 渲染为什么每次都从根组件开始)
    行为可以被打断
    ![alt text](image-76.png)

    在现代游戏中还提出了 decorator 节点来丰富可以执行的行为。
    相当于语法糖
    ![alt text](image-77.png)

    我们还可以使用 precondition 和 blackboard 来提升决策过程的可读性。
    ![alt text](image-78.png)

    目前随着 AI 技术的发展，游戏 AI 也开始使用一些规划(planning)算法来进行决策。这些更先进的算法我们会在后面的课程进行介绍。

    - Planning and Goal
    - Machine Learning

17. 游戏中的人工智能 2 – Advanced Artificial Intelligence

18.
19. 网络游戏的架构 1 – 基础
20. 网络游戏的架构 2 – 进阶
21. 前沿介绍 1 – Data Oriented Programming，Job System
22. 前沿介绍 2 – Motion Matching, Nanite, Lumen
23. 前沿介绍 3 – Procedurally Generated Content
24. 教学项目 Pilot 源码分析：
25. 项目安装
26. 源码分析 1
27. Bevy 游戏引擎分析：
28. 引擎安装
29. 引擎介绍
30. 游戏开发
31. 源码分析 1
32. 源码分析 2
