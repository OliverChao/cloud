package forever

import (
	"github.com/go-ego/riot/types"
	"github.com/sirupsen/logrus"
	"strconv"
)

func InitRiot_Test() {
	text := "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄"
	text1 := "在IMAX影院放映时"
	text2 := "全片以上下扩展至IMAX 1.9：1的宽高比来呈现"
	// 将文档加入索引，docId 从1开始
	searcher.Index("1", types.DocData{Content: text, Attri: "A"})
	searcher.Index("2", types.DocData{Content: text1, Attri: "B"})
	searcher.Index("3", types.DocData{Content: text2, Attri: "C"})
	// 等待索引刷新完毕
	searcher.Flush()
}

func InitRiot() {
	client := GetGlobalRedisClient()
	all := client.HGetAll("kinds")
	m1 := all.Val()
	dataM := map[string][]string{}
	for k, v := range m1 {
		if v == "0" {
			continue
		}
		m2 := client.HGetAll(k).Val()
		slist := make([]string, 0)
		for k2, _ := range m2 {
			slist = append(slist, k2)
		}
		dataM[k] = slist
	}
	logrus.Info("dataM ::", dataM)
	AddMultiDocsByMap(dataM)
}

func addNewDoc(indexed uint64, kind, title string) {
	//indexed := searcher.NumIndexed()
	//indexed += 1
	id := strconv.FormatUint(indexed, 10)
	searcher.Index(id, types.DocData{
		Content: title,
		Attri:   kind,
	})
	//searcher.Flush()
}

func AddDoc(kind, title string) {
	indexed := searcher.NumIndexed()
	indexed += 1
	addNewDoc(indexed, kind, title)
	searcher.Flush()
}

func AddMultiDocs(kind string, titles ...string) {
	indexed := searcher.NumIndexed()
	indexed += 1
	for k, v := range titles {
		id := indexed + uint64(k)
		addNewDoc(id, kind, v)
	}
	searcher.Flush()
}

func AddMultiDocsByMap(m map[string][]string) {
	for k, v := range m {
		AddMultiDocs(k, v...)
	}
	searcher.Flush()
}

func AndSearch(query string) *types.SearchDoc {
	output := _search("must", query)
	logrus.Info(output.Docs)
	return output
}

func OrSearch(query string) *types.SearchDoc {
	output := _search("should", query)
	logrus.Info(output)
	return output
}

func _search(searchType string, query string) *types.SearchDoc {
	var logic types.Logic
	if searchType == "must" {
		logic = types.Logic{
			Must: true,
		}
	} else {
		logic = types.Logic{
			Should: true,
		}
	}
	req := types.SearchReq{
		Text:    query,
		Logic:   logic,
		Timeout: 3000,
	}
	output := searcher.SearchDoc(req)
	return &output
}

func DeleteDoc(docId string) {
	go searcher.RemoveDoc(docId)
}

//func DeleteDocByNameKind(kind,name string){
//	var logic = types.Logic{
//		Must:   true,
//		Expr:   types.Expr{
//			Must:   nil,
//			Should: nil,
//			NotIn:  nil,
//		},
//	}
//}
