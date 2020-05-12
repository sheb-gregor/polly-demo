## A simple pub-sub 

The assignment needs to cover: 

- Source code control usage (github, etc); 

- Implementation: 

  - code structure, style, and organization; 
  - code annotations; 

- Testing; 

- Minimal project instructions (readme). 

Code <u>*performance, efficient memory usage, efficient algorithm*</u> should be the top guiding factors in designing the solution. If either is not addressed, explain the current limitations and how exactly those should be tackled (i.e. data structures, algorithms, etc). 

The assignment should not take more than a few hours at max. Try to timebox it to 2 hours; if you feel something will take too much time, just explain what/how would you do if you had enough time in the code/project description. 

Use a publicly accessible git service (github, bitbucket, etc) for the project source code. Feel free to use mostly any third-party modules you deem needed but please do not use something which implements pub-sub for you (for instance, redis). 

The project instructions should explain how to run the code, how it was tested, what was omitted due to lack of time. 

The goal is build a subsystem (go package) for a simple pub-sub polling service, which will be usedwithanHTTP-basedinterface.O nlythepub/subcomponent/packageisascopeofthis assignment; knowing it will be used through HTTP should shape the design decisions (like which HTTP methods are most appropriate for which interface methods). 

  The following are concepts and fundamental functions of this pub-sub service: 

- It allows to subscribe to and publish to topics; 
- A topic is identified by name; 
- Publishing to a topic notifies all current subscribers (when subscriber polls); 
- Subscribing to a topic does not deliver previously published messages; 
- Subscribing to a non-existing topic should create one; 
- Messages are stored in-memory; 
- The number of messages stored should not be limited. 

The assignment itself is about designing a package, which implements the logic without HTTP bindings. 

The interface of the module exposes the following API: 

- For publishers: 
	1. `publish(topicName, jsonBody)` - always success; 

- For subscribers: 
	1. `subscribe(topicName, subscriberName)` — always success;  
	2. `unsubscribe(topicName, subscriberName)` — always success; 
	3. `poll(topicName, subscriberName)` — returns either next unseen message or no message if all are seen or error if no subscription found. 

A published message stays in memory until all current subscribers have polled it or until all current subscribers have unsubscribed. 

Finally, provide the answer to the following questions: 

- What is the message publish algorithm complexity in big-O notation? 
- What is the message poll algorithm complexity in big-O notation? 
- What is the memory (space) complexity in big-O notation for the algorithm? 
