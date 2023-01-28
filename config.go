package binit

type Config struct {
	WORKDIR         string
	SKIP_SIGNAL_LOG string
	BEFORE          string
	AFTER           string
	PRE_STOP_SIGNAL string
	PRE_STOP        string
}
