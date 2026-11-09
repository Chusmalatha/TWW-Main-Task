# Campaign Events – Engineering Assignment

## 1. My Understanding of the Problem

The main goal of this project is to build a backend application that receives events from different delivery providers and provides campaign statistics.

For example, a company sends marketing emails to its customers. After sending the emails, the delivery provider informs our application about what happened to each email.

There are four types of events:

- Sent – The message was sent.
- Delivered – The message reached the customer.
- Opened – The customer opened the message.
- Clicked – The customer clicked a link in the message.

Our application should receive these events, store them in a database and provide the total counts for each event type through an API.

It should also handle duplicate events, invalid data and events that arrive in the wrong order.

## 2. Assumptions

While working on this assignment, I am making the following assumptions:

- Every event has a unique event ID provided by the delivery provider.
- I will use the event ID to identify duplicate events.
- Only four event types are allowed: sent, delivered, opened and clicked.
- Event ID, campaign ID, contact ID, type and timestamp are required fields.
- Metadata is optional.
- If the same event is received multiple times, it should be stored and counted only once.
- If a batch contains both valid and invalid events, I will process the valid events and reject the invalid ones.
- Events may arrive in any order, so I will not depend on their arrival order.
- Opened and clicked counts will represent the total number of unique events, not the number of unique customers.

## 3. Things That Are Not Clear in the Requirements

While going through the assignment, I noticed a few things that are not clearly mentioned:

- If the same event comes more than once, should we simply ignore it or return a separate message saying it is a duplicate?
- If a request contains 10 events and one of them has incorrect data, should we reject all 10 or accept the valid ones?
- If a customer opens the same email multiple times, should we count every open or count that customer only once?
- What should we do if two events have the same event ID but contain different information?
- Should the campaign statistics show the complete history or only the events from a particular time period?
- Is there any limit on how many events a provider can send in one request?

## 4. Questions I Would Like to Ask the Product Manager

Before starting the development, I would like to clarify a few things with the product manager:

- How do you expect us to handle duplicate events when the provider sends them again?
- Do you want to know the total number of times customers opened an email, or only how many customers opened it?
- If some events contain incorrect or missing information, should we process the remaining valid events?
- What should happen if a provider sends an event with an event ID that we have already received but with different details?
- How many events can we expect in a single request?
- Should the campaign statistics be updated immediately whenever a new event is received?

## 5. What I Plan to Prioritise

Since the assignment has a limited time, I will focus on the most important requirements first.

My implementation order will be:

1. Understand the provided Go starter project.
2. Set up the database to store campaign events.
3. Create the event structure and validate incoming data.
4. Implement the POST /events API.
5. Handle duplicate events properly.
6. Implement the GET /campaigns/{campaign_id}/stats API.
7. Test the application with different types of input, including invalid and duplicate events.
8. Work on the debugging task and fix the identified issues.

I will focus on making the main APIs work correctly before spending time on optional features.

## 6. Scaling the System to 100 Million Events Per Day

If the application needs to handle 100 million events every day, the current in-memory approach will not be enough. We need to improve the system to handle a large number of events.

### 6.1 Main Challenges

Some problems I expect when handling this much data are:

- **High traffic:** Many events may arrive at the same time, which can make the API slow.
- **Database performance:** Storing and retrieving millions of events can take more time.
- **Duplicate events:** Providers may send the same event multiple times.
- **Concurrent updates:** Multiple workers may update the same campaign statistics at the same time.
- **Storage:** A large amount of storage will be required to keep event records.

### 6.2 Changes I Would Make

To handle these problems, I would make the following changes:

- Use multiple API instances behind a load balancer to handle incoming requests.
- Use a message queue like Kafka to manage incoming events.
- Process events in the background using multiple workers.
- Use a scalable database instead of storing everything in memory.
- Add database indexes on `event_id`, `campaign_id` and `timestamp` to make searching faster.
- Maintain separate campaign counters instead of calculating statistics from all events every time.

### 6.3 Why Use a Message Queue?

A message queue like Kafka can help manage a large number of events.

The API can receive events and send them to the queue. Background workers can process them and store them in the database.

**Advantages:**
- Can handle large amounts of incoming events.
- Helps manage sudden increases in traffic.
- Allows background processing.
- Keeps the API from doing all the work itself.

**Disadvantages:**
- Makes the system more complex.
- Requires additional monitoring and maintenance.
- Statistics may take a little time to update because events are processed in the background.

### 6.4 Handling Duplicate Events

I would use `event_id` as a unique identifier for every event.

A database unique constraint can prevent the same event from being stored more than once. This is important because providers may retry sending events when they do not receive a successful response.

### 6.5 Maintaining Campaign Statistics

Instead of calculating statistics from millions of event records every time, I would maintain separate counters for each campaign.

Whenever a new valid event is processed, the corresponding counter will be updated.

For example:
- A sent event increases the sent count.
- A delivered event increases the delivered count.
- An opened event increases the opened count.
- A clicked event increases the clicked count.

These updates should be handled safely to avoid incorrect counts when multiple workers are processing events.

### 6.6 Monitoring

I would monitor the following things:

- Number of events received per second.
- Number of successfully processed events.
- Number of duplicate and rejected events.
- Queue size and processing delay.
- API response time.
- Database performance and errors.

This will help us identify problems and improve the system.

### 6.7 What I Would Not Build Initially

I would not add too many complex technologies at the beginning.

I would avoid unnecessary microservices, multiple databases and complicated caching unless they are actually needed.

My first priority would be reliable event processing, correct statistics, duplicate handling and monitoring.

**Conclusion:** To handle 100 million events per day, I would improve the current system using multiple API instances, a message queue, background workers and a scalable database. My main focus would be reliability, performance and accurate campaign statistics.

---

## 7. Marketer Scenario

The marketer noticed that the campaign statistics changed between 10:00 AM and 10:30 AM.

At 10:00 AM:
- Sent: 1,000,000
- Delivered: 970,000
- Opened: 250,000
- Clicked: 20,000

At 10:30 AM:
- Sent: 1,000,000
- Delivered: 975,000
- Opened: 248,000
- Clicked: 20,000

### 7.1 Why Did the Opened Count Decrease?

The opened count decreased from 250,000 to 248,000.

Normally, if we are counting all unique open events, the count should not decrease when new events arrive.

Some possible reasons are:
- The statistics may have been recalculated after removing duplicate or invalid events.
- Some previously counted events may have been corrected or removed.
- There may be a bug in the statistics calculation.
- The system may be using a different counting method, such as unique contacts instead of total open events.

### 7.2 What Would I Check First?

First, I would check the raw event records received between 10:00 AM and 10:30 AM.

Then I would check:
- Whether any previously counted events were removed.
- Whether duplicate events were handled correctly.
- Whether the opened count is based on total events or unique contacts.
- Whether there are any database update or calculation errors.

I would also compare the raw event data with the campaign statistics.

### 7.3 Is It a Bug or Expected Behaviour?

If the system counts total valid open events and never deletes or corrects old events, the opened count should not decrease.

So, in this situation, I would initially treat the decrease as a possible bug and investigate it.

However, if the system intentionally removes invalid or duplicate records, the decrease could be expected.

### 7.4 Why Did the Delivered Count Increase?

The delivered count increased from 970,000 to 975,000.

This is normal because some emails that were previously sent may have been delivered later.

The provider may send delivery events after the initial sent events. Therefore, the delivered count can increase over time.

**Conclusion:** I would investigate the decrease in opened count because it should normally not happen when counting total valid open events. The increase in delivered count is expected because delivery events can arrive later.