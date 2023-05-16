# Benchmarking 

## Prerequesities

- docker desktop

## Steps

1. Build docker file (run the command in the directory of this readme)
`docker build -t nubes-bench:latest  -f ./Dockerfile ../../`
2. Run docker file - choose preferable option:
 i. Detached mode (run the next commands in a terminal available in docker GUI or using docker exec )
    `docker build -dit nubes-bench:latest`
 ii. Run commands in the same terminal
    `docker build -it nubes-bench:latest`
3. [in docker container]: Configure aws account
   `aws configure`
   The command will ask for the access keys, generate the them here:
   ![aws configure screenshot](.//../../images/aws-configure.png)
   Note that the credentials must belong to a user with    **AmazonDynamoDBFullAccess**, **AWSCloudFormationFullAccess**, **AmazonS3FullAccess** and **AWSLambda_FullAccess** permissions to the user.
4. [in docker container]: Move to the scripts folder copied to the docker container
   `cd /nubes/evaluation/scripts`
5.  [in docker container, this directory]: Initialize & seed DynamoDB
   `./db_init.sh`
6.  [in docker container, this directory]: Deploy lambda functions
   `./deploy.sh`
7. [in docker container, this directory]:As the last step before benchmarking, assign URL to *Gateway* and *gateway_baseline* lambda functions on aws.
Fill the corresponding *gateway* variables in the corresponding lua scripts in this directory (*hotel.lua* and *hotel_baseline.lua*)

> `wget https://raw.githubusercontent.com/tiye/json-lua/main/JSON.lua`

8. [in docker container, this directory]: Invoke the wrk2:
   For the baseline:
 `wrk2 -R<request_rate> -d<duration_of_test_seconds>s -s hotel_baseline.lua <URL_to_gateway_baseline>`
 For the nubes:
 `wrk2 -R<request_rate> -d<duration_of_test_seconds>s -s hotel_baseline.lua <URL_to_Gateway>`