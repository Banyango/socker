![logo](./site/logo.png)
# Socker

A lib for doing client/server state logic using websockets.

My game server is using WebRTC DataChannels for handling the data sent from the clients to the servers. This is a twist on WebRTC where a client opens a connection with the server and each sets up a single data channel peer connection through a websocket. This gets us a reliable UDP connection to the browser, which is great because websockets run over TCP and if we use them for game data we'd potentially run into issues. However WebRTC with a websocket setup requires a few back and forth steps based on the state of each side.

`Socker helps facilitate moving a server/client state forwards based on handlers.`


_Anyways this project is probably pretty specific to my needs but I thought I'd throw it up on github because why not._   

#### My Example 
___

|Server|Client|
|------|------|
|Setup websocket| |
| | Client connects to websocket
|Server sends list of buffered scene creation data| |
| |Client creates scene and data from server|
| | Client sends handshake to notify that it's ready to receive webrtc connection info.
|Server gets handshake and sends webrtc data | |
| | Client sets up webrtc connection and sends back webrtc info|
| Server sets remoteDescription from the client| |
| **Server listens for further websocket data** | **Client sends websocket data** |


___

#### Server.go

``` 
// Client connects and sends a message 
connection.Add(func(message []byte) bool {
    // send all buffered data (other players, networked objects, state, etc)                 
    return true
})

// Client sent the handshake for webrtc 
connection.Add(func(message []byte) bool {
    // Create the web rtc data channel
    // Create Offer
    // Send Offer to client.
    return true
})

// Client sent the answer
connection.Add(func(message []byte) bool {
    // Set RemoteDescription to answer.
    // data channel is created.
    return true
})

connection.Add(func(message []byte) bool {
    // Done setting up connection
    // Can use this for any client data that needs to make sure it gets to server. (Disconnect, Some kind of state)
    return false
})


 ```
 ___
 
#### Client.go

``` 
// Server sends buffered data
connection.Add(func(message []byte) bool {
    // Create all entities. 
    // Send webrtc connection request.                
    return true
})

// Server sends offer 
connection.Add(func(message []byte) bool {
    // set remotedescription to offer.
    // create answer.
    // send answer.
    return true
})

connection.Add(func(message []byte) bool {
    // handle server data messages that must reach client. (A player disconnects, joins, etc)
    return false
})
 ```
