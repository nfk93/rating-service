CONTAINER_NAME='postgres_rating_service'

docker container kill $CONTAINER_NAME
docker container rm $CONTAINER_NAME
docker build -t pg_rating_service_test_image -f pg_docker_file ..
docker run --name $CONTAINER_NAME -d -p 5432:5432 pg_rating_service_test_image