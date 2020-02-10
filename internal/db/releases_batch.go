package db

func (mgr *AppDatabaseMgr) InsertBatchNewReleases(releases []*Release) error {
	const query = "INSERT OR IGNORE INTO releases " +
		"(created_at,artist_id,title,poster,released,store_name,store_id,type,explicit) " +
		"VALUES (:created_at, :artist_id, :title, :poster, :released, :store_name, :store_id, :type, :explicit)"
	_, err := mgr.newdb.NamedExec(query, releases)
	return err
}