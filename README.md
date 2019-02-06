# Wave Messaging Management Microservice

[![Maintainability](https://api.codeclimate.com/v1/badges/cab714ac578daa72645c/maintainability)](https://codeclimate.com/github/WAVE-CBSSPOIMME/wave-messaging-management-service/maintainability)

Wave Instant Messaging Management Microservice aims to provide a secure ACL management interface on top of VerneMQ MQTT broker.
The system is designed with integration in mind in order to provide developers with an easy to integrate messaging service in their applications.

## Table of Contents
- [Wave Messaging Management Microservice](#Wave-messaging-management-microservice)
    - [Table of Contents](#table-of-contents)
    - [Config](#config)
    - [External/Internal Mapping](#externalinternal-mapping)
        - [Redis Stores](#redis-stores)
    - [Authentication & Authorization](#authentication--authorization)
        - [Authentication](#authentication)
            - [External Authentication Endpoint](#external-authentication-endpoint)
            - [MQTT Authentication](#mqtt-authentication)
            - [VerneMQ ACL](#vernemq-acl)
        - [Authorization](#authorization)
            - [Private Conversations](#private-conversations)
            - [Group Conversations](#group-conversations)

## Config

Wave requires these environment variables to be set :

|              Name             |                          Description                          |
|:-----------------------------:|:-------------------------------------------------------------:|
|   WAVE_CONFIG_FILE_PATH     |           Absolute path to your config.json file              |

Config file must have the following structure : 

```json
{
    "authenticationCheckEndpoint": "https://www.myapp.com/myexternalauthendpoint",
    "tokenValidationRegex": "mytokenregex"
}
```

|              Field            |                          Description                          |
|:-----------------------------:|:-------------------------------------------------------------:|
|   authenticationCheckEndpoint | External authentication endpoint provided by your application |
|   tokenValidationRegex        |           Token format validation regular expression          |

## External/Internal Mapping

In order to be able to accept any kind of authentication system (JSON Web Token, Sessions, ...) we decided to map your application user identifiers and tokens with our own internal structures.

The token remain preserved but application user identifiers are mapped with an internal Wave user ID
under the form of an UUID V4. This allows us to easily map one user to its according topic.

These mappings are stored in a Redis instance.

### Redis Stores

|    Type   |            Key           |                           Value                           |
|:---------:|:------------------------:|:---------------------------------------------------------:|
| Key-Value |      session:{token}     |                   {internalWaveUserID}                  |
|    Hash   | mapping:{originalUserID} | token {token} internalWaveUserID {internalWaveUserID} |

## Authentication  & Authorization

>**Authentication** is the process of ascertaining
> that somebody really is who he claims to be while
>**Authorization** refers to rules that determine 
>who is allowed to do what. E.g. Matteo may be 
>authorized to publish and subscribe to a MQTT topic, while Terry is only authorised to publish on it.

### Authentication 

Wave is designed in a way that its system doesn't need access to existing project databases, that for obvious security and integration reasons. 

The identity **verification process must then be managed by your application side**.

#### External Authentication Endpoint
Each request emitted to this service must contain the authentication token of your application user as an HTTP header (See API Documentation section). 

The service will then verify the token authenticity by calling an external endpoint (Provided by your application) responsible to verify the received token and send back informations. 

Internally, Wave system will keep an internal mapping in a Redis instance matching your application users IDs and tokens with our internal identifiers and password.

Wave requires that your web server endpoint respond to requests in a certain predefined format. Our system will issue a `GET` request with a `token` HTTP header containing your application user token and an `empty body`.

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

|   MQTT   |        Wave        |
|:--------:|:--------------------:|
| clientID | internalWaveUserID |
| username | internalWaveUserID |
| password |         token        |


#### VerneMQ ACL

VerneMQ ACLs are stored in a MongoDB Collection named `vmq_acl_auth` with the following schema : 

```
{
	"_id" : ObjectId("5bc897e0a9e22f4273dd067a"),
	"mountpoint" : "",
	"client_id" : "cff1c5b7-9508-49fa-af8a-a4009ac5f27f",
	"username" : "cff1c5b7-9508-49fa-af8a-a4009ac5f27f",
	"passhash" : "$2a$14$.q06pmH.XAlvIDuV9SsIFOSqI/Zf6OE3SfC4r4cm3uwc05cC0i70G",
	"publish_acl" : [
		{
			"pattern" : "conversations/private/cff1c5b7-9508-49fa-af8a-a4009ac5f27f/+"
		}
	],
	"subscribe_acl" : [
		{
			"pattern" : "conversations/private/+/cff1c5b7-9508-49fa-af8a-a4009ac5f27f"
		}
	]
}
```
Note `passhash` field is a [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) hash of the token.

Unfortunately, MQTT doesn't provide any way to identify the sender of a message in a trustful manner. 

To circumvent that, we designed a topic hierarchy granting us the ability to trustfully identify the sender of a message. Even if MQTT topic wildcards are used by the user to subscribe topics, a user will always receive the precise topic that forwarded the message. 

This brings the user the ability to know exactly who's the sender of a message.

### Authorization

#### Private Conversations

Each user gets assigned a sender topic prefix path under the form of `conversations/private/{internalWaveUserID}`. 

In order for the subscriber to be able to trust the sender of a message a user can only publish on `conversations/private/{internalWaveUserID}/+` topic wildcard. 

Each user will then subscribe the `conversations/private/+/{internalWaveUserID}` topic wildcard in order to be able to receive message. 

In MQTT terms this means that on creation the user gets assigned the following ACLs :

|              Topic                             | Publish |     Subscribe      |
|:----------------------------------------------:|:-------:|:------------------:|
| conversations/private/{internalWaveuserID}/+ |    ✅    |   ✅ <sup>1</sup>  |
| conversations/private/+/{internalWaveuserID} |    ❌    |   ✅               |

<sup>1</sup> _Implicit due to wildcard subscription._


#### Group Conversations

Each group conversation is assigned a topic with path `conversations/group/{groupID}`. 

In order for the subscriber to be able to trust the sender of a message a user can only publish on `conversations/group/{groupID}/{internalWaveUserID}` topic. 

Then each group members will have to subscribe the `conversations/group/{groupID}/+`  topic wildcard in order to receive messages from all members.

This means the following MQTT ACLs are created for each user on group creation :
  
|              Topic                                     | Publish | Subscribe            |
|:------------------------------------------------------:|:-------:|:--------------------:|
| conversations/group/{groupID}/{internalWaveuserID}/+ |    ✅    |     ✅ <sup>1</sup>  |
| conversations/group/{groupID}/+                        |    ❌     |    ✅               |

<sup>1</sup> _Implicit due to wildcard subscription._

Note that you'll have to discard messages client side as one user sending a message will also receive its own message due to this configuration. 
