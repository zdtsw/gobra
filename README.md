To build
>cd logic; go build -o gobra
To run
>cd $GIT_ROOT; logic/gobra

To test
>go run logic/*.go

To build image for formal release
>docker build
      --build-arg version=1.2.3
      --build-arg key=$AWS_SECRET_ACCESS_KEY
      --build-arg id=$AWS_ACCESS_KEY_ID
      -t $CONTAINER_LATEST_IMAGE .