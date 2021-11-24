Install swag
>go get -u github.com/swaggo/swag/cmd/swag

Generate swagger
>cd $GIT_ROOT; swag init -g logic/main.go

To build
>cd $GIT_ROOT/logic; go build -o gobra
To run
>cd $GIT_ROOT; logic/gobra

To test
>export  AWS_PROFILE=<profile>
>go run logic/*.go

To build image for formal release
>docker build
      --build-arg version=1.2.4
      --build-arg key=$AWS_SECRET_ACCESS_KEY
      --build-arg id=$AWS_ACCESS_KEY_ID
      -t $CONTAINER_LATEST_IMAGE .