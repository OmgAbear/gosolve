**This project has been developed on Ubuntu 20.04.4 LTS**

*It is not intended to run on any platform/chip*  

## Running the service:

- Install make if not available
  - Run ``make`` in root for a full list of available commands. 
## Run with docker compose is recommended https://docs.docker.com/get-docker/

  ``make docker-run``  
    
  The docker container port to host mapping is 8080:8080 for Go server and 8081:8081 for frontend.  
  You can configure the Go port of the server but this will not change the docker mapping.  
  The dockerfile for Go copies the entire project root and builds/runs the server there as is.  
  I did not use a multistage build.

**Without docker:**  
  If you have the required go version you can ``make run`` 
and this will start the go server locally.  
  Frontend can also be built locally if you have the right tools. It was tested locally with node v20.12.2 under nvm, built with esbuild, files were served with http-server.
## Structure:

- The Makefile contains a list of commands to build/start/test/clean artifacts on local machines.
It also contains commands to run the relevant servers in docker containers    

**Project structure**
    
    -root
        - cmd holds any available command mains - in this case, a basic http server
        - dev holds development relevant docs and tools
        - config holds the yaml configuration file
        - frontend holds the relevant frontend code. Normally would be a separate repository.
        - internal holds different layers
            - infrastructure holds the repository
            - http_interface holds the dtos, handlers and routing setup
            - config holds the internal structs and setup for configuration as well as logger
        - static holds the asset (input.txt, the file containing numbers)
        - Makefile - contains commands for running/testing/cleaning artefacts as well as starting the service(s) using docker compose.
**There are no reusable packages made available externally, thus, the use of internal**

## Tests ##
You can find a couple of tests in the infrastructure package for the purpose of displaying how I tend to tackle testing.
I generally prefer to limit scope of such testing to core/domain functionality while ensuring that other layers are tested in different stages.

## Validation ##
The only validation is for string to int conversion when calling the endpoint.
The path ```{value}``` is expected to fit into an int (platform dependant).

## Http Calls ##
There's a postman collection in dev folder.
