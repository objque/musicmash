package subscriptions

func runAppleWorker(jobs <-chan []string) {
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

func runYandexWorker(jobs <-chan []string) {
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
