package middlewares

// 文件元信息结构
type FileMeta struct {
	FileSha1 string //文件唯一标志
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta //key:hash code  ; value: filemeta

// 初始化：
func init() {
	fileMetas = make(map[string]FileMeta)
}

// 接口：filemeta更新;新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 通过hash code获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
