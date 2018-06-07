# Remind Me Plugin

## Details

* This plugin takes a time and message from a user and once that time has passed, it notifies the user of the message they wanted to be reminded with.
  * Time durations currently supported
    * minute(s)
    * hour(s)
    * day(s)

```
CMD: /remindme
Help: /remindme {time} {message}
```

### Examples

```
/remindme 1 minute "something I want to know about"
---------------
Your reminder is set for 2018 Jun 6 14:03:06 UTC
...One Minute Later...
something I want to know about
```

```
/remindme 2 hours "something I want to know about"
---------------
Your reminder is set for 2018 Jun 6 14:03:06 UTC
...Two Hours Later...
something I want to know about
```

```
/remindme 1 day "something I want to know about"
---------------
Your reminder is set for 2018 Jun 6 14:03:06 UTC
...One Day Later...
something I want to know about
```