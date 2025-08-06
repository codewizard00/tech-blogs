1. How Node js works well as a single threaded language to operate so many requests at a time.
2. What are bottlenecks of Node js while handling too many requests at once.


So before deep into how node js works well with thousands of requests at once.
Node js is single threaded non blocking runtime engine for javascript.
Node js Server Code. 
1. Create Socket 
2. Bind Socket with Port 
3. Start Listen On the socket 
4. Accept Connection


Why do we need to know this because this must raise lot of question in your mind like single threaded means using a single thread ohh if a request take 100ms the 10000 request will gonna take (10000*0.1)seconds of time 
which around 1000 seconds but this will never be the case untill you take help from any intern. 

So how it all works. To make such little time Node js architecture and our operating system play a major role init.

How, lets understands this in little bit depth of node js.

Once a request reach to server it goes to the sync queue where it wait for the ack. when the ack. recieves it went to accept queue from accept queue one more this all happening at the os level currently. If you want to know about this process more do read about the syn flood attack.
After receiving to accept queue a notification is send via epoll from os to our code (libuv) that one request want to something from your code. Libuv took that request and send as callback to callback stack then all process follows. Now those task which are asynchronous in nature will handle by node js or libuv and libuv has worker threads such that the main thread will not get blocked and start accpet new connection.

At the end we know that as being single threaded uses task delegation to handle large number of request once.

Now i have question how it maintains order like concurrent requests taken how did node js know whose output is this or not.
And after know the output how it gives back to the accept queue or whereever i need to send the response i want to know that process in depth

Each invocation of code has its own for res object, what takes the identity of which socket fd it belongs to.
The res object is responsible to write the response to the correct socket fd in accept queue.


What new i learn today 
TCP 5 tuple 
Source Port, Source IP, Destination, Destination Port, Protocol

