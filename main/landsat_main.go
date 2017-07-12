package main

import "landsat"

func main() {
	scanChan := landsat.GenerateScanChan()
	landsat.ProcessScan(scanChan)
}
