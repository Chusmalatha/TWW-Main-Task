# BUGS.md

## Bug 1: Duplicate Events Were Counted Again

**What was the problem?**

The program was checking duplicate event IDs, but only within the current batch. When a new batch started, it forgot the IDs from the previous batch. Because of this, some duplicate events could be counted again.

**Why did it happen?**

The `seen` map was created inside the `processBatch()` function. So, every time this function ran, a new empty map was created.

**How did I fix it?**

I moved the duplicate tracking map outside the function and made it global. Now the program remembers the event IDs even when it moves to the next batch.

**How did I test it?**

I ran the program using the provided `events.jsonl` file and compared its output with `expected_output.txt`. The results matched.

---

## Bug 2: Unique Opens Were Calculated Incorrectly

**What was the problem?**

The program was checking only the contact ID when counting unique opens. If the same contact opened two different campaigns, the program counted that contact only once across both campaigns.

**Why did it happen?**

The code used only `ContactID` as the key in the `openedBy` map. But each campaign should have its own unique open count.

**How did I fix it?**

I created a key using both the campaign ID and contact ID.

```go
key := ev.CampaignID + ":" + ev.ContactID
```

Now the same contact can be counted once in each campaign.

**How did I test it?**

I checked the `unique_opens` values for all campaigns against the expected output. They matched.

---

## Bug 3: Delivered Events Were Grouped Under the Wrong Date

**What was the problem?**

The program was using my computer's local timezone to calculate the daily delivered count. This could put an event under the wrong date.

**Why did it happen?**

The code used `Local()` instead of `UTC()`.

**How did I fix it?**

I changed the code to:

```go
day := ev.Timestamp.UTC().Format("2006-01-02")
```

Now delivered events are grouped using their UTC date.

**How did I test it?**

I compared the daily delivered counts and dates with `expected_output.txt`. They matched.

---

## Bug 4: Multiple Workers Updated Counters at the Same Time

**What was the problem?**

The program uses multiple workers to process events. Sometimes, two workers could update the same campaign counter at the same time. This could cause incorrect counts.

**Why did it happen?**

There was no proper synchronization when workers updated the shared statistics.

**How did I fix it?**

I added a mutex using Go's `sync.Mutex`. Before updating the counters, a worker locks the mutex. After updating, it unlocks it.

This makes sure that only one worker updates the shared counters at a time.

**How did I test it?**

I ran the program with the Go race detector:

```powershell
$env:CGO_ENABLED = "1"
go run -race . events.jsonl
```

The program completed without reporting any data races.

---

## Final Verification

After fixing all four bugs, I ran the program again and compared its output with the provided `expected_output.txt` file.

The comparison showed no differences.

I also ran the race detector successfully. All four issues were fixed while keeping the original program structure.