package sqlite

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestErrorTranslator(t *testing.T) {
	// This is the DSN of the in-memory SQLite database for these tests.
	const InMemoryDSN = "file:testdatabase?mode=memory&cache=shared"

	// This is the example object for testing the unique constraint error
	type Article struct {
		ArticleNumber string `gorm:"unique"`
	}

	db, err := gorm.Open(&Dialector{DSN: InMemoryDSN}, &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		TranslateError: true})

	if err != nil {
		t.Errorf("Expected Open to succeed; got error: %v", err)
	}
	if db == nil {
		t.Errorf("Expected db to be non-nil.")
	}

	err = db.AutoMigrate(&Article{})
	if err != nil {
		t.Errorf("Expected to migrate database models to succeed: %v", err)
	}

	err = db.Create(&Article{ArticleNumber: "A00000XX"}).Error
	if err != nil {
		t.Errorf("Expected first create to succeed: %v", err)
	}

	err = db.Create(&Article{ArticleNumber: "A00000XX"}).Error
	if err == nil {
		t.Errorf("Expected second create to fail.")
	}

	if err != gorm.ErrDuplicatedKey {
		t.Errorf("Expected error from second create to be gorm.ErrDuplicatedKey: %v", err)
	}
}
