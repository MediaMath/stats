.PHONY:	test golint

# Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

test: golint
	go get ./...
	golint ./...
	go vet ./...
	go test -race $(TEST_VERBOSITY) ./...

golint:
	go get github.com/golang/lint/golint

