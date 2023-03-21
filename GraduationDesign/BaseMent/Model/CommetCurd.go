package Model

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"gorm.io/gorm"
)

// CheckCommentLike Check A Comment Is Like or Not
func CheckCommentLike(CommentId uint, UserId string) bool {
	var like CommentLike
	Config.DB.Where("comment_id = ? and user_id = ?", CommentId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// LikeComment Like a Comment
func LikeComment(ArticleId uint, id uint, userId string) int {
	var comment Comment
	var like CommentLike
	Config.DB.Where("user_id = ? and comment_id = ? and article_id = ?", userId, id, ArticleId).First(&like)
	if like.ID == 0 {
		like = CommentLike{CommentId: id, UserId: userId, ArticleId: ArticleId, Like: true}
		err = Config.DB.Create(&like).Error
		if err != nil {
			return utils.ERROR
		}
	} else {
		if like.Like == true {
			return utils.ERROR
		}
		err = Config.DB.Model(&like).Update("like", true).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&comment).Where("id = ? and article_id = ?", id, ArticleId).Update("like_count", gorm.Expr("like_count+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// DisLikeComment DisLike a Comment
func DisLikeComment(ArticleId uint, id uint, userId string) int {
	var comment Comment
	var like CommentLike
	Config.DB.Where("user_id = ? and comment_id = ? and article_id = ?", userId, id, ArticleId).First(&like)
	if like.ID == 0 {
		return utils.SUCCESS
	} else {
		if like.Like == false {
			return utils.ERROR
		}
		err = Config.DB.Model(&like).Update("like", false).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&comment).Where("id = ? and article_id = ?", id, ArticleId).Update("like_count", gorm.Expr("like_count- ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CreatAComment creat a comment
func CreatAComment(data *Comment) int {
	err = Config.DB.Create(&data).Error
	err = Config.DB.Model(&Article{}).Where("id = ?", data.ArticleId).Update("comment_count", gorm.Expr("comment_count+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// DeleteComment del a comment
func DeleteComment(id uint) int {
	var comment Comment
	err = Config.DB.Select("CommentLike").Where("id = ?", id).Find(&comment).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// GetaComment get a comment
func GetaComment(id uint) (Comment, int) {
	var comment Comment
	Config.DB.Where("id = ?", id).First(&comment)
	if comment.ID <= 0 || id <= 0 {
		return comment, utils.ERROR
	}
	return comment, utils.SUCCESS
}

// GetUserAllComments GetUserAllCommentsByUser_id
func GetUserAllComments(id string, pageSize int, pageNum int) ([]Comment, int64) {
	var comments []Comment
	var total int64
	if id != "" {
		Config.DB.Where("user_id = ?", id).Order("creat_time DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comments).Offset(-1).Count(&total)
		return comments, total
	}
	return comments, 0
}

// GetUserAllCommentsByName GetUserAllCommentsByName
func GetUserAllCommentsByName(UserName string, pageSize int, pageNum int) ([]Comment, int64) {
	var comments []Comment
	var total int64
	if UserName != "" {
		Config.DB.Where("username = ?", UserName).Order("creat_time DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comments).Offset(-1).Count(&total)
		return comments, total
	}
	return comments, 0
}

// SearchComments Search TheComments
func SearchComments(title string, pageSize int, pageNum int) ([]Comment, int, int64) {
	var comments []Comment
	var total int64
	err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
		"%"+title+"%").Find(&comments).Count(&total).Error
	if err != nil {
		return comments, utils.ERROR, 0
	}
	return comments, utils.SUCCESS, total
}

// GetCommentListByArticleId get a list
func GetCommentListByArticleId(Id uint, pageSize int, pageNum int) ([]Comment, int64) {
	var comments []Comment
	var total int64
	//var replys []Comment
	//var data []Comment
	_ = Config.DB.Where("article_id = ? and r_id = 0", Id).Order("created_at DESC").Find(&comments).Count(&total).Error
	for i := 0; i < len(comments); i++ {
		_ = Config.DB.Where("parent_id = ? ", comments[i].ID).Order("created_at DESC").Find(&comments[i].Replys)
	}
	//Config.DB.Where("article_id = ? and r_id = 0", Id).Order("creat_time DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comments).Offset(-1).Count(&total)
	//floors := make([]uint, total)
	//for i, comment := range comments {
	//	floors[i] = comment.ID
	//}
	//Config.DB.Where("parent_id in (?)", floors).Order("created_at DESC").Find(&replys)
	//comments = append(comments, replys...)
	//data = CreatCommentOrder(comments)
	return comments, total
}

// GetCommentListByArticleIdT get a list
func GetCommentListByArticleIdT(Id uint, id2 uint, pageSize int, pageNum int) ([]Comment, int64) {
	var comments []Comment
	var total int64
	//var replys []Comment
	//var data []Comment
	_ = Config.DB.Where("article_id = ? and type =? and r_id = 0", Id, id2).Order("created_at DESC").Find(&comments).Count(&total).Error
	for i := 0; i < len(comments); i++ {
		_ = Config.DB.Where("parent_id = ? ", comments[i].ID).Order("created_at DESC").Find(&comments[i].Replys)
	}
	//Config.DB.Where("article_id = ? and r_id = 0", Id).Order("creat_time DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comments).Offset(-1).Count(&total)
	//floors := make([]uint, total)
	//for i, comment := range comments {
	//	floors[i] = comment.ID
	//}
	//Config.DB.Where("parent_id in (?)", floors).Order("created_at DESC").Find(&replys)
	//comments = append(comments, replys...)
	//data = CreatCommentOrder(comments)
	return comments, total
}

// GetCommentListInAComment Get CommentList In A Comment
func GetCommentListInAComment(Id uint, RId uint, pageSize int, pageNum int) ([]Comment, int64) {
	var comments []Comment
	var total int64
	var replys []Comment
	var data []Comment
	Config.DB.Where("id = ? and r_id = ?", Id, RId).Order("created_at DESC").Find(&comments).Count(&total)
	//Config.DB.Where("id = ? and r_id = ?", Id, RId).Order("created_at DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&comments).Offset(-1).Count(&total)
	floors := make([]uint, total)
	for i, comment := range comments {
		floors[i] = comment.ID
	}
	Config.DB.Where("parent_id in (?)", floors).Order("created_at DESC").Find(&replys)
	comments = append(comments, replys...)
	data = CreatCommentOrder(comments)
	return data, total
}

// CreatCommentOrder CommentOrder
func CreatCommentOrder(data []Comment) []Comment {
	var Oders []Comment
	var parents, childs []Comment
	for _, comment := range data {
		if comment.ParentId == 0 {
			parents = append(parents, comment)
		}
		childs = append(childs, comment)
	}
	for _, comment := range parents {
		RecursiveComment(&comment, childs)
		Oders = append(Oders, comment)
	}
	return Oders
}

// RecursiveComment Recursive to a list
func RecursiveComment(tree *Comment, nodes []Comment) {
	for _, comment := range nodes {
		if comment.ParentId == 0 {
			continue
		}
		if tree.ID == comment.RId {
			RecursiveComment(&comment, nodes)
			tree.Replys = append(tree.Replys, comment)
		}
	}
	return
}
