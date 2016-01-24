Pomodoro Crew
===============

Pomodoro Crew is an efficient Time Management tool intending to solve or minimize team communication and productivity problems.

[Presentation about Pomodoro Crew] (https://docs.google.com/presentation/d/1qNrgNp0SvOotwlRot-Yjs2EZp0lyZDd_FmhU-s-WV1c/edit?usp=sharing)


##Used tools

We choose [Cayle] (https://github.com/google/cayley) since it gives us flexibility and fits the best our need to store relations among different entities and [Bolt] (https://github.com/boltdb/bolt) as embedded storage engine as simple, fast and reliable solution. 

##Technical details

User signs in to the application and obtains JWT token which is used as a session identifier lateron. He can run and stop timer and a server will update all other followers and users within his teams.

##What was done:


* Authorization with JWT token in cookies

* Synchronization across browser windows and different user via websockets

* Users, teams and their relations are stored in graph database 

##Plans to do


* Registration

* Integration with Jira

* Slack integration for messaging and status updates

* Team customization

* Timer customization


