# consumer
A golang rabbitmq consumer
 
This project (consumer) and another project of mine (producer) are 2 sides of same coin.

- The consumer is forever waiting for messages on specific queue, 
  and does something.

To make this project work locally run 2 docker commands to start redis & rabbit:
1) docker run -d -p 15672:15672 -p 5672:5672 -p 5671:5671 --hostname my-rabbitmq --name my-rabbitmq-container rabbitmq
2) docker run --name my-redis-container -p 6379:6379 -d redis

Enjoy the code and please, feel free to do pull requests to the code.

### Upcoming Features:
* Retry connection to rabbit.
* Thread safe.

### Whats new?
* Go modules.
* Decomposition to saperate files. 
* health check Endpoint.