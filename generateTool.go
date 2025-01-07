package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"strconv"
	"sync"
	"text/template"
)

// 定义生成代码的模板
const codeTemplate = `
package main

{{range .Formulas}}
// {{.FuncName}} 是生成的计算函数
func {{.FuncName}}(params ...int64) int64 {
	{{range $index, $var := .Variables}}
	{{$var}} := params[{{$index}}] // 动态映射变量
	{{end}}
	return {{.FormulaLogic}}
}
{{end}}

// 方法名-方法的映射字典
var FormulaDict = map[string]func(...int64) int64{
	{{range .Formulas}}
	"{{.FuncName}}": {{.FuncName}},{{end}}
}
`

// 公式信息
type FormulaInfo struct {
	FuncName     string   // 方法名
	FormulaLogic string   // 公式逻辑
	Variables    []string // 公式中的变量
}

// 全局唯一 ID 生成器和缓存
var (
	idCounter    int
	idLock       sync.Mutex
	formulaCache = make(map[string]string) // 缓存：公式逻辑 -> 方法名
	cacheLock    sync.RWMutex
)

// 生成唯一的 FuncName（带查重）
func generateFuncName(formula string) (string, bool) {
	// 先检查缓存中是否已存在
	cacheLock.RLock()
	if funcName, ok := formulaCache[formula]; ok {
		cacheLock.RUnlock()
		return funcName, true // 返回已存在的方法名和 true 表示重复
	}
	cacheLock.RUnlock()

	// 生成新的唯一 ID
	idLock.Lock()
	defer idLock.Unlock()
	idCounter++
	funcName := "Formula_" + strconv.Itoa(idCounter)

	// 存入缓存
	cacheLock.Lock()
	formulaCache[formula] = funcName
	cacheLock.Unlock()

	return funcName, false // 返回新生成的方法名和 false 表示不重复
}

// 解析公式并提取变量
func extractVariables(formula string) ([]string, error) {
	// 解析公式为AST
	expr, err := parser.ParseExpr(formula)
	if err != nil {
		return nil, fmt.Errorf("failed to parse formula: %v", err)
	}

	// 提取变量
	variables := make(map[string]bool)
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			variables[ident.Name] = true
		}
		return true
	})

	// 将变量按字母顺序排序（保证一致性）
	var sortedVars []string
	for v := range variables {
		sortedVars = append(sortedVars, v)
	}
	return sortedVars, nil
}

// 生成 Go 代码并保存到文件
func generateFormulaCode(formulas []FormulaInfo, outputPath string) error {
	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// 解析模板
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// 模板数据
	data := map[string]interface{}{
		"Formulas": formulas,
	}

	// 生成代码并写入文件
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to generate code: %v", err)
	}

	fmt.Printf("Code generated successfully and saved to %s\n", outputPath)
	return nil
}

func main() {
	// 定义多个公式
	rawFormulas := []string{
		"a + b - c",     // 示例逻辑
		"a * b + c/d",   // 示例逻辑
		"a + b - c",     // 重复公式
		"a + b * c - d", // 新公式
	}

	// 自动提取变量并生成 FuncName
	var formulas []FormulaInfo
	for _, formula := range rawFormulas {
		variables, err := extractVariables(formula)
		if err != nil {
			fmt.Printf("Error extracting variables for formula: %v\n", err)
			return
		}

		funcName, isDuplicate := generateFuncName(formula)
		if isDuplicate {
			fmt.Printf("Skipping duplicate formula: %s\n", formula)
			// 跳过重复公式
			continue
		}

		formulas = append(formulas, FormulaInfo{
			FuncName:     funcName,
			FormulaLogic: formula,
			Variables:    variables,
		})
	}

	// 生成代码并保存到文件
	outputPath := "formula_generated.go"
	if err := generateFormulaCode(formulas, outputPath); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
