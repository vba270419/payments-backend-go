# Payments Backend

An GO implementation of a simple payments service which supports the following operations:
- fetch a payment resource
- create, update abd delete a payment resource
- list a collection of payment resources
- persist resource state

Please refer to _API_Documentation.pdf_ to learn more about the payments API design.

## Running an application

1) Current implementation uses MongoDB for persisting payment resources, therefore the pre-requisite to run an application is to have a running instance of a MongoDB. For example, it can be a local MongoDB docker process:
    ```
    docker run -d -p 27017:27017 --name payments_mongodb mongo
    ```  
2) Build an application using the simple go command: 
    ```
    go build -o payments-server.go
    ```
3) To run an application with default parameters you can simply start the application built in the previous step:
    ```
    ./payments-server.go
    ```    
    or, you can pass custom MongoDB connection properties as well as some other configurable properties to an application using _--conf_ flag. As an argument a path to a json configuration file has to be provided. 
   ```
    ./payments-server.go --conf=<path_to_json_file>
   ``` 
   If _--conf_ flag is not set, application by default reads _config/server.json_ file.The full list of supported properties is shown in Implementation details(#implementation-details) section 

## Implementation details

| Property          | Description                   | Default value |
| --------          | -----------                   | ------ |
|**server_host**    |server TCP address to listen on|127.0.0.1|
|**server_port**    |server port number             |8000|  
|**server_timeout** |the maximum duration for reading and writing requests before http server times out (in seconds)|15|
|**mongodb_host**   |MongoDB instance host address|127.0.0.1|
|**mongodb_port**   |MongoDB instance port number|27017|  
|**mongodb_timeout**|the maximum duration for querying and persisting payment resources before MongoDB session times out (in seconds)|10|

## Out of the scope
Current implementation does **not** not support:
- user authentication and authorization
- secure MongoDB connection
- pagination and filtering in get all payment resources