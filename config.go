package config

import (
	"errors"
	"github.com/go-ini/ini"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Configurator interface {
	Get(keys string, defaults ...interface{}) (interface{}, error)
	Set(key string, value string) error
	All() map[string]*ini.File
}

var config *Config

type Config struct {
	directory string
	configs   map[string]*ini.File
	delimiter string
	extension string
}

// 构造函数
func NewConfig(directory string) (*Config, error) {
	// singleton
	if config != nil {
		return config, nil
	}

	directory, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	config := &Config{directory: directory, delimiter: `.`, extension: `.conf`}

	// loadAll
	err = config.loadAll()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// 获取指定key配置
func (this *Config) Get(keys string, defaults ...interface{}) (interface{}, error) {

	keySlices := parseKey(keys, this.delimiter)

	defaultValue := ``
	if len(defaults) >= 1 {
		defaultValue = defaults[0].(string)
	}

	length := len(keySlices)

	var cfg *ini.File
	if _, ok := this.configs[keySlices[0]]; !ok {
		return defaultValue, nil
	} else {
		cfg = this.configs[keySlices[0]]
	}

	if length == 1 {
		return cfg, nil
	} else if length == 2 {
		if !cfg.Section(ini.DefaultSection).HasKey(keySlices[1]) {
			return defaultValue, nil
		}
		return cfg.Section(ini.DefaultSection).GetKey(keySlices[1])
	} else {
		// 取得可能是section的slice拼接
		sectionName := strings.Join(keySlices[1:length-1], this.delimiter)
		section, err := cfg.GetSection(sectionName)
		// 如果section不存在拼接全部的
		if err != nil {
			sectionName = sectionName + `.` + keySlices[length-1:][0]
			section, err = cfg.GetSection(sectionName)
			if err != nil {
				return nil, err
			}

			return section, nil
		}

		// 否则最后一位肯定是key
		key := keySlices[length-1:][0]
		if section.HasKey(key) {
			return section.GetKey(key)
		}

		return defaultValue, nil
	}
}

// 设置配置
func (this *Config) Set(keys string, value string) error {

	keySlices := parseKey(keys, this.delimiter)

	length := len(keySlices)

	if length == 1 {
		return errors.New("error")
	}

	var err error

	if _, ok := this.configs[keySlices[0]]; !ok {
		this.configs[keySlices[0]], err = loadConf(this.fullPath(keySlices[0]))
		if err != nil {
			return err
		}
	}

	if length == 2 {
		_, err = this.configs[keySlices[0]].Section(ini.DefaultSection).NewKey(keySlices[1], value)
	} else {
		_, err = this.configs[keySlices[0]].Section(strings.Join(keySlices[1:length-1], this.delimiter)).NewKey(keySlices[length-1], value)
	}

	if err != nil {
		return err
	}

	file, err := os.OpenFile(this.fullPath(keySlices[0]), os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = this.configs[keySlices[0]].WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

// 获取所有配置
func (this *Config) All() map[string]*ini.File {
	return this.configs
}

// 一次加载所有配置文件
func (this *Config) loadAll() error {
	// init File map
	this.configs = make(map[string]*ini.File)

	err := filepath.Walk(this.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)

		if strings.Contains(filename, this.extension) && !info.IsDir() {
			fileKey := strings.Replace(filename, this.extension, "", 1)

			cfg, err := loadConf(path)
			if err != nil {
				return err
			}

			this.configs[fileKey] = cfg
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 通过文件名（不包含扩展名），得到当前的完整路径
func (this *Config) fullPath(filename string) string {
	return path.Clean(this.directory + "/" + filename + this.extension)
}

// 加载配置文件
func loadConf(filename string) (*ini.File, error) {

	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			defer file.Close()
		}
	}

	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// 解析数据key
func parseKey(key string, delimiter string) []string {
	return strings.Split(key, delimiter)
}
