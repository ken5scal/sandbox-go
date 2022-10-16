package controllers_test

import (
	"github.com/ken5scal/api-go-mid-level/controllers"
	"github.com/ken5scal/api-go-mid-level/controllers/testdata"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var aCon *controllers.ArticleController

func TestMain(m *testing.M) {
	ser := testdata.NewServiceMock()
	aCon = controllers.NewArticleController(ser)

	m.Run()
}
