To build
>cd logic; go build -o gobra
To run
>cd $GIT_ROOT; logic/gobra

To test
>go run logic/*.go

To build image for formal release
>docker build --build-arg version=1.2.3 .

