.PHONY: build-lambda clean 

# to build binaries to zip it and uploade as aws lambda resource 
build-lambda: 
	GOOS=linux go build -o ./build/safenotes ./cmd/main.go

pkg-lambda: 
	cd ./build && zip -r safenotes.zip safenotes

clean:
	cd ./build && rm -rf *