package app

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

// CliOptions 从命令行中读取参数获取配置选项的抽象
type CliOptions interface {
	// 添加flag到指定的FlagSet obj
	// AddFlags(fs *pflag.FlagSet)
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}

// ConfigurableOptions 从配置文件中读取参数获取配置选项的抽象
type ConfigurableOptions interface {
	// ApplyFlags 解析命令行或者配置文件到options 实例
	ApplyFlags()
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
