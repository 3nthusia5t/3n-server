package sqlite

import (
	"backend/article"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

// Test if database initialization works.
func TestInit(t *testing.T) {

	filePath := "../../rsrc/test_database/database.db"

	os.Remove(filePath)

	// Call your Init function
	dbManager := Init(filePath)
	if dbManager == nil || dbManager.db == nil {
		t.Errorf("Failed to initialize database")
	}
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Init did not create file: %s", filePath)
	}

	err = dbManager.db.Close()
	if err != nil {
		t.Errorf("Database failed to close %s", err.Error())
	}
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Database wasn't removed properly: %s", err.Error())
	}
}

func TestDbManager_CreateRecord(t *testing.T) {
	filePath := "../../rsrc/test_database/database.db"
	dbManager := Init(filePath)

	//close and remove database
	cleanup := func(mDb *DbManager) {
		mDb.db.Close()
		os.Remove(mDb.path)
	}

	metaDes := "asdfasdf"
	editInt := uint64(123123)

	type args struct {
		a article.Article
	}
	tests := []struct {
		name    string
		db      *DbManager
		args    args
		wantErr bool
	}{
		{
			name: "Simple record",
			db:   dbManager,
			args: args{
				a: *article.NewFromScratch("asdasd", "asdasdas", "./foobar", []string{"new", "cybersec"}, "firendly-url", 14325345, nil, nil, true),
			},
			wantErr: false,
		},
		{
			name: "Full record",
			db:   dbManager,
			args: args{
				a: *article.NewFromScratch("asd", "asdasdas", "./foobar", []string{"new", "cybersec"}, "firendly-url", 123, &editInt, &metaDes, true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.CreateRecord(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("DbManager.CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	cleanup(dbManager)
}
