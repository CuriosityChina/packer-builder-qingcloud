package qingcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/template/interpolate"
	"os"
)

// Config 配置
type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	// 秘钥配置信息
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`

	// ImageID 映像ID
	// 此映像将作为主机的模板。可传青云提供的映像ID，或自己创建的映像ID
	ImageID string `mapstructure:"image_id"`
	// InstanceType 主机类型，参考 [Instance Types](https://docs.qingcloud.com/api/common/includes/instance_type.html#instance-type)
	// 如果请求中指定了 instance_type，cpu 和 memory 参数可略过。
	// 如果请求中没有 instance_type，则 cpu 和 memory 参数必须指定。
	// 如果请求参数中既有 instance_type，又有 cpu 和 memory，则以 cpu, memory 的值为准。
	InstanceType string `mapstructure:"instance_type"`
	// CPU core
	// 有效值为: 1, 2, 4, 8, 16
	CPU int `mapstructure:"cpu"`
	// Memory 内存
	// 有效值为: 1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768
	Memory int `mapstructure:"memory"`
	// InstanceClass 主机性能类型: 性能型:0 ,超高性能型:1
	InstanceClass string `mapstructure:"instance_class"`

	NeedUserdata string `mapstructure:"need_userdata"`
	Userdata     string `mapstructure:"userdata"`
	UserdataPath string `mapstructure:"userdata_path"`
	UserdataFile string `mapstructure:"userdata_file"`
	Zone         string `mapstructure:"zone"`

	SSHUsername string `mapstructure:"ssh_username"`

	ctx interpolate.Context
}

// NewConfig 创建新的配置
func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	warnings := []string{}

	err := config.Decode(c, &config.DecodeOpts{
		Interpolate: true,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return nil, warnings, err
	}

	// 如果 APIKey / APISecret 为空，则从环境变量中获取
	if c.APIKey == "" {
		c.APIKey = os.Getenv("QINGCLOUD_API_KEY")
	}
	if c.APISecret == "" {
		c.APISecret = os.Getenv("QINGCLOUD_API_SECRET")
	}

	// ...

	return c, warnings, nil
}
