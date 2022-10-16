package services

import (
	"database/sql"
	"errors"
	"github.com/ken5scal/api-go-mid-level/apperrors"
	"github.com/ken5scal/api-go-mid-level/models"
	"github.com/ken5scal/api-go-mid-level/repositories"
)

// PostArticleHandlerで使うことを想定したサービス
// 引数の情報をもとに新しい記事を作り、結果を返却
func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "fail to record data")
		return models.Article{}, err
	}
	return newArticle, nil
}

// ArticleListHandlerで使うことを想定したサービス
// 指定pageの記事一覧を返却
func (s *MyAppService) GetArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.SelectArticleList(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return nil, err
	}

	if len(articleList) == 0 {
		err := apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}
	return articleList, nil
}

// ArticleDetailHandlerで使うことを想定したサービス
// 指定IDの記事情報を返却
func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	//article, err := repositories.SelectArticleDetail(s.db, articleID)
	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error
	// avoiding race condition by Lock (sync.Mutex)
	//var amu sync.Mutex    // ロックを司る
	//var cmu sync.Mutex    // ロックを司る
	//var wg sync.WaitGroup // ゴルーチンの待ち合わせを司る（メインルーチンが終わっても、他go routine終了をまつ）
	//wg.Add(2)
	type articleResult struct {
		article models.Article
		err     error
	}
	articleChan := make(chan articleResult)
	defer close(articleChan)

	go func(ch chan<- articleResult, db *sql.DB, articleID int) {
		//go func(db *sql.DB, articleID int) {
		//defer wg.Done()
		//amu.Lock()
		article, articleGetErr = repositories.SelectArticleDetail(db, articleID)
		ch <- articleResult{article: article, err: articleGetErr}
		//amu.Unlock()
		//}(s.db, articleID)
	}(articleChan, s.db, articleID)

	type commentResult struct {
		commentList *[]models.Comment
		err         error
	}
	commentChan := make(chan commentResult)
	defer close(commentChan)
	go func(ch chan<- commentResult, db *sql.DB, articleID int) {
		//go func(db *sql.DB, articleID int) {
		//defer wg.Done()
		//cmu.Lock()
		commentList, commentGetErr = repositories.SelectCommentList(db, articleID)
		ch <- commentResult{commentList: &commentList, err: commentGetErr}
		//cmu.Unlock()
	}(commentChan, s.db, articleID)
	//wg.Wait()

	for i := 0; i < 2; i++ {
		select {
		case ar := <-articleChan:
			article, articleGetErr = ar.article, ar.err
		case cr := <-commentChan:
			commentList, commentGetErr = *cr.commentList, cr.err
		}
	}

	if articleGetErr != nil {
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			articleGetErr = apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, articleGetErr
		}
		articleGetErr = apperrors.GetDataFailed.Wrap(articleGetErr, "fail to get data")
		return models.Article{}, articleGetErr
	}

	if commentGetErr != nil {
		commentGetErr = apperrors.GetDataFailed.Wrap(commentGetErr, "fail to get data")
		return models.Article{}, commentGetErr
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

// PostNiceHandlerで使うことを想定したサービス
// 指定IDの記事のいいね数を+1して、結果を返却
func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateNiceNum(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NoTargetData.Wrap(err, "does not exist target article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "fail to update nice count")
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum + 1,
		CreatedAt: article.CreatedAt,
	}, nil
}
