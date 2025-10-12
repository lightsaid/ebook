package app

// type wrapper map[string]interface{}

// func msgWrapp(msg string) wrapper {
// 	return wrapper{"msg": msg}
// }

// func dataWrapp(data interface{}) wrapper {
// 	return wrapper{"data": data}
// }

// func (a *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
// 	maxBytes := 2 << 20 // 2 MB
// 	// 限制请求体大小为2MB1⃣️内
// 	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

// 	// 使用请求体创建一个解码器
// 	dec := json.NewDecoder(r.Body)

// 	err := dec.Decode(dst)
// 	if err != nil {

// 		logger.ErrorfoLog.Println("readJSON failed: " + err.Error())

// 		var syntaxError *json.SyntaxError

// 		switch {
// 		case errors.As(err, &syntaxError):
// 			return errors.New("请输入JSON格式请求体")
// 		case errors.Is(err, io.EOF):
// 			return errors.New("请求体不能为空")
// 		case err.Error() == "http: request body too large":
// 			return errors.New("请求体过大")
// 		default:
// 			return errors.New("未知错误，请检查参数")
// 		}
// 	}

// 	return nil
// }

// func (a *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) {
// 	// 添加请求头，如果有 headers => []http.Header => []map[string][]string
// 	if len(headers) > 0 {
// 		for i := range headers {
// 			for key, val := range headers[i] {
// 				w.Header()[key] = val
// 			}
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	// NOTE: 设置请求最终状态，用于中间件 logRequest 获取记录
// 	w.Header().Set("Stauts", http.StatusText(status))
// 	w.WriteHeader(status)
// 	json.NewEncoder(w).Encode(data)
// }

// func (a *application) errorResponse(w http.ResponseWriter, err error) {
// 	logger.ErrorfoLog.Println("errorResponse -> ", err)
// 	var status = http.StatusInternalServerError
// 	// 判断是否是数据库操作错误
// 	if ok := dbrepo.IsCustomDBError(err); ok {
// 		if errors.Is(err, dbrepo.ErrNotFound) {
// 			status = http.StatusNotFound
// 		} else {
// 			status = http.StatusUnprocessableEntity
// 		}
// 		a.writeJSON(w, status, msgWrapp(err.Error()))
// 		return
// 	}

// 	// 判断是否是上传错误
// 	if fileupload.IsUploaderError(err) {
// 		a.writeJSON(w, http.StatusUnprocessableEntity, msgWrapp(err.Error()))
// 		return
// 	}

// 	a.writeJSON(w, status, "服务开了小差，请稍后重试～")
// }

// func (a *application) methodNotAllowed(w http.ResponseWriter, allowMethod ...string) {
// 	header := http.Header(make(map[string][]string))

// 	var methods []string

// 	methods = append(methods, allowMethod...)

// 	header.Add("Allow", strings.Join(methods, ","))

// 	a.writeJSON(
// 		w,
// 		http.StatusMethodNotAllowed,
// 		wrapper{"error": http.StatusText(http.StatusMethodNotAllowed), "allow": allowMethod},
// 		header,
// 	)
// }
