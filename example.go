
package main


// Formula_1 是生成的计算函数
func Formula_1(params ...int64) int64 {
	
	c := params[0] // 动态映射变量
	
	a := params[1] // 动态映射变量
	
	b := params[2] // 动态映射变量
	
	return a + b - c
}

// Formula_2 是生成的计算函数
func Formula_2(params ...int64) int64 {
	
	a := params[0] // 动态映射变量
	
	b := params[1] // 动态映射变量
	
	c := params[2] // 动态映射变量
	
	d := params[3] // 动态映射变量
	
	return a * b + c/d
}

// Formula_3 是生成的计算函数
func Formula_3(params ...int64) int64 {
	
	d := params[0] // 动态映射变量
	
	a := params[1] // 动态映射变量
	
	b := params[2] // 动态映射变量
	
	c := params[3] // 动态映射变量
	
	return a + b * c - d
}


// 方法名-方法的映射字典
var FormulaDict = map[string]func(...int64) int64{
	
	"Formula_1": Formula_1,
	"Formula_2": Formula_2,
	"Formula_3": Formula_3,
}
