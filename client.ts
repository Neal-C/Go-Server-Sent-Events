const EVENT_SOURCE = new EventSource("http://localhost:4000/stream");

EVENT_SOURCE.addEventListener("message",(data) => {
    console.log("Awww we were sent data 💖")
    //? data is raw bytes ?
    //? data &&= data.toLocaleString() ?
    console.log(data);
});

EVENT_SOURCE.addEventListener("error", (error) => {
    throw ( "panic!");
});