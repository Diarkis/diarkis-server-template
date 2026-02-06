module github.com/Diarkis/diarkis-server-template/examples/csar/dgs

go 1.24
toolchain go1.24.0

require (
	github.com/Diarkis/diarkis v1.3.2
	github.com/magefile/mage v1.15.0
)

require golang.org/x/sys v0.33.0 // indirect

replace github.com/Diarkis/diarkis => /Users/Shared/Diarkis/diarkis
