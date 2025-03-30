# dwimc

Simple device last known location service (Dude, where is my car?), De-wim-C in short.  
This service was designed for personal self hosted use.

## Motivation

I'm not driving on daily basis these days, remembering where I've parked my car and
searching for it tends to be time consuming.

## Common usage

View location: [Home Assistant](https://www.home-assistant.io/) app or its web interface, requires HA cloud or standalone deployment.  
Posting location: [Tasker for Android](https://play.google.com/store/apps/details?id=net.dinglisch.android.taskerm) or [MacroDroid - Device Automation](https://play.google.com/store/apps/details?id=com.arlosoft.macrodroid) with Bluetooth event trigger (detecting when car has turned off - Phone disconnects from the audio panel).

## Envs

Required envs passed to [Dockerfile](./Dockerfile):

```bash
-- DATABASE_URI - mongodb uri, default is local
-- DATABASE_NAME - db name
-- PORT - service port to listen on
-- SECRET_API_KEY - auth key which will be compared with requests X-API-Key header content.
```

More envs can be found in the: [env_example](./env_example) file

> **_NOTE:_** `SECRET_API_KEY` is a global auth key to all clients and not used to generate / validate client specific auth key. Basic / JWT authentication was not implemented yet.

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

```bash
docker-compose up --build
```

## Client API

> Note that `X-API-Key` header value must be the same as `SECRET_API_KEY` value.

### Post device location

Create a new device / updates an existing one by performing:

```bash
curl --location 'http://localhost:1337/api/devices/' \
   --header 'Content-Type: application/json' \
   --header 'X-API-Key: ••••••' \
   --data '{
      "serial": "serenity-123",
      "name": "serenity spacecraft 123"
   }'
```

Response:

```json
{
    "data": {
        "id": "67e97602e9621df49430c290",
        "created_at": "2025-03-30T16:49:06.453Z",
        "updated_at": "2025-03-30T16:49:20.337Z",
        "serial": "serenity-123",
        "name": "serenity spacecraft 123"
    },
    "error": null
}
```

Then post the actual location to this device:

```bash
curl --location 'http://localhost:1337/api/devices/67e97602e9621df49430c290/locations' \
--header 'Content-Type: application/json' \
--header 'X-API-Key: ••••••' \
--data '{
    "latitude": 32.179111,
    "longitude": 34.916111
}'
```

Response:

```json
{
    "data": {
        "success": true
    },
    "error": null
}
```

### Get last device location

Call:

```bash
curl --location 'http://localhost:1337/api/devices/67e97602e9621df49430c290/locations/latest' \
--header 'X-API-Key: ••••••'
```

Response:

```json
{
    "data": {
        "id": "67e977f5fc86793b73a0e161",
        "created_at": "2025-03-30T16:57:25.203Z",
        "updated_at": "2025-03-30T16:57:25.203Z",
        "device_id": "67e97602e9621df49430c290",
        "latitude": 32.179111,
        "longitude": 34.916111
    },
    "error": null
}
```
