// 它在基类中定义了一个算法的框架，允许子类在不修改结构的情况下重写算法的特定步骤

trait IReactComponent {
    fn component_did_mount(&self);
    fn component_will_unmount(&self);
    fn render(&self);
}
