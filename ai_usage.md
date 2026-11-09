# AI Usage Report

## 1. How I Used AI

I used ChatGPT as a support tool while working on this assignment.

I mainly used it for:

- Understanding the assignment requirements and breaking them into smaller tasks.
- Understanding Go concepts such as HTTP handlers, JSON handling, maps and mutexes.
- Getting help with API implementation and error handling.
- Understanding the debugging issues in the provided Go program.
- Improving the documentation and making the explanations clear.

## 2. Suggestions I Changed or Adapted

One suggestion was to use a database for storing events. However, the assignment allowed in-memory storage, so I decided to use maps with a mutex to keep the implementation simple and focus on the main requirements.

I also adjusted the testing commands to use port 8080 because that was the port used by my application.

I did not blindly follow every suggestion. I checked the code and tested the changes before accepting them.

## 3. What AI Helped Me Understand

AI helped me understand why duplicate events should not be counted more than once.

I also learned why mutexes are important when multiple goroutines update shared data. The race detector helped me verify that the concurrency issue was fixed.

I understood how campaign statistics can be maintained and why large-scale systems may need message queues and background workers.

## 4. What I Generated and Then Tested or Debugged

I used AI-generated guidance while implementing parts of the Go backend and fixing the debugging program.

I tested the APIs using PowerShell requests and checked how the application handled valid events, duplicate events, invalid events and conflicting event IDs.

For the debugging task, I fixed four issues and compared the actual output with the provided expected output.

I also ran the Go race detector to check for concurrency problems.

AI helped me understand the problems, but I verified the changes by running the program and checking the results myself.