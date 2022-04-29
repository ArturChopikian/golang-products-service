# grpc-server

Need to write a gRPC server in Golang (1.13+), with persistent storage, implements 2 methods:
- Fetch(URL) - request an external CSV file with a list of products at an external address. 
The CSV file looks like PRODUCT NAME, PRICE. The last price of each product must be stored 
in the database with the date of the request. You also need to save a lot of changes 
in the price of the product.
- List(paging params,sorting params) - get a paginated list of products with their prices,
the number of price changes and the dates of their last update. Provide all sorting options 
to implement the interface in the form of an infinite scroll.

For start project need to define the list of env variable:
>SERVER_HOST="localhost"<br>
>SERVER_PORT="50051"<br>
>SERVER_NETWORK="tcp"<br>
>CSV_SERVER_HOST="localhost"<br>
>CSV_SERVER_PORT="8090"<br>
>CSV_SERVER_FOLDER="files"<br>
>MONGODB_USER="root"<br>
>MONGODB_PASSWORD="root"<br>
>MONGODB_URI="url to mongoDB server"<br>
>LOG_PREFIX="server"<br>
