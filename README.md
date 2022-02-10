#使用client-go 中的三种 client：

Clientset
Clientset 是我们最常用的 client，你可以在它里面找到 kubernetes 目前所有原生资源对应的 client。 获取方式一般是，指定 group 然后指定特定的 version，然后根据 resource 名字来获取到对应的 client。

Dynamic Client
Dynamic client 是一种动态的 client，它能同时处理 kubernetes 所有的资源。并且同时，它也不同于 clientset，dynamic client 返回的对象是一个 map[string]interface{}，如果一个 controller 中需要控制所有的 API，可以使用dynamic client，目前它被用在了 garbage collector 和 namespace controller。

RESTClient
RESTClient 是 clientset 和 dynamic client 的基础，前面这两个 client 本质上都是 RESTClient，它提供了一些 RESTful 的函数如 Get()，Put()，Post()，Delete()。由 Codec 来提供序列化和反序列化的功能。

Informer
● 等待所有的 cache 同步完成: 这是为了避免生成大量无用的资源，比如 replica set controller 需要watch replica sets 和 pods, 在 cache 还没有同步完之前，controller 可能为一个 replica set 创建了大量重复的 pods，因为这个时候 controller 觉得目前还没有任何的 pods。
● 修改 resource 对象前先 deepcopy 一份: 在 Informer 这个模型中，我们的 resource 一般是从本地 cache 中取出的，而本地的 cache 对于用户来说应该是 read-only 的，因为它可能是与其他的 informer 共享的，如果你直接修改 cache 中的对象，可能会引起读写的竞争。
● 处理 DeletedFinalStateUnknown 类型对象: 当你的收到一个删除事件时，这个对象有可能不是你想要的类型，即它可能是一个 DeletedFinalStateUnknown，你需要单独处理它。
● 注意 informer 的 resync 行为， informer 会定期从 apiserver resync 资源，这时候会收到大量重复的更新事件，这个事件有一个特点就是更新的 Object 的 ResourceVersion 是一样的，将这种不必要的更新过滤掉。
● 在创建事件中注意 Object 已经被删除的情况: 在 Controller 重启的过程中，可能会有一些对象被删除了，重启后，Controller 会收到这些已删除对象的创建事件，请把这些对象正确地删除。
● SharedInformer: 建议使用 SharedInformer, 它会在多个 Informer 中共享一个本地 cache，这里有一 个 factory 来方便你编写一个新的 Informer。

Lister
Lister 是用来帮助我们访问本地 cache 的一个组件。

Informer的工作流程：
1. 创建一个控制器
● 为控制器创建 workqueue
● 创建 informer, 为 informer 添加 callback 函数，创建 lister

2. 启动控制器
● 启动 informer
● 等待本地 cache sync 完成后， 启动 workers

3. 当收到变更事件后，执行 callback 
● 等待事件触发
● 从事件中获取变更的 Object
● 做一些必要的检查
● 生成 object key，一般是 namespace/name 的形式
● 将 key 放入 workqueue 中

4. worker loop
● 等待从 workqueue 中获取到 item，一般为 object key
● 用 object key 通过 lister 从本地 cache 中获取到真正的 object 对象
● 做一些检查
● 执行真正的业务逻辑
● 处理下一个 item
