package yaml

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// マッピングする構造体です。
// どういうわけか、大文字じゃないとダメのようです。
type Test struct {
	First  string
	Second string
}

func main() {
	buf, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// []byte を []Test に変換します。
	data, err := ReadOnStruct(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

// yaml形式の[]byteを渡すと[]Testに変換してくれる関数です。
func ReadOnStruct(fileBuffer []byte) ([]Test, error) {
	var data []Test
	// []map[string]string のときと使う関数は同じです。
	// いい感じにマッピングしてくれます。
	err := yaml.Unmarshal(fileBuffer, &data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}
