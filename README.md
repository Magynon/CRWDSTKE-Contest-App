# CRWDSTKE-Contest-App

GO API that allows you to interact with a PostgreSQL database.

It features both simple and batch queries using multithreading.
GET, POST, PATCH, DELETE requests are supported.

````
           _________                   _________
 O        |         |  single / batch |         |
\|/  ---> |   API   | --------------> |  STORE  | ---> PostgreSQL
 |        |_________| get/post/patch/ |_________|
/ \                      delete
````

