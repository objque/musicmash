package db

func setup() {
	DbMgr = NewFakeDatabaseMgr()
}

func teardown() {
	_ = DbMgr.DropAllTables()
	_ = DbMgr.Close()
}
