# statestore-explorer
Explorer for bee statestore

```
Statestore explorer
----------------------
>> h

Actions:
        get             Get value for key
        count           Get count for prefix
        list            List values for prefix

```

With list we can browse more values using `Enter` key

Count and list take optional prefix

```
>> count <prefix>

>> count         (returns `all`)

>> list <prefix>

>> list          (returns `all`)

```

Also we can set start index. For this we need to use prefix explicitly. For no prefix use `all`

```
>> list all 100
....
....
....
---- Press 'Enter' to load more, any other key to exit ---- 150 50
```

