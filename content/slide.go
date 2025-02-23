// https://github.com/yurisasuke/golang-decode-string-json-to-int/blob/main/main.go
package content

type Slide struct {
	Background string                 `json:"background"`
	Bottom     map[string]interface{} `json:"bottom"`
	Duration   StringInt              `json:"duration"`
}
