.PHONY: clean
clean:  ## Remove build artifacts
	rm -rf build erc20 node_modules

.PHONY: bindings
bindings:  ## Create the Go bindings
	yarn install
	yarn truffle compile
	yarn truffle run abigen Erc20
	mkdir -p erc20
	abigen --bin=abigenBindings/bin/Erc20.bin --abi=abigenBindings/abi/Erc20.abi --pkg=erc20 --out=erc20/erc20.go
