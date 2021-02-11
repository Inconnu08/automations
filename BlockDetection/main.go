package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type NginxBlock struct {
	StartLine   string
	EndLine     string
	AllContents string
	// split lines by \n on AllContents,
	// use make to create *[],
	// first create make([]*Type..)
	// then use &var to make it *
	AllLines          *[]*string
	NestedBlocks      []*NginxBlock
	TotalBlocksInside int
}

// ngBlock *NginxBlock
func IsBlock(line string) bool {
	rex := regexp.MustCompile(`{((?:[^{}]|{{[^}]*}})*)}`)
	out := rex.FindAllStringSubmatch(line, -1)
	if len(out) != 0 {
		return true
	}
	return false
}

func IsLine(line string) bool {
	b, _ := regexp.MatchString(`(\r\n|\r|\n)`, line)
	return b
}

func HasComment(line string) bool {
	b, _ := regexp.MatchString(`\B(\#\s?[a-zA-Z]+\b)`, line)
	return b
}

type NginxBlocks struct {
	blocks      *[]*NginxBlock
	AllContents string
	// split lines by \n on AllContents
	AllLines *[]*string
}

func GetNginxBlock(lines *[]*string, startIndex, endIndex, recursionMax int) *NginxBlock {
	return nil
}

func GetNginxBlocks(configContent string) *NginxBlocks {
	var ngx NginxBlocks
	ngx.AllContents = configContent
	var lines []*string
	sc := bufio.NewScanner(strings.NewReader(configContent))
	for sc.Scan() {
		text := sc.Text()
		lines = append(lines, &text)
	}
	fmt.Println(lines)
	ngx.AllLines = &lines

	return &ngx
}

func main() {
	ConfigContent := `server { # simple reverse-proxy
    listen       80;
    server_name  domain2.com www.domain2.com;
    access_log   logs/domain2.access.log  main;

    # serve static files
    location ~ ^/(images|javascript|js|css|flash|media|static)/  {
      root    /var/www/virtual/big.server.com/htdocs;
      expires 30d;
    }

    # pass requests for dynamic content to rails/turbogears/zope, et al
    location / {
      proxy_pass      http://127.0.0.1:8080;
    }
}`

	ngx := GetNginxBlocks(ConfigContent)
	fmt.Println(*ngx)

}

//func blockParser(str string) {
//	length := len(str)
//	var stack []int
//	var result []string
//
//	for i := 0; i < length; i++ {
//
//		if str[i] == '{' {
//			stack = append(stack, i)
//		}
//
//		if str[i] == '}' {
//			open := stack[len(stack)-1]
//			fmt.Println(stack)
//			stack = stack[:len(stack)-1]
//			fmt.Println(stack)
//			result = append(result, str[open:i+1])
//			//fmt.Println(str[open: i-open-1])
//		}
//	}
//
//	for i, v := range result {
//		fmt.Println(i)
//		fmt.Println(v)
//	}
//
//	fmt.Println(len(result))
//}

