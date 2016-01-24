Pomodoro Crew project
===============

Our project concentrate on team communication and team pruductivity problem.
Pomodoro Crew its efficient Time Management tool.
Presentation about Pomodoro Crew https://docs.google.com/presentation/d/1qNrgNp0SvOotwlRot-Yjs2EZp0lyZDd_FmhU-s-WV1c/edit?usp=sharing


Used tools
As db we choose Cayley https://github.com/google/cayley and Bol https://github.com/boltdb/boltt solution because we were interested in using graph database, both are written in Go. Using Cayley gave us possibility to display progress of users from different groups. 
#Technical details

User sign in to our application and after it receiving JWT token into cookies which using for session. User could run and stop timer and websockets will update all users with whom he connected by teams.

What was done:
Login for user
Authorization with JWT token in cookies
Running websockets with streaming current stage time
Plans to do
Registration
Integration with Jira
Team customization
Settings page with customization of stages
Messeging and integration with messages
