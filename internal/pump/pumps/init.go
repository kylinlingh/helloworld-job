package pumps

var availablePumps map[string]PumpBackend

func init() {
	availablePumps = make(map[string]PumpBackend)

	availablePumps["csv"] = &CSVPump{}

}
