# aTimer

(c) 2020 Orange Number.
Written by Gon Yi. <https://gonyyi.com/copyright.txt>

aTimer is a library for jobs constantly running in the background.


---

## Intro

ATimer has two structs -- a `bucket` and a `job`.
A job can run itself without a bucket.
A bucket holds jobs together and control together for start and stop.


---

## Bucket

A bucket is a group of timer jobs controlled together.



---

## Job

A job is individual process of timer.
It can be created using `NewJob(func(*Job), time.Time) *Job`

A job can be started with/without a bucket by `Start(interval time.Duration)`
and can be stopped by `Stop()`.

A job can have its name set. This will be useful when used within the bucket.
If no name is given and linked to bucket, random name will be assigned.
A name can be checked by `GetName()string` or renamed by `Rename(string)`