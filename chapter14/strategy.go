package chapter14

type IStrategy interface {
	do(int, int) int
}

type add struct{}

func (impl *add) do(a, b int) int {
	return a + b
}

type reduce struct{}

func (impl *reduce) do(a, b int) int {
	return a - b
}

// 具体策略的执行者
type Operator struct {
	strategy IStrategy
}

// 设置策略
func (impl *Operator) setStrategy(strategy IStrategy) {
	impl.strategy = strategy
}

// 调用策略中的方法
func (impl *Operator) calculate(a, b int) int {
	return impl.strategy.do(a, b)
}
