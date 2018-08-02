package db

func setup() {
	DbMgr = NewFakeDatabaseMgr()
}

func teardown() {
	DbMgr.DropAllTables()
	DbMgr.Close()
}
