.PHONY : all
all: house

BUILDOS=
ifeq ("$(OS)", "linux")
	BUILDOS= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
endif

CUR_PATH=${MY_REPO}/house/

clean:
	rm ./bin/house

house:
	${BUILDOS} go build -o ./bin/house ./src/main.go

dirs := $(shell ls ${CUR_PATH}/common/typedef/pb)
inject:
	$(foreach N,$(dirs),protoc-go-inject-tag -input=${CUR_PATH}/common/typedef/pb/${N} -XXX_skip=gorm;)
	# protoc-go-inject-tag -input=${CUR_PATH}/common/typedef/pb/*.pb.go -XXX_skip=gorm

test:
	go test -v ./src/scraping/

proto:
	protoc -I=${GOPATH}/src/github.com/golang/protobuf  -I=${CUR_PATH}/common/typedef/proto/ --go_out=${CUR_PATH} ${CUR_PATH}/common/typedef/proto/*.proto

protoall: proto inject
