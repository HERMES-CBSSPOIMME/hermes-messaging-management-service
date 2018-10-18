# Hermes Messaging Management Microservice

HERMES Instant Messaging Management Microservice aims to provide a secure ACL management interface on top of VerneMQ MQTT broker.
The system is designed with integration in mind in order to provide the developer with an easy to integrate a messaging service in his/her applications.




## Authentication  & Authorization

>**Authentication** is the process of ascertaining
> that somebody really is who he claims to be while
>**Authorization** refers to rules that determine 
>who is allowed to do what. E.g. Matteo may be 
>authorized to publish and subscribe to a MQTT topic, while Terry is only authorised to publish it.

### Authentication 

As HERMES is designed in a way that its system doesn't need access to existing project databases, that for obvious security and integration reasons. 

The identity **verification process must then be managed by your application side**.

#### External Authentication Endpoint
Each request emitted to this service must contain the authentication token of your application user as an HTTP header (See API Documentation section). 

The service will then verify the token authenticity by calling an external endpoint (Provided by your application) responsible to verify the received token and send back informations. 

Internally, HERMES system will keep an internal mapping in a Redis instance matching your application users IDs and tokens with our internal identifiers and password.

HERMES requires that your web server endpoint respond to requests in a certain predefined format. Our system will issue a `GET` request with a `token` HTTP header containing your application user token and an `empty body`.

Your endpoint should respond the following :

**Invalid Token**

```
HTTP/1.1 400 Bad request
Content-Type: application/json
Date: Thu, 18 Oct 2018 10:35:54 GMT
Content-Length: 0
```

**Valid Token**

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Oct 2018 10:35:54 GMT
Content-Type: application/json
Content-Length: 34

{"userID":"put_the_userID_here"}
```

#### MQTT Authentication

At the MQTT level each user credentials are represented with the following mapping :

|   MQTT   |        HERMES        |
|:--------:|:--------------------:|
| clientID | internalHermesUserID |
| username | internalHermesUserID |
| password |         token        |

VerneMQ ACLs are stored in a MongoDB Collection named `vmq_auth_acl` with the following schema : 

```json
{
	"_id" : "5bc8541703e034b6be3b8870",
	"mountpoint" : "",
	"client_id" : "3bf186ba-ede0-4e1e-8305-64767942229a",
	"username" : "3bf186ba-ede0-4e1e-8305-64767942229a",
	"passhash" : "$2a$14$uBvhPcRP2NGo4eMm.nAIgu6jIZXUMe9DZ9zdLvK03xGypzvBYZvYS",
	"publish_acl" : [
		"conversations/private/+"
	],
	"subscribe_acl" : [
		"conversations/private/3bf186ba-ede0-4e1e-8305-64767942229a"
	]
}
```
Note `passhash` field is a [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) hash of the token.

### Authorization

MQTT Topics authorization are managed the following way :

#### Private Conversations
Each user gets assigned a topic with path `conversations/private/{internalHermesUserID}`. Only this user can subscribe the topic but everyone can publish in it. In MQTT terms this means that on creation the user gets assigned the following ACLs :

 |              Topic                          | Publish | Subscribe |
|:--------------------------------------------:|:-------:|:---------:|
| conversations/private/{internalHermesuserID} |    ❌    |     ✅    |
|     conversations/private/+                  |    ✅    |     ❌    |


#### Group Conversations

Each group conversation is assigned a topic with path `conversations/group/{groupID}` that members can publish and subscribe to.

