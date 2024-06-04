package core

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// CheckImplementation 检查给定的 struct 是否实现了指定的接口
func CheckImplementation(moduleDir, interfaceName, structName string) (bool, error) {
	// 加载整个模块的类型信息
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  moduleDir, // 指定你的模块根目录
	}
	pkgs, err := packages.Load(cfg, "...")
	if err != nil {
		return false, err
	}

	var myInterface, myStruct types.Type

	// 遍历所有包，查找接口和结构体
	for _, pkg := range pkgs {
		for _, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			switch obj.Name() {
			case interfaceName:
				if _, ok := obj.Type().Underlying().(*types.Interface); ok {
					myInterface = obj.Type()
				}
			case structName:
				myStruct = obj.Type()
			}
		}
	}

	if myInterface == nil || myStruct == nil {
		return false, fmt.Errorf("could not find %s or %s", interfaceName, structName)
	}

	// 检查 struct 是否实现了接口
	return types.Implements(myStruct, myInterface.Underlying().(*types.Interface)), nil
}
