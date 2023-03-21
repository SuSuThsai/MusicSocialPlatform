package utils

import "strconv"

func GetCacheNameArticle(id uint) string {
	return "article_" + strconv.FormatUint(uint64(id), 10)
}

func GetCacheNameComments(id uint) string {
	return "AArticleComments_" + strconv.FormatUint(uint64(id), 10)
}

func GetCacheNameCommentONE(id uint, id1 uint) string {
	return "articleACommentOne_" + strconv.FormatUint(uint64(id), 10) + strconv.FormatUint(uint64(id1), 10)
}

func GetCacheNameArticleLike(id uint, id1 string) string {
	return "articleLike_" + strconv.FormatUint(uint64(id), 10) + "_" + id1
}

func GetCacheNameCommentLike(id uint, id1 string) string {
	return "commentLike_" + strconv.FormatUint(uint64(id), 10) + "_" + id1
}

func GetCacheNameMusicLike(id uint, id1 uint) string {
	return "musicLike_" + strconv.FormatUint(uint64(id), 10) + "_" + strconv.FormatUint(uint64(id1), 10)
}

func GetCacheNameMusicLikeCount(id uint) string {
	return "musicLikeCount_" + strconv.FormatUint(uint64(id), 10)
}

func GetCacheNameMusicListLike(id uint, id1 uint) string {
	return "musicListLike_" + strconv.FormatUint(uint64(id), 10) + "_" + strconv.FormatUint(uint64(id1), 10)
}

func GetCacheNameMusicRank(id string) string {
	return "musicRank_" + id
}

func GetCacheNameMusicListRank(id string) string {
	return "musicListRank_" + id
}
func GetCacheNameTopic(id string) string {
	return "Rank_" + id
}
