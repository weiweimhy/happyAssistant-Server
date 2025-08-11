package customUtils

func ToMongoDocs[T any](data []*T) []interface{} {
	docs := make([]interface{}, len(data))
	for i, item := range data {
		docs[i] = item
	}
	return docs
}
