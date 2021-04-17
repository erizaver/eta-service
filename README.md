# eta-service

Eta service returns closest carID and eta of this car to the point specified in request.

The Easiest way to launch application, is through docker-compose, just type

`make docker-compose-up`

To launch it locally you have to specify environment values location by adding 

`--local-config=*PATH_TO_LOCAL_FILE*`

or use default, for now its hardcoded as well
(default values contains redis host for docker-compose, so change it to localhost). 

This was made in order to be able to change environmental variables.

Or you can simply type 

`make run` 

to launch it locally.

Or run in containers separately

`make docker-build`

`make docker-run`

Also, you can run redis in a docker.

`redis-run`

But you will have to specify networking if you`ll launch them separately(just use compose=) ).

After that you can try examples in _example_ folder.

If you need to change something, you can fix proto file, and launch

`make generate`

Also, Makefile contains generators for mocks and external clients.