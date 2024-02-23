package xslices

func Concat(values ...any) []string {
	var res []string
	for _, v := range values {
		switch typ := v.(type) {
		case []string:
			res = append(res, typ...)
		case string:
			res = append(res, typ)
		}
	}
	return res
}
