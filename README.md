# Feedback service

## Task Description

Your task is to create a simple microservice for processing customer feedback data. The microservice will be responsible for managing customer feedback, storing it in a PostgreSQL database, and publishing it to a Kafka topic. Additionally, the microservice implements caching using Memcached to improve performance.

### Requirements

1. Create a Go microservice with the following endpoints:

   * POST /feedback: allows the user to submit a new feedback item. The endpoint should accept JSON in the following format:

   ```json
   {
      "customer_name": "Andrii Derkach",
      "email": "andrsj.derkach@gmail.com",
      "feedback_text": "It's a really good test task.",
      "source": "https://t.me/andrsj"
   }
   ```

   > The microservice should validate the input data, create a unique identifier for the feedback item, save the item in the PostgreSQL database, and publish it to the Kafka topic.

   * GET /feedback/{id}: allows the user to retrieve a specific feedback item by its identifier. The microservice should first check the Memcached cache for the item. If the item is not in the cache, it should retrieve it from the PostgreSQL database and store it in the cache for future requests.

2. Use the following technology stack: (if needed for the task):
   * Go
   * PostgreSQL for data storage
   * Kafka for message publishing
   * Memcached for caching
   * JavaScript and React for the interface
   * Microservice communication via HTTP
   * JWT for authentication

3. The microservice should return appropriate error messages if the endpoint is used incorrectly or if there are issues with the data.

4. You should use any relevant libraries or frameworks that are convenient for you.

### Bonus points

* Implement pagination for the /feedback endpoint to limit the number of results returned.
* Add unit tests for the microservice endpoints.
* Implement authentication and authorization using JWT.
* Use Docker to containerize the microservice.

### Submission

Please provide the source code for your solution, including instructions for configuration and running the program.

**Good luck!**

## Solution

### Endpoints

* `GET /` - Status code 200 for verify the viability of server

![Status output](/img/status.png)

---

* `GET /token?minutes=10&role=all` - Generator of JSON Web Tokens
  * minutes:
    * int
    * default = 10
  * role:
    * string
    * available: `get`, `post`, `all`

Text | Image
---- | -----
Input for generating JWT | ![Token input](/img/tokenInput.png)
JWT Response | ![Token result](/img/token.png)
Invalid value | `{"error": "error while checking minutes: wrong value for minutes param '-500': invalid minutes parameter"}`
Invalid role | `{"error": "error while checking role: wrong role 'none': invalid role parameter"}`

---
JWT Errors:

Errors | Messages
------ | -----
Invalid JWT | `{"error": "wrong token authorization: parsing error: token is malformed: could not JSON decode header: invalid character 'ÿ' looking for beginning of value"}`
Wrong access role | `{"error": "validating token error: validating 'role' error: wrong role for 'POST': token has wrong role"}`
Token Expired | `{"error": "validating token error: validating 'expiredAt' error: expired (122): token is expired"}`

---

* `GET /feedbacks` - GET all feedbacks from DB [No cache, no paginated values]

Text | Image
---- | -----
Response for GET | ![Response for GET /feedbacks](/img/GETallFeedbacks.png)

---

* `GET /feedback/{id}` - GET one specific feedback by ID
  * id - UUID string (that parsed into `uuid.UUID`)

Text | Image
---- | -----
Input/Output | ![GET Headers input](/img/GETidFeedback.png)
Output/Headers no cached | ![GET Header output NoCache](/img/GETidFeedback2.png)
Output/Headers cached | ![GET Header output WithCache](/img/GETidFeedback3.png)
Wrong format ID | `{"error": "can't parse the ID: invalid UUID length: 35"}`
Invalid character | `{"error": "can't parse the ID: invalid UUID format"}`
Missing ID | `{"error": "getting by ID: failed to get feedback from DB: record not found"}`

---

* `GET /p-feedbacks?limit=10&next=""` - Paginated version of `/feedbacks`
  * limit:
    * int
    * EXPECTED > 0
  * next:
    * string
    * **NEED TO BE IN UUID format**

Text | Image
---- | -----
Output| ![Response](/img/GETpFeedbacks.png)
Headers | ![Response headers](/img/GETpFeedbacks2.png)
Next request/response | ![Middle request/response](/img/GETpFeedbacks3.png)

Error | Message
----- | -------
Wrong type of limit | `{"error":"error while check limit: wrong limit param 'a': invalid limit parameter"}`
Wrong limit value | `{"error":"error while check limit: wrong limit param '-1 < 0': invalid limit parameter"}`
Wrong format of next | `{"error":"error while check next: wrong format of next ID: invalid next parameter"}`
No values after next | `{"error": "next 'b7970cbb-0a70-41bb-99d7-52ec3884d3d5': no values after 'next'"}`

---

* `POST /feedback` - CREATE one feedback

Text | Image
---- | -----
Body | ![Body of request](/img/POSTinput.png)
Headers | ![Headers of request](/img/POSTinput2.png)
Output | ![Response](/img/POSToutput.png)
Invalid email (no @ for example) | `{"error": "validating feedback error: invalid email address"}`
Invalid source URL (no protocol for example) | `{"error": "validating feedback error: invalid source URL"}`

---

* `/l?time=0` - Just a handler that sleep, used to test "graceful shutdown" (actually a timeout signal interrupt)

## How to run?

In the [Makefile](/Makefile) I include a lot of different commands:

* `make build` - build app on local machine
* `make brun` - build and run on local machine
* `make db-up | kafka-up | cache-up` - app separate service
* `make consumer` - CLI consumer for TOPIC that used for producing in App
* `make up` - Start whole dockerize project
* `make down` - Stop it
* `make go` - rerun Go app container without rebuild
* `make bgo` - rebuild and run Docker image for Go app
* e.t.c.

## Suggestions for improvement

1. Unit tests! Integration tests! E2E tests! No manual testing!
2. More effective error handling [no wrapped messages]
3. Communicating between running requests for graceful shutdown
4. Make writing to DB and Kafka in goro [not 100% sure]
5. Move generating UUID for models inside DB and return this value back to APP [Currently the ID is generated in Repository]

## Issues

If we are running the Docker Compose, we can see that the Go Web Server try to connect to not-ready DB.
I tried to fix it with [this shell/bash script](https://github.com/vishnubob/wait-for-it/) but not successful.
