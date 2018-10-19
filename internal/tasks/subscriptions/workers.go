package subscriptions

func runAppleWorker(id int, jobs <-chan []string) {
	for {
		select {
		case artists, ok := <-jobs:
			if !ok {
				return
			}

			appleLinker.SearchArtists(artists)
		}
	}
}

func runYandexWorker(id int, jobs <-chan []string) {
	for {
		select {
		case artists, ok := <-jobs:
			if !ok {
				return
			}

			yandexLinker.SearchArtists(artists)
		}
	}
}
