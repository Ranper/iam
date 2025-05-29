package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"

	genericoptions "github.com/Ranper/iam/internal/pkg/options"
	"github.com/Ranper/iam/internal/pkg/server"
	"github.com/Ranper/iam/pkg/log"
)

type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server" mapstructure:"server"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	Log                     *log.Options                           `json:"log" mapstructure:"log"`
}

// NewOptions 创建一个带[默认参数]的选项
func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		Log:                     log.NewOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.SecureServing.AddFlags(fss.FlagSet("secure"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
