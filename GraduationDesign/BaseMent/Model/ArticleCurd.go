package Model

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	time2 "time"
)

// ForwardArticle Forward A Article
func ForwardArticle(articleId uint, UserId string) int {
	var forward Forward
	Config.DB.Where("user_id = ?", UserId).First(&forward)
	if forward.UserId != "" || forward.ArticleId != 0 {
		return utils.SUCCESS
	}
	err = Config.DB.Model(&Article{}).Where("article_id = ?", articleId).Update("forward", gorm.Expr("forward+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// RecordReadCount Record The ReadCount
func RecordReadCount(ArticleId uint) int {
	Config.DB.Model(&Article{}).Where("id = ?", ArticleId).Update("read_count", gorm.Expr("read_count+ ?", 1))
	return utils.SUCCESS
}

// CheckArticleLike Check A Article Is Like or Not
func CheckArticleLike(ArticleId uint, UserId string) bool {
	var like ArticleLike
	Config.DB.Where("article_id = ? and user_id = ?", ArticleId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// LikeArticle Like A Article
func LikeArticle(ArticleId uint, userId string) int {
	var article Article
	var like ArticleLike
	Config.DB.Where("user_id = ? and article_id = ?", userId, ArticleId).First(&like)
	if like.ID == 0 {
		like = ArticleLike{ArticleId: ArticleId, UserId: userId, Like: true}
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
	err = Config.DB.Model(&article).Where("id = ?", ArticleId).Update("like_count", gorm.Expr("like_count+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// DisLikeArticle DisLike A Article
func DisLikeArticle(ArticleId uint, userId string) int {
	var article Article
	var like ArticleLike
	Config.DB.Where("user_id = ? and article_id = ?", userId, ArticleId).First(&like)
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
	err = Config.DB.Model(&article).Where("id = ?", ArticleId).Update("like_count", gorm.Expr("like_count- ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CreatArticle Creat A Article
func CreatArticle(data *Article) (int, uint) {
	err = Config.DB.Create(&data).Error
	if err != nil {
		return utils.ERROR, 0
	}
	if data.Topic1 != "" {
		CreatTopic(data.Topic1, 1)
	}
	if data.Topic2 != "" {
		CreatTopic(data.Topic2, 1)
	}
	if data.Topic3 != "" {
		CreatTopic(data.Topic3, 1)
	}
	return utils.SUCCESS, data.ID
}

// CreatTopic Creat A Topic
func CreatTopic(data string, num int) int {
	var topic Topic
	Config.DB.Where("t_name = ?", data).Find(&topic)
	if topic.ID == 0 {
		topic = Topic{TName: data, Use: 1, CommentCount: 0, ReadCount: 0}
		err = Config.DB.Create(&topic).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&topic).Where("t_name = ?", data).Update("use", gorm.Expr("use+ ?", num)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// UpdateTopic Update A Topic
func UpdateTopic(data string, commentNum int, readNum int) int {
	var topic Topic
	Config.DB.Where("t_name = ?", data).Find(&topic)
	if topic.ID == 0 {
		topic = Topic{TName: data, Use: 1, CommentCount: 0, ReadCount: 0}
		err = Config.DB.Create(&data).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&topic).Where("t_name = ?", data).Update("comment_count", gorm.Expr("comment_count+ ?", commentNum)).Error
	err = Config.DB.Model(&topic).Where("t_name = ?", data).Update("read_count", gorm.Expr("read_count+ ?", readNum)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditArticlePicture  Edit ArticlePicture
func EditArticlePicture(id uint, ArticlePicture string) int {
	var article Article
	var maps = make(map[string]interface{})
	maps["img"] = ArticlePicture
	err = Config.DB.Model(&article).Where("id = ?", id).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func ScheduledArticleTask() {
	spec := "0 0 3 1 * ?"
	spec1 := "0 0 2 * * ?"
	utils.ScheduledUpdateTask(ClearTheTopicStore, spec)
	utils.ScheduledUpdateTask(CountTopics, spec1)
}

func GetTheTopicList() ([]Topic, int64) {
	var Topics []Topic
	Config.DB.Limit(100).Order("comment_count+read_count+use DESC").Find(&Topics)
	if len(Topics) == 0 {
		return nil, 0
	}
	return Topics, int64(len(Topics))
}

// ClearTheTopicStore Clean every 30 days 3:00 to start
func ClearTheTopicStore() {
	time := time2.Now().Add(-3 * time2.Hour)
	err = Config.DB.Where("Updated_At <= ?", time).Delete(&Topic{}).Error
}

// CountTopics every 1 day or 12 hours to start
func CountTopics() {
	var articles []Article
	Config.DB.Find(&articles)
	for _, article := range articles {
		var top Topic
		if article.Topic1 != "" {
			Config.DB.Select("t_name = ?", article.Topic1).First(&top)
			UpdateTopic(top.TName, article.CommentCount, article.ReadCount)
		}
		if article.Topic2 != "" {
			Config.DB.Select("t_name = ?", article.Topic2).First(&top)
			UpdateTopic(top.TName, article.CommentCount, article.ReadCount)
		}
		if article.Topic2 != "" {
			Config.DB.Select("t_name = ?", article.Topic2).First(&top)
			UpdateTopic(top.TName, article.CommentCount, article.ReadCount)
		}
	}
}

// DeleteArticle Delete A Article
func DeleteArticle(id uint) int {
	var article Article
	Config.DB.Where("id = ?", id).First(&article)
	if article.ID == 0 {
		return utils.ErrorArticleNotExist
	}
	err = Config.DB.Preload("Comments.Replys").Preload("Comments.CommentLike").Select(clause.Associations).Delete(&article).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CheckArticle Check A Article is exit or not
func CheckArticle(id uint) int {
	var article Article
	Config.DB.Where("id = ?", id).First(&article)
	if article.ID == 0 {
		return utils.ErrorArticleNotExist
	}
	return utils.SUCCESS
}

// CheckAArticle Check A Article
func CheckAArticle(id uint) (Article, int) {
	var article Article
	Config.DB.Where("id = ?", id).First(&article)
	if article.ID == 0 {
		return article, utils.ErrorArticleNotExist
	}
	return article, utils.SUCCESS
}

// EditArticle Edit A Article
func EditArticle(id uint, data *Article) (Article, int) {
	var article Article
	var maps = make(map[string]interface{})
	maps["content"] = data.Content
	maps["img"] = data.Img
	err = Config.DB.Model(&article).Where("id = ?", id).Updates(&maps).Error
	if err != nil {
		return article, utils.ERROR
	}
	return article, utils.SUCCESS
}

// SearchActivitiesAndComments Search The activities And Comments
func SearchActivitiesAndComments(title string, pageSize int, pageNum int) ([]Article, []Comment, int, int64) {
	var articles []Article
	var comments []Comment
	var total int64
	var total2 int64
	err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
		"%"+title+"%").Find(&articles).Count(&total).Error
	if err != nil {
		return articles, comments, utils.ERROR, 0
	}
	err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
		"%"+title+"%").Find(&comments).Count(&total2).Error
	if err != nil {
		return articles, comments, utils.ERROR, 0
	}
	return articles, comments, utils.SUCCESS, total2 + total
}

// SearchActivities Search The activities
func SearchActivities(title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articles []Article
	var total int64
	err = Config.DB.Order("created_at DESC").Where("content LIKE ?",
		"%"+title+"%").Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
	//	"%"+title+"%").Find(&articles).Count(&total).Error
	if err != nil {
		return articles, utils.ERROR, 0
	}
	return articles, utils.SUCCESS, total
}

// SearchArticleDays Search The activities
func SearchArticleDays(title int, pageSize int, pageNum int) ([]Article, int, int64) {
	startOfDay := time2.Now().AddDate(0, 0, -title).Truncate(24 * time2.Hour)
	endOfDay := startOfDay.Add(24 * time2.Hour)
	var articles []Article
	var total int64
	err = Config.DB.Order("created_at DESC").Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
	//	"%"+title+"%").Find(&articles).Count(&total).Error
	if err != nil {
		return articles, utils.ERROR, 0
	}
	return articles, utils.SUCCESS, total
}

// SearchTopics Search The activities
func SearchTopics(title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articles []Article
	var total int64
	err = Config.DB.Order("created_at DESC").Where("topic1 LIKE ? or topic2 LIKE ? or topic3 LIKE ?",
		"%"+title+"%", "%"+title+"%", "%"+title+"%").Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("content LIKE ?",
	//	"%"+title+"%").Find(&articles).Count(&total).Error
	if err != nil {
		return articles, utils.ERROR, 0
	}
	return articles, utils.SUCCESS, total
}

// SearchUserActivities Search User Activities
func SearchUserActivities(userid string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articles []Article
	var total int64
	err = Config.DB.Order("Created_At DESC").Where("user_id =  ? ",
		userid).Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id =  ? ",
	//	userid).Find(&articles).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return articles, utils.ERROR, 0
	}
	return articles, utils.SUCCESS, total
}

// SearchUserActivities7Days Search User Activities 7Days
func SearchUserActivities7Days(userid string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articles []Article
	var total int64
	err = Config.DB.Order("Created_At DESC").Where("user_id =  ? and created_at >= ?",
		userid, time2.Now().AddDate(0, 0, -7)).Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id =  ? ",
	//	userid).Find(&articles).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return articles, utils.ERROR, 0
	}
	return articles, utils.SUCCESS, total
}

// SearchUserActivitiesAndComments Search User activities And Comments
func SearchUserActivitiesAndComments(userId string, pageSize int, pageNum int) ([]Article, []Comment, int, int64) {
	var articles []Article
	var comments []Comment
	var total int64
	var total2 int64
	err = Config.DB.Order("Created_At DESC").Where("user_id = ?",
		userId).Find(&articles).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ?",
	//	userId).Find(&articles).Count(&total).Error
	if err != nil {
		return articles, comments, utils.ERROR, 0
	}
	err = Config.DB.Order("Created_At DESC").Where("user_id = ?",
		userId).Find(&comments).Count(&total2).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ?",
	//	userId).Find(&comments).Count(&total2).Error
	if err != nil {
		return articles, comments, utils.ERROR, 0
	}
	return articles, comments, utils.SUCCESS, total2 + total
}

// UpLoadPicture UpLoad Picture
func UpLoadPicture() {

}
