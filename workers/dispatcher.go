package workers

func StartDispatcher(nworkers int) {
	startDeltaDispatcher(nworkers)
	startFileDispatcher(nworkers)
}
