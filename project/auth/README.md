# Authentication & Authorization

## Why do we cover Auths?
Authentication and authorization are absolutely important topics in web development(or any platform where user identification is necessary). Everything in our lives requires identification of who you are. Your bank account, your work data in your office PC, your social network account, etc. - all these are impossible to operate without indentifying it is you who are using them. You would hardly discover any application without **auths**(from now on we would like to refer to both of them all as *auths*, since these two usually go together).

Thus, for web developers authentication and authorization is unavoidable. No matter how complicated they are, it is necessary for developers to understand how they work and fix possible problems if necessary. Therefore, we naturally choose these topics as good exercises for writing server-side Go code.

## How do we cover Auths?
However, they are so important that tons of books have been already written for only these topics. Therefore, in this project we will only cover the very basic ideas of auths, and our focus will be mostly on how to implement them rather than deeply covering the theories behind.  

We tailored the topics as follows based on the progress of a reader, so that one can exercise one's knowledge leared from the book:
- [Basic](./basic/README.md)
- [Intermediate]()

The codebases of the two sections could overlap largely. On the one hand, we want continuity of the topics. But on the other handm we want to separate the code so that readers can choose which one to read once they finish either basic or intermediate level. So the completed code version will be in the [Intermediate]() section, and you can use it for reference later. 

## What do we cover for Auths?
There are so many sub-topics inside Auths, and obviously, we cannot cover all of them. However, our gola is to cover one of the most widely used auth protocol - [OAuth 2.0](https://oauth.net/2/). All those "Google Login" or "Apple Login" are based on this protocol. So it is definitely worth studying the protocol and getting our hands dirty.

But we might as well go over a few other simple schemes first, if necessary. 

# References
- [MDN Guide](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Authentication)
- [API Security in Action](https://www.amazon.com/API-Security-Action-Neil-Madden/dp/1617296023)
