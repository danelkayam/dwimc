# dwimc
Simple device last known location service (Dude, where is my car?), De-wim-C in short.  
This service was designed for personal use.

## Motivation
I'm not driving on daily basis these day, remembering where I've parked my car and
searching for it tends to be time consuming.

## Common usage

View location: [Home Assistant](https://www.home-assistant.io/) app or its web interface, requires HA cloud or standalone deployment.
Posting location: [Tasker for Android](https://play.google.com/store/apps/details?id=net.dinglisch.android.taskerm) with Bluetooth event trigger (detecting when car has turned off - Phone disconnects from the audio panel).

## Envs

Required envs passed to [Dockerfile](./Dockerfile):
```
-- DATABASE_URI - mongodb uri, default is local
-- DATABASE_NAME - db name
-- PORT - service port to listen on
-- SECRET_API_KEY - auth key which will be compared with requests X-API-Key header content.
```

> **_NOTE:_** `SECRET_API_KEY` is a global auth key to all clients and not used to generate / validate client specific auth key. Basic / GWT authentication was not implemented yet.


## Local running

Generate random `SECRET_API_KEY`:
```bash
cat /dev/urandom | tr -cd '[:alnum:]' | fold -w ${1:-24} | head -n 1
```

Create local `.env` file:
```bash
SECRET_API_KEY= <---- place value here
```

Build and lift:
```
docker-compose up --build
```


## Client API

> Note that `X-API-Key` header value must be the same as `SECRET_API_KEY` value.

### Upsert
Inserting new device / updating existing one by performing:
```bash
curl -X POST http://localhost:1337/api/devices \
   -H 'Content-Type: application/json' \
   -H 'X-API-Key: secretapikey' \
   -d '{"serial":"serenity:123","name":"dan spaceship - serenity","position":{"latitude":32.0744615,"longitude":34.7911511}}'
```

### Get All
Retriving all devices:
```bash
curl http://localhost:1337/api/devices \
   -H 'Content-Type: application/json' \
   -H 'X-API-Key: secretapikey'
```

### Get One
Retrive single device:
```bash
curl http://localhost:1337/api/devices/serenity:123 \
   -H 'Content-Type: application/json' \
   -H 'X-API-Key: secretapikey'
```

