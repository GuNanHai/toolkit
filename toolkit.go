package toolkit

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// StructJSON -- 与Struct 和 Json相互转化有关的接口
type StructJSON interface {
	Print() //简单的打印自己，仅用于构建接口关系
}

// GenJSONString 将单个struct变量值以类json的字符串值输出
func GenJSONString(obj StructJSON) string {
	e := reflect.ValueOf(obj).Elem()
	result := "{"
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type.String()

		if strings.Contains(varType, "int") {
			intNum := int(e.Field(i).Int())
			varValue := strconv.Itoa(intNum)
			result = result + `"` + varName + `":` + varValue + `,`
		} else if strings.Contains(varType, "float") {
			floatNum := float64(e.Field(i).Float())
			varValue := fmt.Sprintf("%f", floatNum)
			result = result + `"` + varName + `":` + varValue + `,`
		} else {
			varValue := e.Field(i).String()
			result = result + `"` + varName + `":"` + varValue + `",`
		}
	}
	result = result[:len(result)-1] + "}"

	//将json字符串中所有key关键字的首字母变成小写
	re := regexp.MustCompile(`"(\w+)":`)
	indexList := re.FindAllStringSubmatchIndex(result, -1)
	for _, each := range indexList {
		i := each[2]
		result = result[0:i] + strings.ToLower(string(result[i])) + result[i+1:]
	}

	return result
}

// Unmarshall : 将json变为Struct 类型，传入的dstVarPtr应为指针，类型以某struct为元素类型的数组类 |
// StructJSON 接口：要求成员struct数组类拥有  Print()  --无返回值 |
// StructJSON 传入参数这里仅接受struct数组的指针
func Unmarshall(srcJSON []byte, dstVarPtr StructJSON) {
	err := json.Unmarshal(srcJSON, dstVarPtr)
	if err != nil {
		fmt.Println("json unmarshal failed:", err)
	}
}

// FormatJSON : 有些json文件内容不够正规，无法被golang的json包解析，所以需要对原json内容进行格式调整
func FormatJSON(jsonBlob string) string {
	reFindQuotedKeyword := regexp.MustCompile(`"(\w+)":`)
	if reFindQuotedKeyword.MatchString(jsonBlob) {
		return jsonBlob
	}
	//为json内容中所有key配上双引号
	re := regexp.MustCompile(`(\w+[\D]):`)
	result := re.ReplaceAllString(jsonBlob, `"$1":`)
	return result
}

// JSONToStruct : 将字符串JSON变为某个Struct数组类型 |
// StructJSON 接口：要求成员struct数组类拥有  Print()  --无返回值 |
// StructJSON 传入参数这里仅接受struct数组的指针
func JSONToStruct(jsonBlob string, dstStruct StructJSON) {
	formattedJSON := FormatJSON(jsonBlob)
	Unmarshall([]byte(formattedJSON), dstStruct)
}

// =====================================ERROR

// CheckErr :
func CheckErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// GetPkgPath : 获取从github下载到本地的指定名称的package的路径
func GetPkgPath(pkgName string) string {
	pkgPath := ""
	fileSelf, err := os.Open("main.go")
	CheckErr(err)
	b := make([]byte, 90000)
	fileSelf.Read(b)

	re := regexp.MustCompile(`import\s+?\((.|\n)*?\)`)
	packageDeclarationTextList := strings.Split(re.FindString(string(b)), "\n")

	for _, each := range packageDeclarationTextList {
		if strings.Contains(each, pkgName) {
			each = strings.TrimSpace(each)
			each = strings.Trim(each, `"`)
			pkgPath = os.Getenv("GOPATH") + string(os.PathSeparator) + "src" + string(os.PathSeparator) + path.Dir(each) + string(os.PathSeparator) + pkgName
		}
	}
	return pkgPath
}
